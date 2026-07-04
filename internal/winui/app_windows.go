package winui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"program-launch-manager/internal/config"
	"program-launch-manager/internal/process"
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000
	WS_CHILD            = 0x40000000
	WS_TABSTOP          = 0x00010000
	WS_BORDER           = 0x00800000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_VSCROLL          = 0x00200000
	WS_EX_CLIENTEDGE    = 0x00000200
	WS_EX_TOOLWINDOW    = 0x00000080
	ES_AUTOHSCROLL      = 0x0080
	ES_READONLY         = 0x0800
	BS_PUSHBUTTON       = 0x00000000
	BS_AUTOCHECKBOX     = 0x00000003
	LBS_NOTIFY          = 0x00000001
	LBS_STANDARD        = 0x00A00003
	SS_CENTER           = 0x00000001

	SW_HIDE       = 0
	SW_SHOW       = 5
	SW_RESTORE    = 9
	SW_SHOWNORMAL = 1

	WM_CREATE          = 0x0001
	WM_DESTROY         = 0x0002
	WM_MOVE            = 0x0003
	WM_SIZE            = 0x0005
	WM_CLOSE           = 0x0010
	WM_GETMINMAXINFO   = 0x0024
	WM_ERASEBKGND      = 0x0014
	WM_COMMAND         = 0x0111
	WM_TIMER           = 0x0113
	WM_VSCROLL         = 0x0115
	WM_CTLCOLORSTATIC  = 0x0138
	WM_CTLCOLORLISTBOX = 0x0134
	WM_CTLCOLOREDIT    = 0x0133
	WM_CTLCOLORBTN     = 0x0135
	WM_MOUSEWHEEL      = 0x020A
	WM_LBUTTONDBLCLK   = 0x0203
	WM_RBUTTONUP       = 0x0205
	WM_APP             = 0x8000
	WM_SETFONT         = 0x0030
	WM_GETTEXT         = 0x000D
	WM_GETTEXTLENGTH   = 0x000E
	WM_SETREDRAW       = 0x000B
	WM_SETICON         = 0x0080

	BM_GETCHECK = 0x00F0
	BM_SETCHECK = 0x00F1
	BST_CHECKED = 1
	BN_CLICKED  = 0
	EN_CHANGE   = 0x0300

	OFN_PATHMUSTEXIST = 0x00000800
	OFN_FILEMUSTEXIST = 0x00001000
	OFN_EXPLORER      = 0x00080000

	LB_ADDSTRING     = 0x0180
	LB_GETCURSEL     = 0x0188
	LB_GETTEXT       = 0x0189
	LB_GETTEXTLEN    = 0x018A
	LB_SETCURSEL     = 0x0186
	LB_SETITEMHEIGHT = 0x01A0
	LB_GETITEMDATA   = 0x0199
	LB_SETITEMDATA   = 0x019A
	CBN_DBLCLK       = 2
	LBN_DBLCLK       = 2

	SIZE_MINIMIZED           = 1
	COLOR_WINDOW             = 5
	SB_VERT                  = 1
	SB_LINEUP                = 0
	SB_LINEDOWN              = 1
	SB_PAGEUP                = 2
	SB_PAGEDOWN              = 3
	SB_THUMBPOSITION         = 4
	SB_THUMBTRACK            = 5
	SB_TOP                   = 6
	SB_BOTTOM                = 7
	SIF_RANGE                = 0x0001
	SIF_PAGE                 = 0x0002
	SIF_POS                  = 0x0004
	SIF_ALL                  = SIF_RANGE | SIF_PAGE | SIF_POS
	WHEEL_DELTA              = 120
	RDW_INVALIDATE           = 0x0001
	RDW_ERASE                = 0x0004
	RDW_ALLCHILDREN          = 0x0080
	IDC_ARROW                = 32512
	IDI_APPLICATION          = 32512
	ICON_SMALL               = 0
	ICON_BIG                 = 1
	APP_ICON_ID              = 101
	SM_XVIRTUAL              = 76
	SM_YVIRTUAL              = 77
	SM_CXSCREEN              = 0
	SM_CXVIRTUAL             = 78
	SM_CYVIRTUAL             = 79
	MONITOR_DEFAULTTONEAREST = 2

	WM_TRAY_ICON    = WM_APP + 10
	WM_STATUS_READY = WM_APP + 11
	WM_SAVE_DONE    = WM_APP + 12
	ID_TIMER        = 2001
	TRAY_ID         = 1

	NIM_ADD     = 0
	NIM_DELETE  = 2
	NIF_MESSAGE = 1
	NIF_ICON    = 2
	NIF_TIP     = 4

	ID_START = 100
	ID_STOP  = 101
	ID_ADD   = 102
	ID_ALL   = 103
	ID_NONE  = 104
	ID_SAVE  = 105

	rowBase          = 1000
	rowStep          = 10
	rowStartY        = 96
	rowHeight        = 64
	rowBottomPadding = 16
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")

	procDefWindowProc     = user32.NewProc("DefWindowProcW")
	procRegisterClassEx   = user32.NewProc("RegisterClassExW")
	procCreateWindowEx    = user32.NewProc("CreateWindowExW")
	procShowWindow        = user32.NewProc("ShowWindow")
	procUpdateWindow      = user32.NewProc("UpdateWindow")
	procGetMessage        = user32.NewProc("GetMessageW")
	procTranslateMessage  = user32.NewProc("TranslateMessage")
	procDispatchMessage   = user32.NewProc("DispatchMessageW")
	procPostQuitMessage   = user32.NewProc("PostQuitMessage")
	procPostMessage       = user32.NewProc("PostMessageW")
	procSendMessage       = user32.NewProc("SendMessageW")
	procSetWindowText     = user32.NewProc("SetWindowTextW")
	procGetWindowText     = user32.NewProc("GetWindowTextW")
	procGetWindowTextLen  = user32.NewProc("GetWindowTextLengthW")
	procMoveWindow        = user32.NewProc("MoveWindow")
	procRedrawWindow      = user32.NewProc("RedrawWindow")
	procDestroyWindow     = user32.NewProc("DestroyWindow")
	procLoadCursor        = user32.NewProc("LoadCursorW")
	procLoadIcon          = user32.NewProc("LoadIconW")
	procSetTimer          = user32.NewProc("SetTimer")
	procKillTimer         = user32.NewProc("KillTimer")
	procSetScrollInfo     = user32.NewProc("SetScrollInfo")
	procGetClientRect     = user32.NewProc("GetClientRect")
	procGetWindowRect     = user32.NewProc("GetWindowRect")
	procMonitorFromWindow = user32.NewProc("MonitorFromWindow")
	procGetMonitorInfo    = user32.NewProc("GetMonitorInfoW")
	procGetSystemMetrics  = user32.NewProc("GetSystemMetrics")
	procIsIconic          = user32.NewProc("IsIconic")
	procIsWindowVisible   = user32.NewProc("IsWindowVisible")
	procFillRect          = user32.NewProc("FillRect")
	procMessageBox        = user32.NewProc("MessageBoxW")
	procEnableWindow      = user32.NewProc("EnableWindow")
	procSetForeground     = user32.NewProc("SetForegroundWindow")
	procCreatePopupMenu   = user32.NewProc("CreatePopupMenu")
	procAppendMenu        = user32.NewProc("AppendMenuW")
	procTrackPopupMenu    = user32.NewProc("TrackPopupMenu")
	procGetCursorPos      = user32.NewProc("GetCursorPos")
	procDestroyMenu       = user32.NewProc("DestroyMenu")
	procGetDlgItem        = user32.NewProc("GetDlgItem")
	procGetParent         = user32.NewProc("GetParent")

	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	procRtlMoveMemory    = kernel32.NewProc("RtlMoveMemory")
	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procSetTextColor     = gdi32.NewProc("SetTextColor")
	procSetBkColor       = gdi32.NewProc("SetBkColor")
	procCreateFont       = gdi32.NewProc("CreateFontW")
	procShellNotifyIcon  = shell32.NewProc("Shell_NotifyIconW")
	procGetOpenFileName  = comdlg32.NewProc("GetOpenFileNameW")
)

