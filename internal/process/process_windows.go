package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procEnumWindows      = user32.NewProc("EnumWindows")
	procShellExecuteW    = shell32.NewProc("ShellExecuteW")
	procShellExecuteExW  = shell32.NewProc("ShellExecuteExW")
	procOpenProcess      = kernel32.NewProc("OpenProcess")
	procTerminateProcess = kernel32.NewProc("TerminateProcess")
)

const (
	SEE_MASK_FLAG_NO_UI = 0x00000400
	PROCESS_TERMINATE   = 0x0001
)

type Info struct {
	PID  uint32
	Name string
}

type shellExecuteInfo struct {
	CbSize     uint32
	FMask      uint32
	Hwnd       uintptr
	Verb       *uint16
	File       *uint16
	Parameters *uint16
	Directory  *uint16
	Show       int32
	HInstApp   uintptr
	IDList     uintptr
	Class      *uint16
	HkeyClass  uintptr
	HotKey     uint32
	Icon       uintptr
	Process    uintptr
}

func EnsureAdmin() bool {
	if isAdmin() {
		return true
	}
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	verb, err := syscall.UTF16PtrFromString("runas")
	if err != nil {
		return false
	}
	file, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		return false
	}
	dir, err := syscall.UTF16PtrFromString(filepath.Dir(exe))
	if err != nil {
		return false
	}
	procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, uintptr(unsafe.Pointer(dir)), win.SW_SHOWNORMAL)
	// Whether elevation starts successfully or the user cancels UAC, this
	// non-elevated process must exit. Only the elevated child may continue.
	return false
}

func isAdmin() bool {
	shell32 := windows.NewLazySystemDLL("shell32.dll")
	proc := shell32.NewProc("IsUserAnAdmin")
	ret, _, _ := proc.Call()
	return ret != 0
}

func ProcessNameForPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(path), "shell:") {
		return ""
	}
	base := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(base))
	if ext == ".exe" {
		return base
	}
	return ""
}

func Launch(path string, isUWP bool) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("程序路径为空")
	}
	if isUWP {
		target := path
		if !strings.HasPrefix(strings.ToLower(target), "shell:") {
			target = "shell:AppsFolder\\" + target
		}
		return shellExecute(target, "")
	}
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("程序文件不存在: %s", path)
		}
		return fmt.Errorf("检查程序文件 %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("程序路径是目录: %s", path)
	}
	dir := filepath.Dir(path)
	return shellExecute(path, dir)
}

func shellExecute(target, dir string) error {
	file, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	var dirPtr *uint16
	if dir != "" && dir != "." {
		dirPtr, err = syscall.UTF16PtrFromString(dir)
		if err != nil {
			return err
		}
	}
	verb := utf16Ptr("open")
	info := shellExecuteInfo{
		CbSize:    uint32(unsafe.Sizeof(shellExecuteInfo{})),
		FMask:     SEE_MASK_FLAG_NO_UI,
		Verb:      verb,
		File:      file,
		Directory: dirPtr,
		Show:      win.SW_SHOWNORMAL,
	}
	r, _, err := procShellExecuteExW.Call(uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		if err != syscall.Errno(0) {
			return err
		}
		return syscall.EINVAL
	}
	return nil
}

func List() ([]Info, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snap)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	if err := windows.Process32First(snap, &pe); err != nil {
		return nil, err
	}
	var out []Info
	seen := map[string]bool{}
	for {
		name := windows.UTF16ToString(pe.ExeFile[:])
		key := strings.ToLower(name)
		if name != "" && !seen[key] {
			out = append(out, Info{PID: pe.ProcessID, Name: name})
			seen[key] = true
		}
		if err := windows.Process32Next(snap, &pe); err != nil {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}

func IsRunning(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}
	procs, err := List()
	if err != nil {
		return false
	}
	for _, p := range procs {
		if strings.ToLower(p.Name) == name {
			return true
		}
	}
	return false
}

func CloseByName(name string) error {
	return CloseByNames([]string{name})
}

type closeTarget struct {
	pid  uint32
	name string
}

