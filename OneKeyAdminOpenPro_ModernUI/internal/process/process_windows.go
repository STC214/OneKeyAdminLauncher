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
	procWaitForSingleObj = kernel32.NewProc("WaitForSingleObject")
)

const (
	SEE_MASK_FLAG_NO_UI = 0x00000400
	PROCESS_TERMINATE   = 0x0001
	SYNCHRONIZE         = 0x00100000
	WAIT_OBJECT_0       = 0x00000000
	WAIT_TIMEOUT        = 0x00000102
	WAIT_FAILED         = 0xFFFFFFFF
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
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	pids, err := pidsByName(name)
	if err != nil {
		return fmt.Errorf("枚举进程 %s: %w", name, err)
	}
	if len(pids) == 0 {
		return nil
	}
	var closeErrs []error
	for _, pid := range pids {
		if err := postClose(pid); err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("PID %d 发送关闭消息: %w", pid, err))
		}
	}
	deadline := time.Now().Add(2500 * time.Millisecond)
	for time.Now().Before(deadline) {
		alive, err := targetsStillRunning(pids, name)
		if err != nil {
			return errors.Join(append(closeErrs, fmt.Errorf("确认进程状态: %w", err))...)
		}
		if len(alive) == 0 {
			return errors.Join(closeErrs...)
		}
		time.Sleep(150 * time.Millisecond)
	}
	for _, pid := range pids {
		if pidMatchesName(pid, name) {
			if err := terminate(pid); err != nil {
				closeErrs = append(closeErrs, fmt.Errorf("PID %d 强制结束: %w", pid, err))
			}
		}
	}
	alive, err := targetsStillRunning(pids, name)
	if err != nil {
		closeErrs = append(closeErrs, fmt.Errorf("确认强制结束结果: %w", err))
	} else if len(alive) != 0 {
		closeErrs = append(closeErrs, fmt.Errorf("进程仍在运行，PID: %v", alive))
	}
	return errors.Join(closeErrs...)
}

func pidMatchesName(pid uint32, name string) bool {
	entries, err := processEntries()
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.PID == pid {
			return strings.EqualFold(entry.Name, name)
		}
	}
	return false
}

func terminate(pid uint32) error {
	h, _, callErr := procOpenProcess.Call(PROCESS_TERMINATE|SYNCHRONIZE, 0, uintptr(pid))
	if h == 0 {
		return winCallError(callErr)
	}
	defer windows.CloseHandle(windows.Handle(h))
	r, _, callErr := procTerminateProcess.Call(h, 1)
	if r == 0 {
		return winCallError(callErr)
	}
	wait, _, callErr := procWaitForSingleObj.Call(h, 2000)
	switch uint32(wait) {
	case WAIT_OBJECT_0:
		return nil
	case WAIT_TIMEOUT:
		return errors.New("等待进程结束超时")
	case WAIT_FAILED:
		return winCallError(callErr)
	default:
		return fmt.Errorf("等待进程结束返回未知状态 0x%08X", uint32(wait))
	}
}

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func pidsByName(name string) ([]uint32, error) {
	entries, err := processEntries()
	if err != nil {
		return nil, err
	}
	name = strings.ToLower(name)
	var pids []uint32
	for _, entry := range entries {
		if strings.ToLower(entry.Name) == name {
			pids = append(pids, entry.PID)
		}
	}
	return pids, nil
}

func targetsStillRunning(pids []uint32, name string) ([]uint32, error) {
	entries, err := processEntries()
	if err != nil {
		return nil, err
	}
	targets := make(map[uint32]struct{}, len(pids))
	for _, pid := range pids {
		targets[pid] = struct{}{}
	}
	var alive []uint32
	for _, entry := range entries {
		if _, ok := targets[entry.PID]; ok && strings.EqualFold(entry.Name, name) {
			alive = append(alive, entry.PID)
		}
	}
	return alive, nil
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