type point struct{ X, Y int32 }
type rect struct{ Left, Top, Right, Bottom int32 }
type monitorInfo struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
}
type msg struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}
type wndclassex struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}
type notifyIconData struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
}
type scrollInfo struct {
	CbSize    uint32
	FMask     uint32
	NMin      int32
	NMax      int32
	NPage     uint32
	NPos      int32
	NTrackPos int32
}
type openFileName struct {
	StructSize    uint32
	HwndOwner     uintptr
	Instance      uintptr
	Filter        *uint16
	CustomFilter  *uint16
	MaxCustFilter uint32
	FilterIndex   uint32
	File          *uint16
	MaxFile       uint32
	FileTitle     *uint16
	MaxFileTitle  uint32
	InitialDir    *uint16
	Title         *uint16
	Flags         uint32
	FileOffset    uint16
	FileExtension uint16
	DefExt        *uint16
	CustData      uintptr
	FnHook        uintptr
	TemplateName  *uint16
	Reserved      unsafe.Pointer
	Reserved2     uint32
	FlagsEx       uint32
}

type rowControls struct {
	check         uintptr
	path          uintptr
	browse        uintptr
	uwp           uintptr
	status        uintptr
	selectProcess uintptr
	remove        uintptr
	index         int
	visible       bool
}

type appState struct {
	hwnd       uintptr
	font       uintptr
	icon       uintptr
	bgBrush    uintptr
	editBrush  uintptr
	panelBrush uintptr
	cfg        config.File
	rows       []rowControls
	scrollY    int32
	mu         sync.Mutex
	saveMu     sync.Mutex

	statusMu        sync.Mutex
	statusBusy      bool
	runningSnapshot map[string]bool
	updatingUI      bool
	autoSaveDirty   bool
	autoSaveBusy    bool
}

var app *appState

func Run() int {
	cfg, _ := config.Load()
	app = &appState{cfg: cfg}
	if app.cfg.Window.W == 0 {
		app.cfg.Window = config.DefaultWindow()
	}
	return runWindow()
}