// CloseByNames closes all matching processes as one batch. Every process gets
// WM_CLOSE before the shared graceful-close timeout starts, so one slow process
// cannot delay sending WM_CLOSE to the next one.
func CloseByNames(names []string) error {
	wanted := make(map[string]string, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			key := strings.ToLower(name)
			if _, exists := wanted[key]; !exists {
				wanted[key] = name
			}
		}
	}
	if len(wanted) == 0 {
		return nil
	}

	entries, err := processEntries()
	if err != nil {
		return fmt.Errorf("枚举待关闭进程: %w", err)
	}
	var targets []closeTarget
	for _, entry := range entries {
		if name, ok := wanted[strings.ToLower(entry.Name)]; ok {
			targets = append(targets, closeTarget{pid: entry.PID, name: name})
		}
	}
	if len(targets) == 0 {
		return nil
	}

	var closeErrs []error
	for _, target := range targets {
		if err := postClose(target.pid); err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("%s (PID %d) 发送关闭消息: %w", target.name, target.pid, err))
		}
	}
	deadline := time.Now().Add(2500 * time.Millisecond)
	for time.Now().Before(deadline) {
		alive, err := closeTargetsStillRunning(targets)
		if err != nil {
			return errors.Join(append(closeErrs, fmt.Errorf("确认进程状态: %w", err))...)
		}
		if len(alive) == 0 {
			return errors.Join(closeErrs...)
		}
		time.Sleep(150 * time.Millisecond)
	}

	alive, err := closeTargetsStillRunning(targets)
	if err != nil {
		return errors.Join(append(closeErrs, fmt.Errorf("确认强制结束目标: %w", err))...)
	}
	for _, target := range alive {
		if err := terminate(target.pid); err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("%s (PID %d) 强制结束: %w", target.name, target.pid, err))
		}
	}
	forceDeadline := time.Now().Add(2000 * time.Millisecond)
	for time.Now().Before(forceDeadline) {
		alive, err = closeTargetsStillRunning(alive)
		if err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("确认强制结束结果: %w", err))
			return errors.Join(closeErrs...)
		}
		if len(alive) == 0 {
			return errors.Join(closeErrs...)
		}
		time.Sleep(50 * time.Millisecond)
	}
	if len(alive) != 0 {
		var remaining []string
		for _, target := range alive {
			remaining = append(remaining, fmt.Sprintf("%s (PID %d)", target.name, target.pid))
		}
		closeErrs = append(closeErrs, fmt.Errorf("进程仍在运行: %s", strings.Join(remaining, ", ")))
	}
	return errors.Join(closeErrs...)
}

func closeTargetsStillRunning(targets []closeTarget) ([]closeTarget, error) {
	entries, err := processEntries()
	if err != nil {
		return nil, err
	}
	running := make(map[uint32]string, len(entries))
	for _, entry := range entries {
		running[entry.PID] = entry.Name
	}
	var alive []closeTarget
	for _, target := range targets {
		if strings.EqualFold(running[target.pid], target.name) {
			alive = append(alive, target)
		}
	}
	return alive, nil
}

func terminate(pid uint32) error {
	h, _, callErr := procOpenProcess.Call(PROCESS_TERMINATE, 0, uintptr(pid))
	if h == 0 {
		return winCallError(callErr)
	}
	defer windows.CloseHandle(windows.Handle(h))
	r, _, callErr := procTerminateProcess.Call(h, 1)
	if r == 0 {
		return winCallError(callErr)
	}
	return nil
}

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func processEntries() ([]Info, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snap)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	if err := windows.Process32First(snap, &pe); err != nil {
		return nil, err
	}
	var out []Info
	for {
		name := windows.UTF16ToString(pe.ExeFile[:])
		if name != "" {
			out = append(out, Info{PID: pe.ProcessID, Name: name})
		}
		if err := windows.Process32Next(snap, &pe); err != nil {
			break
		}
	}
	return out, nil
}

func postClose(pid uint32) error {
	cb := syscall.NewCallback(func(hwnd win.HWND, lparam uintptr) uintptr {
		var windowPID uint32
		win.GetWindowThreadProcessId(hwnd, &windowPID)
		if windowPID == pid && win.IsWindowVisible(hwnd) {
			win.PostMessage(hwnd, win.WM_CLOSE, 0, 0)
		}
		return 1
	})
	r, _, callErr := procEnumWindows.Call(cb, 0)
	if r == 0 {
		return winCallError(callErr)
	}
	return nil
}

func winCallError(err error) error {
	if err == nil || err == syscall.Errno(0) {
		return syscall.EINVAL
	}
	return err
}
