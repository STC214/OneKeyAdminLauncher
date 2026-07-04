package process

import (
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
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	cwd, _ := os.Getwd()
	dir, _ := syscall.UTF16PtrFromString(cwd)
	r, _, _ := procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, uintptr(unsafe.Pointer(dir)), win.SW_SHOWNORMAL)
	return r > 32
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
		return nil
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
	file, _ := syscall.UTF16PtrFromString(target)
	var dirPtr *uint16
	if dir != "" && dir != "." {
		dirPtr, _ = syscall.UTF16PtrFromString(dir)
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

func CloseByName(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	pids := pidsByName(name)
	if len(pids) == 0 {
		return
	}
	for _, pid := range pids {
		postClose(pid)
	}
	deadline := time.Now().Add(2500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if len(pidsByName(name)) == 0 {
			return
		}
		time.Sleep(150 * time.Millisecond)
	}
	for _, pid := range pidsByName(name) {
		terminate(pid)
	}
}

func terminate(pid uint32) {
	h, _, _ := procOpenProcess.Call(PROCESS_TERMINATE, 0, uintptr(pid))
	if h == 0 {
		return
	}
	defer windows.CloseHandle(windows.Handle(h))
	procTerminateProcess.Call(h, 1)
}

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func pidsByName(name string) []uint32 {
	entries, err := processEntries()
	if err != nil {
		return nil
	}
	name = strings.ToLower(name)
	var pids []uint32
	for _, entry := range entries {
		if strings.ToLower(entry.Name) == name {
			pids = append(pids, entry.PID)
		}
	}
	return pids
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

func postClose(pid uint32) {
	cb := syscall.NewCallback(func(hwnd win.HWND, lparam uintptr) uintptr {
		var windowPID uint32
		win.GetWindowThreadProcessId(hwnd, &windowPID)
		if windowPID == pid && win.IsWindowVisible(hwnd) {
			win.PostMessage(hwnd, win.WM_CLOSE, 0, 0)
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
}