func runWindow() int {
	hinst := getModuleHandle()
	className := utf16Ptr("ProgramLaunchManagerGo")
	icon := loadAppIcon(hinst)
	wc := wndclassex{
		Size:      uint32(unsafe.Sizeof(wndclassex{})),
		WndProc:   syscall.NewCallback(wndProc),
		Instance:  hinst,
		Cursor:    loadCursor(0, IDC_ARROW),
		Icon:      icon,
		IconSm:    icon,
		ClassName: className,
	}
	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	win := app.cfg.Window
	if win.W < 640 {
		win.W = 800
	}
	if minW := startupMinWindowWidth(); win.W < minW {
		win.W = minW
	}
	if win.H < 360 {
		win.H = 594
	}
	x, y := win.X, win.Y
	if !isWindowRectVisible(x, y, win.W, win.H) {
		x, y = 100, 100
	}
	hwnd := createWindow(0, className, utf16Ptr("程序启动管理器"), WS_OVERLAPPEDWINDOW|WS_VISIBLE|WS_VSCROLL, x, y, win.W, win.H, 0, 0, hinst, 0)
	app.hwnd = hwnd
	procShowWindow.Call(hwnd, SW_SHOWNORMAL)
	procUpdateWindow.Call(hwnd)
	var m msg
	for {
		r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
	return int(m.WParam)
}

func wndProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_CREATE:
		app.hwnd = hwnd
		app.bgBrush = createSolidBrush(rgb(43, 43, 43))
		app.editBrush = createSolidBrush(rgb(54, 54, 54))
		app.panelBrush = createSolidBrush(rgb(36, 37, 40))
		app.font = createFont("Microsoft YaHei UI", 14)
		app.icon = loadAppIcon(getModuleHandle())
		procSendMessage.Call(hwnd, WM_SETICON, ICON_BIG, app.icon)
		procSendMessage.Call(hwnd, WM_SETICON, ICON_SMALL, app.icon)
		createToolbar(hwnd)
		rebuildRows(hwnd)
		layout(hwnd, true)
		addTray(hwnd)
		procSetTimer.Call(hwnd, ID_TIMER, 1000, 0)
		requestStatusRefresh()
		return 0
	case WM_SIZE:
		if wparam == SIZE_MINIMIZED {
			hideToTray(hwnd)
			return 0
		}
		layout(hwnd, true)
		return 0
	case WM_GETMINMAXINFO:
		applyMinWindowSize(hwnd, lparam)
		return 0
	case WM_VSCROLL:
		handleScroll(hwnd, loword(uint32(wparam)), hiword(uint32(wparam)))
		return 0
	case WM_MOUSEWHEEL:
		handleMouseWheel(hwnd, hiword(uint32(wparam)))
		return 0
	case WM_TIMER:
		requestStatusRefresh()
		flushAutoSave()
		return 0
	case WM_STATUS_READY:
		applyStatusSnapshot()
		return 0
	case WM_SAVE_DONE:
		message(hwnd, "配置已保存。", "程序启动管理器")
		return 0
	case WM_COMMAND:
		if app != nil && app.updatingUI {
			return 0
		}
		handleCommand(hwnd, loword(uint32(wparam)), hiword(uint32(wparam)), lparam)
		return 0
	case WM_TRAY_ICON:
		if lparam == WM_LBUTTONDBLCLK {
			restoreFromTray(hwnd)
			return 0
		}
		if lparam == WM_RBUTTONUP {
			showTrayMenu(hwnd)
			return 0
		}
	case WM_ERASEBKGND:
		var rc rect
		procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		procFillRect.Call(wparam, uintptr(unsafe.Pointer(&rc)), app.bgBrush)
		return 1
	case WM_CTLCOLORSTATIC, WM_CTLCOLORLISTBOX, WM_CTLCOLOREDIT, WM_CTLCOLORBTN:
		color := rgb(220, 220, 220)
		if isStatusControl(lparam) {
			txt := getText(lparam)
			if txt != "" && txt != "未运行" && txt != "未绑定" {
				color = rgb(70, 220, 90)
			} else {
				color = rgb(150, 150, 150)
			}
		}
		procSetTextColor.Call(wparam, uintptr(color))
		procSetBkColor.Call(wparam, uintptr(rgb(54, 54, 54)))
		return app.editBrush
	case WM_CLOSE:
		saveFromUI()
		removeTray(hwnd)
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		procKillTimer.Call(hwnd, ID_TIMER)
		if app.bgBrush != 0 {
			procDeleteObject.Call(app.bgBrush)
		}
		if app.editBrush != 0 {
			procDeleteObject.Call(app.editBrush)
		}
		if app.panelBrush != 0 {
			procDeleteObject.Call(app.panelBrush)
		}
		if app.font != 0 {
			procDeleteObject.Call(app.font)
		}
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProc.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}

func createToolbar(hwnd uintptr) {
	button(hwnd, "一键开启", 16, 38, 82, 28, ID_START)
	button(hwnd, "一键关闭", 108, 38, 82, 28, ID_STOP)
	button(hwnd, "添加程序", 200, 38, 82, 28, ID_ADD)
	button(hwnd, "全选", 292, 38, 50, 28, ID_ALL)
	button(hwnd, "全不选", 352, 38, 64, 28, ID_NONE)
	button(hwnd, "保存配置", 426, 38, 82, 28, ID_SAVE)
}

func rebuildRows(hwnd uintptr) {
	for _, r := range app.rows {
		for _, h := range []uintptr{r.check, r.path, r.browse, r.uwp, r.status, r.selectProcess, r.remove} {
			if h != 0 {
				procDestroyWindow.Call(h)
			}
		}
	}
	app.rows = nil
}

func layout(hwnd uintptr, erase bool) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	width := rc.Right - rc.Left
	clampScroll(hwnd, rc.Bottom-rc.Top)
	ensureRowPool(hwnd, visibleRowCapacity(rc.Bottom-rc.Top))
	first := firstVisibleIndex()
	for i := range app.rows {
		r := &app.rows[i]
		idx := first + i
		if idx >= len(app.cfg.Programs) {
			setRowVisible(r, false)
			r.index = -1
			continue
		}
		bindRow(r, idx)
		y := int32(rowStartY+idx*rowHeight) - app.scrollY
		visible := y >= rowStartY && y < rc.Bottom
		move(r.check, 20, y+8, 18, 18)
		pathW := width - 550
		if pathW < 250 {
			pathW = 250
		}
		move(r.path, 44, y, pathW, 30)
		x := int32(44) + pathW + 6
		move(r.browse, x, y, 78, 30)
		x += 82
		move(r.uwp, x, y+4, 56, 22)
		x += 60
		statusW := int32(158)
		if width < 760 {
			statusW = 120
		}
		move(r.status, x, y, statusW, 30)
		x += statusW + 6
		move(r.selectProcess, x, y, 86, 30)
		x += 92
		move(r.remove, x, y, 72, 30)
		setRowVisible(r, visible)
	}
	updateScrollBar(hwnd, rc.Bottom-rc.Top)
	redraw(hwnd, erase)
}

func ensureRowPool(hwnd uintptr, count int) {
	for len(app.rows) < count {
		slot := len(app.rows)
		id := rowBase + slot*rowStep
		app.rows = append(app.rows, rowControls{
			check:         checkbox(hwnd, 0, 0, 18, 18, id+1),
			path:          edit(hwnd, "", 0, 0, 260, 30, id+2, false),
			browse:        button(hwnd, "浏览", 0, 0, 78, 30, id+3),
			uwp:           checkboxText(hwnd, "UWP", 0, 0, 56, 22, id+7),
			status:        edit(hwnd, "", 0, 0, 158, 30, id+4, true),
			selectProcess: button(hwnd, "选择进程", 0, 0, 86, 30, id+5),
			remove:        button(hwnd, "×", 0, 0, 72, 30, id+6),
			index:         -1,
			visible:       true,
		})
	}
}

func bindRow(r *rowControls, idx int) {
	if idx < 0 || idx >= len(app.cfg.Programs) {
		return
	}
	if r.index == idx {
		return
	}
	item := app.cfg.Programs[idx]
	app.updatingUI = true
	defer func() { app.updatingUI = false }()
	setTextIfChanged(r.path, item.Path)
	setCheckIfChanged(r.check, item.Enabled)
	setCheckIfChanged(r.uwp, item.IsUWP)
	setTextIfChanged(r.status, currentStatusText(item))
	r.index = idx
}

func handleScroll(hwnd uintptr, code uint16, thumb uint16) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	listH := listViewportHeight(rc.Bottom - rc.Top)
	pos := app.scrollY
	switch code {
	case SB_LINEUP:
		pos -= rowHeight
	case SB_LINEDOWN:
		pos += rowHeight
	case SB_PAGEUP:
		pos -= scrollPage(listH)
	case SB_PAGEDOWN:
		pos += scrollPage(listH)
	case SB_THUMBPOSITION, SB_THUMBTRACK:
		pos = int32(thumb)
	case SB_TOP:
		pos = 0
	case SB_BOTTOM:
		pos = maxScroll(listH)
	}
	setScroll(hwnd, pos)
}

func handleMouseWheel(hwnd uintptr, deltaWord uint16) {
	delta := int16(deltaWord)
	if delta == 0 {
		return
	}
	lines := int32(delta) / WHEEL_DELTA
	if lines == 0 {
		if delta > 0 {
			lines = 1
		} else {
			lines = -1
		}
	}
	setScroll(hwnd, app.scrollY-lines*rowHeight)
}

func setScroll(hwnd uintptr, pos int32) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	listH := listViewportHeight(rc.Bottom - rc.Top)
	max := snapScroll(maxScroll(listH))
	if pos < 0 {
		pos = 0
	}
	pos = snapScroll(pos)
	if pos > max {
		pos = max
	}
	if pos == app.scrollY {
		updateScrollBar(hwnd, rc.Bottom-rc.Top)
		return
	}
	syncFromUI()
	app.scrollY = pos
	layout(hwnd, false)
}

func scrollToBottom(hwnd uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	app.scrollY = snapScroll(maxScroll(listViewportHeight(rc.Bottom - rc.Top)))
}

func clampScroll(hwnd uintptr, clientH int32) {
	max := snapScroll(maxScroll(listViewportHeight(clientH)))
	if app.scrollY > max {
		app.scrollY = max
	}
	if app.scrollY < 0 {
		app.scrollY = 0
	}
	app.scrollY = snapScroll(app.scrollY)
}

func updateScrollBar(hwnd uintptr, clientH int32) {
	listH := listViewportHeight(clientH)
	max := contentHeight() - 1
	if max < 0 {
		max = 0
	}
	page := uint32(listH)
	if page < 1 {
		page = 1
	}
	si := scrollInfo{
		CbSize: uint32(unsafe.Sizeof(scrollInfo{})),
		FMask:  SIF_ALL,
		NMin:   0,
		NMax:   max,
		NPage:  page,
		NPos:   app.scrollY,
	}
	procSetScrollInfo.Call(hwnd, SB_VERT, uintptr(unsafe.Pointer(&si)), 1)
}

func scrollPage(clientH int32) int32 {
	page := clientH
	if page < rowHeight {
		return rowHeight
	}
	return snapScroll(page)
}

func maxScroll(clientH int32) int32 {
	max := contentHeight() - clientH
	if max < 0 {
		return 0
	}
	return max
}

func contentHeight() int32 {
	if len(app.cfg.Programs) == 0 {
		return 0
	}
	return int32(len(app.cfg.Programs)*rowHeight + rowBottomPadding)
}

func listViewportHeight(clientH int32) int32 {
	h := clientH - rowStartY
	if h < rowHeight {
		return rowHeight
	}
	return h
}

func visibleRowCapacity(clientH int32) int {
	if len(app.cfg.Programs) == 0 {
		return 0
	}
	capacity := int(listViewportHeight(clientH)/rowHeight) + 2
	if capacity > len(app.cfg.Programs) {
		capacity = len(app.cfg.Programs)
	}
	if capacity < 0 {
		return 0
	}
	return capacity
}

func currentClientHeight(hwnd uintptr) int32 {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	return rc.Bottom - rc.Top
}

func firstVisibleIndex() int {
	if app.scrollY <= 0 {
		return 0
	}
	idx := int(app.scrollY / rowHeight)
	if idx < 0 {
		return 0
	}
	if idx > len(app.cfg.Programs) {
		return len(app.cfg.Programs)
	}
	return idx
}

func currentStatusText(item config.ProgramItem) string {
	name := effectiveProcess(item)
	if name == "" {
		return "未绑定"
	}
	app.statusMu.Lock()
	running := app.runningSnapshot[strings.ToLower(name)]
	app.statusMu.Unlock()
	if running {
		return name
	}
	return "未运行"
}

func snapScroll(pos int32) int32 {
	if pos <= 0 {
		return 0
	}
	return (pos / rowHeight) * rowHeight
}

func setRowVisible(r *rowControls, visible bool) {
	if r.visible == visible {
		return
	}
	show := uintptr(SW_HIDE)
	if visible {
		show = SW_SHOW
	}
	for _, h := range []uintptr{r.check, r.path, r.browse, r.uwp, r.status, r.selectProcess, r.remove} {
		procShowWindow.Call(h, show)
	}
	r.visible = visible
}

func handleCommand(hwnd uintptr, id uint16, code uint16, lparam uintptr) {
	needsRefresh := false
	switch id {
	case ID_START:
		syncFromUI()
		markAutoSaveDirty()
		for _, item := range cloneConfig(app.cfg).Programs {
			if item.Enabled {
				go process.Launch(item.Path, item.IsUWP)
			}
		}
		needsRefresh = true
	case ID_STOP:
		syncFromUI()
		markAutoSaveDirty()
		for _, item := range cloneConfig(app.cfg).Programs {
			if item.Enabled {
				if name := effectiveProcess(item); name != "" {
					go process.CloseByName(name)
				}
			}
		}
		needsRefresh = true
	case ID_ADD:
		syncFromUI()
		app.cfg.Programs = append(app.cfg.Programs, config.ProgramItem{Enabled: true})
		rebuildRows(hwnd)
		scrollToBottom(hwnd)
		layout(hwnd, true)
		markAutoSaveDirty()
		needsRefresh = true
	case ID_ALL:
		setAllEnabled(true)
		syncFromUI()
		markAutoSaveDirty()
	case ID_NONE:
		setAllEnabled(false)
		syncFromUI()
		markAutoSaveDirty()
	case ID_SAVE:
		syncFromUI()
		saveConfigNotifyAsync(hwnd, cloneConfig(app.cfg))
	default:
		if id >= rowBase {
			slot := int((id - rowBase) / rowStep)
			part := int((id - rowBase) % rowStep)
			if slot < 0 || slot >= len(app.rows) {
				return
			}
			idx := app.rows[slot].index
			if idx < 0 || idx >= len(app.cfg.Programs) {
				return
			}
			switch part {
			case 1, 7:
				if code == BN_CLICKED {
					syncFromUI()
					scheduleAutoSave()
				}
			case 2:
				if code == EN_CHANGE {
					scheduleAutoSave()
				}
			case 3:
				if path := chooseFile(hwnd); path != "" {
					setText(app.rows[slot].path, path)
					if app.cfg.Programs[idx].SelectedProcess == "" {
						app.cfg.Programs[idx].ProcessName = process.ProcessNameForPath(path)
					}
					syncFromUI()
					markAutoSaveDirty()
					needsRefresh = true
				}
			case 5:
				if name := chooseProcess(hwnd); name != "" {
					syncFromUI()
					app.cfg.Programs[idx].SelectedProcess = name
					markAutoSaveDirty()
					needsRefresh = true
				}
			case 6:
				syncFromUI()
				app.cfg.Programs = append(app.cfg.Programs[:idx], app.cfg.Programs[idx+1:]...)
				rebuildRows(hwnd)
				clampScroll(hwnd, currentClientHeight(hwnd))
				layout(hwnd, true)
				markAutoSaveDirty()
				needsRefresh = true
			}
			_ = code
			_ = lparam
		}
	}
	if needsRefresh {
		requestStatusRefresh()
	}
}

func saveFromUI() {
	syncFromUI()
	_ = saveConfig(cloneConfig(app.cfg))
}

func saveConfigNotifyAsync(hwnd uintptr, cfg config.File) {
	go func() {
		_ = saveConfig(cfg)
		procPostMessage.Call(hwnd, WM_SAVE_DONE, 0, 0)
	}()
}

func scheduleAutoSave() {
	syncFromUI()
	markAutoSaveDirty()
}

func setAllEnabled(enabled bool) {
	check := uintptr(0)
	if enabled {
		check = BST_CHECKED
	}
	app.updatingUI = true
	defer func() { app.updatingUI = false }()
	for _, r := range app.rows {
		if r.index >= 0 && procSendMessageUint(r.check, BM_GETCHECK, 0, 0) != check {
			procSendMessage.Call(r.check, BM_SETCHECK, check, 0)
		}
	}
	for i := range app.cfg.Programs {
		app.cfg.Programs[i].Enabled = enabled
	}
}

func saveConfig(cfg config.File) error {
	if app != nil {
		app.saveMu.Lock()
		defer app.saveMu.Unlock()
	}
	return config.Save(cfg)
}

func markAutoSaveDirty() {
	app.mu.Lock()
	app.autoSaveDirty = true
	app.mu.Unlock()
}

func flushAutoSave() {
	app.mu.Lock()
	if !app.autoSaveDirty || app.autoSaveBusy {
		app.mu.Unlock()
		return
	}
	cfg := cloneConfig(app.cfg)
	app.autoSaveDirty = false
	app.autoSaveBusy = true
	app.mu.Unlock()

	go func() {
		_ = saveConfig(cfg)
		app.mu.Lock()
		app.autoSaveBusy = false
		app.mu.Unlock()
	}()
}

func cloneConfig(cfg config.File) config.File {
	out := cfg
	if cfg.Programs != nil {
		out.Programs = append([]config.ProgramItem(nil), cfg.Programs...)
	}
	return out
}

func syncFromUI() {
	app.mu.Lock()
	defer app.mu.Unlock()
	for _, r := range app.rows {
		if r.index < 0 || r.index >= len(app.cfg.Programs) {
			continue
		}
		item := &app.cfg.Programs[r.index]
		item.Path = getText(r.path)
		item.Enabled = procSendMessageUint(r.check, BM_GETCHECK, 0, 0) == BST_CHECKED
		item.IsUWP = procSendMessageUint(r.uwp, BM_GETCHECK, 0, 0) == BST_CHECKED || strings.HasPrefix(strings.ToLower(item.Path), "shell:")
		if item.ProcessName == "" {
			item.ProcessName = process.ProcessNameForPath(item.Path)
		}
	}
	if isWindowStateSavable(app.hwnd) {
		var wr rect
		procGetWindowRect.Call(app.hwnd, uintptr(unsafe.Pointer(&wr)))
		if wr.Right > wr.Left && wr.Bottom > wr.Top {
			app.cfg.Window = config.WindowState{X: wr.Left, Y: wr.Top, W: wr.Right - wr.Left, H: wr.Bottom - wr.Top}
		}
	}
}

func requestStatusRefresh() {
	syncFromUI()

	app.statusMu.Lock()
	if app.statusBusy {
		app.statusMu.Unlock()
		return
	}
	app.statusBusy = true
	app.statusMu.Unlock()

	hwnd := app.hwnd
	go func() {
		running := map[string]bool{}
		if procs, err := process.List(); err == nil {
			for _, p := range procs {
				running[strings.ToLower(p.Name)] = true
			}
		}

		app.statusMu.Lock()
		app.runningSnapshot = running
		app.statusBusy = false
		app.statusMu.Unlock()

		procPostMessage.Call(hwnd, WM_STATUS_READY, 0, 0)
	}()
}

func applyStatusSnapshot() {
	app.statusMu.Lock()
	running := map[string]bool{}
	for k, v := range app.runningSnapshot {
		running[k] = v
	}
	app.statusMu.Unlock()

	app.updatingUI = true
	defer func() { app.updatingUI = false }()

	for _, r := range app.rows {
		if r.index < 0 || r.index >= len(app.cfg.Programs) {
			continue
		}
		item := app.cfg.Programs[r.index]
		name := effectiveProcess(item)
		if name == "" {
			setTextIfChanged(r.status, "未绑定")
			continue
		}
		if running[strings.ToLower(name)] {
			setTextIfChanged(r.status, name)
		} else {
			setTextIfChanged(r.status, "未运行")
		}
	}
}

func effectiveProcess(item config.ProgramItem) string {
	if item.SelectedProcess != "" {
		return item.SelectedProcess
	}
	if item.ProcessName != "" {
		return item.ProcessName
	}
	return process.ProcessNameForPath(item.Path)
}

func chooseFile(hwnd uintptr) string {
	var buf [4096]uint16
	filterData := syscall.StringToUTF16("程序 (*.exe;*.lnk)\x00*.exe;*.lnk\x00所有文件 (*.*)\x00*.*\x00\x00")
	filter := &filterData[0]
	title := utf16Ptr("选择程序")
	ofn := openFileName{
		StructSize: uint32(unsafe.Sizeof(openFileName{})),
		HwndOwner:  hwnd,
		Filter:     filter,
		File:       &buf[0],
		MaxFile:    uint32(len(buf)),
		Title:      title,
		Flags:      OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST,
	}
	r, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(&ofn)))
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func chooseProcess(parent uintptr) string {
	procs, _ := process.List()
	d := newProcessDialog(parent, procs)
	return d.run()
}

type processDialog struct {
	parent    uintptr
	hwnd      uintptr
	title     uintptr
	subtitle  uintptr
	countText uintptr
	list      uintptr
	result    string
	done      bool
	procs     []process.Info
	titleFont uintptr
	smallFont uintptr
	listFont  uintptr
}

var dlg *processDialog

func newProcessDialog(parent uintptr, procs []process.Info) *processDialog {
	return &processDialog{parent: parent, procs: procs}
}

func (d *processDialog) run() string {
	dlg = d
	hinst := getModuleHandle()
	className := utf16Ptr("ProgramLaunchProcessDialog")
	icon := loadAppIcon(hinst)
	wc := wndclassex{Size: uint32(unsafe.Sizeof(wndclassex{})), WndProc: syscall.NewCallback(processDlgProc), Instance: hinst, Cursor: loadCursor(0, IDC_ARROW), Icon: icon, IconSm: icon, ClassName: className}
	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	x, y := centeredDialogPos(d.parent, 620, 680)
	d.hwnd = createWindow(WS_EX_TOOLWINDOW, className, utf16Ptr("选择进程"), WS_CAPTION|WS_SYSMENU|WS_VISIBLE, x, y, 620, 680, d.parent, 0, hinst, 0)
	procEnableWindow.Call(d.parent, 0)
	var m msg
	for !d.done {
		r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
	procEnableWindow.Call(d.parent, 1)
	procSetForeground.Call(d.parent)
	return d.result
}

func processDlgProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_CREATE:
		dlg.titleFont = createFontWeight("Microsoft YaHei UI", 22, 600)
		dlg.smallFont = createFontWeight("Microsoft YaHei UI", 13, 400)
		dlg.listFont = createFontWeight("Microsoft YaHei UI", 15, 400)
		dlg.title = staticText(hwnd, "选择进程", 24, 22, 260, 32, 300)
		dlg.subtitle = staticText(hwnd, "绑定后用于运行状态检测和一键关闭", 24, 58, 420, 24, 304)
		dlg.countText = staticText(hwnd, fmt.Sprintf("当前进程  %d", len(dlg.procs)), 470, 60, 120, 22, 305)
		procSendMessage.Call(dlg.title, WM_SETFONT, dlg.titleFont, 1)
		procSendMessage.Call(dlg.subtitle, WM_SETFONT, dlg.smallFont, 1)
		procSendMessage.Call(dlg.countText, WM_SETFONT, dlg.smallFont, 1)
		dlg.list = createChild(WS_EX_CLIENTEDGE, "LISTBOX", "", WS_CHILD|WS_VISIBLE|WS_TABSTOP|WS_VSCROLL|WS_BORDER|LBS_NOTIFY, 24, 104, 572, 480, hwnd, 301)
		procSendMessage.Call(dlg.list, WM_SETFONT, dlg.listFont, 1)
		procSendMessage.Call(dlg.list, LB_SETITEMHEIGHT, 0, 28)
		for i, p := range dlg.procs {
			text := fmt.Sprintf("%s    PID %d", p.Name, p.PID)
			row, _, _ := procSendMessage.Call(dlg.list, LB_ADDSTRING, 0, uintptr(unsafe.Pointer(utf16Ptr(text))))
			if row != ^uintptr(0) {
				procSendMessage.Call(dlg.list, LB_SETITEMDATA, row, uintptr(i))
			}
		}
		procSendMessage.Call(dlg.list, LB_SETCURSEL, 0, 0)
		button(hwnd, "确定", 410, 606, 86, 34, 302)
		button(hwnd, "取消", 510, 606, 86, 34, 303)
		return 0
	case WM_COMMAND:
		id := loword(uint32(wparam))
		code := hiword(uint32(wparam))
		if id == 302 || (id == 301 && code == LBN_DBLCLK) {
			sel := int(procSendMessageUint(dlg.list, LB_GETCURSEL, 0, 0))
			if sel >= 0 {
				item := procSendMessageUint(dlg.list, LB_GETITEMDATA, uintptr(sel), 0)
				if item != ^uintptr(0) && int(item) < len(dlg.procs) {
					dlg.result = dlg.procs[int(item)].Name
				}
			}
			dlg.done = true
			procDestroyWindow.Call(hwnd)
			return 0
		}
		if id == 303 {
			dlg.done = true
			procDestroyWindow.Call(hwnd)
			return 0
		}
	case WM_CLOSE:
		dlg.done = true
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		if dlg.titleFont != 0 {
			procDeleteObject.Call(dlg.titleFont)
			dlg.titleFont = 0
		}
		if dlg.smallFont != 0 {
			procDeleteObject.Call(dlg.smallFont)
			dlg.smallFont = 0
		}
		if dlg.listFont != 0 {
			procDeleteObject.Call(dlg.listFont)
			dlg.listFont = 0
		}
		return 0
	case WM_ERASEBKGND:
		var rc rect
		procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
		procFillRect.Call(wparam, uintptr(unsafe.Pointer(&rc)), app.bgBrush)
		return 1
	case WM_CTLCOLORSTATIC, WM_CTLCOLORLISTBOX, WM_CTLCOLOREDIT, WM_CTLCOLORBTN:
		color := rgb(220, 220, 220)
		if lparam == dlg.title {
			color = rgb(245, 246, 248)
		}
		if lparam == dlg.subtitle || lparam == dlg.countText {
			color = rgb(158, 166, 176)
		}
		procSetTextColor.Call(wparam, uintptr(color))
		if msg == WM_CTLCOLORLISTBOX {
			procSetBkColor.Call(wparam, uintptr(rgb(36, 37, 40)))
			return app.panelBrush
		}
		procSetBkColor.Call(wparam, uintptr(rgb(43, 43, 43)))
		return app.bgBrush
	}
	r, _, _ := procDefWindowProc.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}

func addTray(hwnd uintptr) {
	var data notifyIconData
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.HWnd = hwnd
	data.UID = TRAY_ID
	data.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	data.UCallbackMessage = WM_TRAY_ICON
	data.HIcon = app.icon
	copy(data.SzTip[:], syscall.StringToUTF16("程序启动管理器"))
	procShellNotifyIcon.Call(NIM_ADD, uintptr(unsafe.Pointer(&data)))
}

func removeTray(hwnd uintptr) {
	var data notifyIconData
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.HWnd = hwnd
	data.UID = TRAY_ID
	procShellNotifyIcon.Call(NIM_DELETE, uintptr(unsafe.Pointer(&data)))
}

func hideToTray(hwnd uintptr) { procShowWindow.Call(hwnd, SW_HIDE) }
func restoreFromTray(hwnd uintptr) {
	procShowWindow.Call(hwnd, SW_RESTORE)
	procSetForeground.Call(hwnd)
}

func showTrayMenu(hwnd uintptr) {
	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	menu, _, _ := procCreatePopupMenu.Call()
	procAppendMenu.Call(menu, 0, 401, uintptr(unsafe.Pointer(utf16Ptr("显示"))))
	procAppendMenu.Call(menu, 0, 402, uintptr(unsafe.Pointer(utf16Ptr("退出"))))
	procSetForeground.Call(hwnd)
	cmd, _, _ := procTrackPopupMenu.Call(menu, 0x0100, uintptr(p.X), uintptr(p.Y), 0, hwnd, 0)
	procDestroyMenu.Call(menu)
	if cmd == 401 {
		restoreFromTray(hwnd)
	}
	if cmd == 402 {
		saveFromUI()
		removeTray(hwnd)
		procDestroyWindow.Call(hwnd)
	}
}

func button(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return control(parent, "BUTTON", text, WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_PUSHBUTTON, x, y, w, h, id)
}
func staticText(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return control(parent, "STATIC", text, WS_CHILD|WS_VISIBLE, x, y, w, h, id)
}
func checkbox(parent uintptr, x, y, w, h int32, id int) uintptr {
	return control(parent, "BUTTON", "", WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_AUTOCHECKBOX, x, y, w, h, id)
}
func checkboxText(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return control(parent, "BUTTON", text, WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_AUTOCHECKBOX, x, y, w, h, id)
}
func edit(parent uintptr, text string, x, y, w, h int32, id int, readonly bool) uintptr {
	style := uint32(WS_CHILD | WS_VISIBLE | WS_TABSTOP | WS_BORDER | ES_AUTOHSCROLL)
	if readonly {
		style |= ES_READONLY
	}
	return control(parent, "EDIT", text, style, x, y, w, h, id)
}
func control(parent uintptr, class, text string, style uint32, x, y, w, h int32, id int) uintptr {
	hwnd := createChild(0, class, text, style, x, y, w, h, parent, id)
	if app != nil && app.font != 0 {
		procSendMessage.Call(hwnd, WM_SETFONT, app.font, 1)
	}
	return hwnd
}
func createChild(exStyle uint32, class, text string, style uint32, x, y, w, h int32, parent uintptr, id int) uintptr {
	return createWindow(exStyle, utf16Ptr(class), utf16Ptr(text), style, x, y, w, h, parent, uintptr(id), getModuleHandle(), 0)
}
func createWindow(exStyle uint32, class, text *uint16, style uint32, x, y, w, h int32, parent, menu, inst, param uintptr) uintptr {
	ret, _, _ := procCreateWindowEx.Call(uintptr(exStyle), uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(text)), uintptr(style), uintptr(x), uintptr(y), uintptr(w), uintptr(h), parent, menu, inst, param)
	return ret
}
func move(hwnd uintptr, x, y, w, h int32) {
	procMoveWindow.Call(hwnd, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 0)
}
func redraw(hwnd uintptr, erase bool) {
	flags := uintptr(RDW_INVALIDATE | RDW_ALLCHILDREN)
	if erase {
		flags |= RDW_ERASE
	}
	procRedrawWindow.Call(hwnd, 0, 0, flags)
}
func setCheckIfChanged(hwnd uintptr, checked bool) {
	want := uintptr(0)
	if checked {
		want = BST_CHECKED
	}
	if procSendMessageUint(hwnd, BM_GETCHECK, 0, 0) != want {
		procSendMessage.Call(hwnd, BM_SETCHECK, want, 0)
	}
}
func setText(hwnd uintptr, text string) {
	procSetWindowText.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(text))))
}
func setTextIfChanged(hwnd uintptr, text string) {
	if getText(hwnd) == text {
		return
	}
	setText(hwnd, text)
}
func getText(hwnd uintptr) string {
	if hwnd == 0 {
		return ""
	}
	n, _, _ := procGetWindowTextLen.Call(hwnd)
	buf := make([]uint16, n+1)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), n+1)
	return syscall.UTF16ToString(buf)
}
func procSendMessageUint(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	r, _, _ := procSendMessage.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}
func message(hwnd uintptr, text, title string) {
	procMessageBox.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr(title))), 0)
}
func utf16Ptr(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }
func getModuleHandle() uintptr  { r, _, _ := procGetModuleHandle.Call(0); return r }
func loadCursor(inst uintptr, id int) uintptr {
	r, _, _ := procLoadCursor.Call(inst, uintptr(id))
	return r
}
func loadIcon(inst uintptr, id int) uintptr {
	r, _, _ := procLoadIcon.Call(inst, uintptr(id))
	return r
}
func loadAppIcon(inst uintptr) uintptr {
	icon := loadIcon(inst, APP_ICON_ID)
	if icon != 0 {
		return icon
	}
	return loadIcon(0, IDI_APPLICATION)
}
func createSolidBrush(c uint32) uintptr { r, _, _ := procCreateSolidBrush.Call(uintptr(c)); return r }
func rgb(r, g, b byte) uint32           { return uint32(r) | uint32(g)<<8 | uint32(b)<<16 }
func createFont(face string, size int32) uintptr {
	return createFontWeight(face, size, 400)
}
func createFontWeight(face string, size int32, weight int32) uintptr {
	r, _, _ := procCreateFont.Call(uintptr(-size), 0, 0, 0, uintptr(weight), 0, 0, 0, 1, 0, 0, 0, 0, uintptr(unsafe.Pointer(utf16Ptr(face))))
	return r
}
func loword(v uint32) uint16 { return uint16(v & 0xffff) }
func hiword(v uint32) uint16 { return uint16((v >> 16) & 0xffff) }

func isStatusControl(hwnd uintptr) bool {
	if app == nil {
		return false
	}
	for _, r := range app.rows {
		if r.status == hwnd {
			return true
		}
	}
	return false
}

func isWindowStateSavable(hwnd uintptr) bool {
	visible, _, _ := procIsWindowVisible.Call(hwnd)
	iconic, _, _ := procIsIconic.Call(hwnd)
	return visible != 0 && iconic == 0
}

func applyMinWindowSize(hwnd uintptr, lparam uintptr) {
	if lparam == 0 {
		return
	}
	minW := monitorWidth(hwnd) / 2
	writeInt32(lparam+24, minW)
	writeInt32(lparam+28, 360)
}

func startupMinWindowWidth() int32 {
	if w := getSystemMetric(SM_CXSCREEN); w > 0 {
		return w / 2
	}
	return 640
}

func writeInt32(addr uintptr, value int32) {
	procRtlMoveMemory.Call(addr, uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(value))
}

func monitorWidth(hwnd uintptr) int32 {
	monitor, _, _ := procMonitorFromWindow.Call(hwnd, MONITOR_DEFAULTTONEAREST)
	if monitor != 0 {
		mi := monitorInfo{CbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
		ok, _, _ := procGetMonitorInfo.Call(monitor, uintptr(unsafe.Pointer(&mi)))
		if ok != 0 && mi.RcMonitor.Right > mi.RcMonitor.Left {
			return mi.RcMonitor.Right - mi.RcMonitor.Left
		}
	}
	if w := getSystemMetric(SM_CXSCREEN); w > 0 {
		return w
	}
	return 1280
}

func isWindowRectVisible(x, y, w, h int32) bool {
	if w <= 0 || h <= 0 {
		return false
	}
	vx := getSystemMetric(SM_XVIRTUAL)
	vy := getSystemMetric(SM_YVIRTUAL)
	vw := getSystemMetric(SM_CXVIRTUAL)
	vh := getSystemMetric(SM_CYVIRTUAL)
	if vw <= 0 || vh <= 0 {
		return x >= 0 && y >= 0
	}
	return x < vx+vw && x+w > vx && y < vy+vh && y+h > vy
}

func centeredDialogPos(parent uintptr, w, h int32) (int32, int32) {
	var wr rect
	if parent != 0 {
		procGetWindowRect.Call(parent, uintptr(unsafe.Pointer(&wr)))
	}
	if wr.Right <= wr.Left || wr.Bottom <= wr.Top {
		return 220, 180
	}
	x := wr.Left + ((wr.Right-wr.Left)-w)/2
	y := wr.Top + ((wr.Bottom-wr.Top)-h)/2
	if !isWindowRectVisible(x, y, w, h) {
		return 220, 180
	}
	return x, y
}

func getSystemMetric(index int32) int32 {
	r, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r)
}

func _unused() {
	_ = filepath.Separator
	_ = strconv.IntSize
	_ = strings.Builder{}
	_ = procGetDlgItem
	_ = procGetParent
}
