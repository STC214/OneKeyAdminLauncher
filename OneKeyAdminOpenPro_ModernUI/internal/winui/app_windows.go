package winui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"program-launch-manager/internal/config"
	"program-launch-manager/internal/process"
)

const (
	appName = "程序启动管理器"

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
	BS_OWNERDRAW        = 0x0000000B
	BS_FLAT             = 0x00008000
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
	WM_PAINT           = 0x000F
	WM_CLOSE           = 0x0010
	WM_CANCELMODE      = 0x001F
	WM_GETMINMAXINFO   = 0x0024
	WM_DRAWITEM        = 0x002B
	WM_ERASEBKGND      = 0x0014
	WM_COMMAND         = 0x0111
	WM_TIMER           = 0x0113
	WM_VSCROLL         = 0x0115
	WM_CTLCOLORSTATIC  = 0x0138
	WM_CTLCOLORLISTBOX = 0x0134
	WM_CTLCOLOREDIT    = 0x0133
	WM_CTLCOLORBTN     = 0x0135
	WM_MOUSEMOVE       = 0x0200
	WM_LBUTTONDOWN     = 0x0201
	WM_LBUTTONUP       = 0x0202
	WM_MOUSEWHEEL      = 0x020A
	WM_CAPTURECHANGED  = 0x0215
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

	SIZE_MINIMIZED             = 1
	COLOR_WINDOW               = 5
	SB_VERT                    = 1
	SB_LINEUP                  = 0
	SB_LINEDOWN                = 1
	SB_PAGEUP                  = 2
	SB_PAGEDOWN                = 3
	SB_THUMBPOSITION           = 4
	SB_THUMBTRACK              = 5
	SB_TOP                     = 6
	SB_BOTTOM                  = 7
	SIF_RANGE                  = 0x0001
	SIF_PAGE                   = 0x0002
	SIF_POS                    = 0x0004
	SIF_ALL                    = SIF_RANGE | SIF_PAGE | SIF_POS
	WHEEL_DELTA                = 120
	RDW_INVALIDATE             = 0x0001
	RDW_ERASE                  = 0x0004
	RDW_ALLCHILDREN            = 0x0080
	ODS_SELECTED               = 0x0001
	TRANSPARENT                = 1
	DT_CENTER                  = 0x0001
	DT_VCENTER                 = 0x0004
	DT_SINGLELINE              = 0x0020
	IDC_ARROW                  = 32512
	IDI_APPLICATION            = 32512
	ICON_SMALL                 = 0
	ICON_BIG                   = 1
	APP_ICON_ID                = 101
	SM_XVIRTUAL                = 76
	SM_YVIRTUAL                = 77
	SM_CXSCREEN                = 0
	SM_CXVIRTUAL               = 78
	SM_CYVIRTUAL               = 79
	MONITOR_DEFAULTTONEAREST   = 2
	ERROR_CLASS_ALREADY_EXISTS = 1410

	WM_TRAY_ICON      = WM_APP + 10
	WM_STATUS_READY   = WM_APP + 11
	WM_AUTOSAVE_ERR   = WM_APP + 13
	WM_OPERATION_DONE = WM_APP + 14
	ID_TIMER          = 2001
	TRAY_ID           = 1

	DWMWA_USE_IMMERSIVE_DARK_MODE_OLD = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE     = 20
	DWMWA_CAPTION_COLOR               = 35
	DWMWA_TEXT_COLOR                  = 36

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

	MB_ABORTRETRYIGNORE = 0x00000002
	IDRETRY             = 4
	IDIGNORE            = 5

	rowBase           = 1000
	rowStep           = 10
	rowStartY         = 108
	rowHeight         = 58
	rowBottomPadding  = 20
	modernMinWidth    = 880
	scrollbarWidth    = 8
	scrollbarMargin   = 10
	scrollbarMinThumb = 42
)

// appVersion is replaced by the build script through -ldflags -X.
var appVersion = "dev"

func appWindowTitle() string { return appName + " " + appVersion }

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	dwmapi   = syscall.NewLazyDLL("dwmapi.dll")

	procDefWindowProc     = user32.NewProc("DefWindowProcW")
	procRegisterClassEx   = user32.NewProc("RegisterClassExW")
	procCreateWindowEx    = user32.NewProc("CreateWindowExW")
	procShowWindow        = user32.NewProc("ShowWindow")
	procUpdateWindow      = user32.NewProc("UpdateWindow")
	procGetMessage        = user32.NewProc("GetMessageW")
	procRegisterWindowMsg = user32.NewProc("RegisterWindowMessageW")
	procTranslateMessage  = user32.NewProc("TranslateMessage")
	procDispatchMessage   = user32.NewProc("DispatchMessageW")
	procPostQuitMessage   = user32.NewProc("PostQuitMessage")
	procPostMessage       = user32.NewProc("PostMessageW")
	procSendMessage       = user32.NewProc("SendMessageW")
	procBeginPaint        = user32.NewProc("BeginPaint")
	procEndPaint          = user32.NewProc("EndPaint")
	procInvalidateRect    = user32.NewProc("InvalidateRect")
	procSetWindowText     = user32.NewProc("SetWindowTextW")
	procGetWindowText     = user32.NewProc("GetWindowTextW")
	procGetWindowTextLen  = user32.NewProc("GetWindowTextLengthW")
	procMoveWindow        = user32.NewProc("MoveWindow")
	procRedrawWindow      = user32.NewProc("RedrawWindow")
	procDrawText          = user32.NewProc("DrawTextW")
	procFrameRect         = user32.NewProc("FrameRect")
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
	procSetCapture        = user32.NewProc("SetCapture")
	procReleaseCapture    = user32.NewProc("ReleaseCapture")

	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	procRtlMoveMemory    = kernel32.NewProc("RtlMoveMemory")
	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procSetTextColor     = gdi32.NewProc("SetTextColor")
	procSetBkColor       = gdi32.NewProc("SetBkColor")
	procSetBkMode        = gdi32.NewProc("SetBkMode")
	procCreateFont       = gdi32.NewProc("CreateFontW")
	procShellNotifyIcon  = shell32.NewProc("Shell_NotifyIconW")
	procGetOpenFileName  = comdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtError  = comdlg32.NewProc("CommDlgExtendedError")
	procDwmSetWindowAttr = dwmapi.NewProc("DwmSetWindowAttribute")
)

var taskbarCreatedMsg uint32

type point struct{ X, Y int32 }
type rect struct{ Left, Top, Right, Bottom int32 }
type monitorInfo struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
}
type drawItemStruct struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   uintptr
	HDC        uintptr
	RcItem     rect
	ItemData   uintptr
}
type paintStruct struct {
	HDC         uintptr
	Erase       int32
	Paint       rect
	Restore     int32
	IncUpdate   int32
	RGBReserved [32]byte
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
	uwpLabel      uintptr
	status        uintptr
	selectProcess uintptr
	remove        uintptr
	index         int
	visible       bool
}

type appState struct {
	hwnd             uintptr
	font             uintptr
	icon             uintptr
	bgBrush          uintptr
	editBrush        uintptr
	panelBrush       uintptr
	cfg              config.File
	rows             []rowControls
	scrollY          int32
	scrollDrag       bool
	scrollDragOffset int32
	mu               sync.Mutex
	saveQueue        chan saveRequest

	statusMu         sync.Mutex
	statusBusy       bool
	runningSnapshot  map[string]bool
	updatingUI       bool
	autoSaveDirty    bool
	autoSaveBusy     bool
	autoSaveErr      string
	autoSaveErrSeen  bool
	nextSaveSeq      uint64
	lastSaveSeq      uint64
	controlCreateErr error
	trayAvailable    bool
	timerActive      bool

	operationMu      sync.Mutex
	operationResults []operationResult
	launchBatches    int
	launchFailures   []string
}

type operationResult struct {
	heading  string
	failures []string
}

type saveRequest struct {
	seq  uint64
	cfg  config.File
	done chan error
}

var app *appState

func Run() int {
	cfg, err := config.Load()
	if err != nil {
		message(0, "配置加载失败，原文件未被修改。\n\n"+err.Error(), "程序启动管理器")
		return 1
	}
	app = &appState{cfg: cfg, saveQueue: make(chan saveRequest, 8)}
	go runSaveWorker(app.saveQueue)
	if app.cfg.Window.W == 0 {
		app.cfg.Window = config.DefaultWindow()
	}
	return runWindow()
}

func runWindow() int {
	if r, _, _ := procRegisterWindowMsg.Call(uintptr(unsafe.Pointer(utf16Ptr("TaskbarCreated")))); r != 0 {
		taskbarCreatedMsg = uint32(r)
	}
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
	if err := registerWindowClass(&wc); err != nil {
		message(0, "注册主窗口类失败：\n\n"+err.Error(), "程序启动管理器")
		return 1
	}

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
	hwnd, err := createWindowChecked(0, className, utf16Ptr(appWindowTitle()), WS_OVERLAPPEDWINDOW, x, y, win.W, win.H, 0, 0, hinst, 0)
	if err != nil {
		message(0, "创建主窗口失败：\n\n"+err.Error(), "程序启动管理器")
		return 1
	}
	applyDarkTitleBar(hwnd)
	app.hwnd = hwnd
	procShowWindow.Call(hwnd, SW_SHOWNORMAL)
	procUpdateWindow.Call(hwnd)
	var m msg
	for {
		status, err := nextMessage(&m)
		if err != nil {
			message(hwnd, "读取窗口消息失败：\n\n"+err.Error(), "程序启动管理器")
			return 1
		}
		if status == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
	return int(m.WParam)
}

func wndProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	if taskbarCreatedMsg != 0 && msg == taskbarCreatedMsg {
		app.trayAvailable = addTray(hwnd)
		if !app.trayAvailable {
			restoreFromTray(hwnd)
		}
		return 0
	}
	switch msg {
	case WM_CREATE:
		app.hwnd = hwnd
		app.bgBrush = createSolidBrush(rgb(18, 20, 24))
		app.editBrush = createSolidBrush(rgb(30, 34, 40))
		app.panelBrush = createSolidBrush(rgb(38, 43, 52))
		app.font = createFont("Segoe UI", 15)
		app.icon = loadAppIcon(getModuleHandle())
		procSendMessage.Call(hwnd, WM_SETICON, ICON_BIG, app.icon)
		procSendMessage.Call(hwnd, WM_SETICON, ICON_SMALL, app.icon)
		createToolbar(hwnd)
		rebuildRows(hwnd)
		if app.controlCreateErr != nil {
			return ^uintptr(0)
		}
		layout(hwnd, true)
		app.trayAvailable = addTray(hwnd)
		if !app.trayAvailable {
			message(hwnd, "系统托盘图标创建失败；窗口将保持普通最小化，不会隐藏到托盘。", "程序启动管理器")
		}
		if timer, _, _ := procSetTimer.Call(hwnd, ID_TIMER, 1000, 0); timer != 0 {
			app.timerActive = true
		} else {
			message(hwnd, "创建状态刷新定时器失败；状态不会自动刷新，配置仍会在手动保存或退出时保存。", "程序启动管理器")
		}
		requestStatusRefresh()
		return 0
	case WM_SIZE:
		if wparam == SIZE_MINIMIZED {
			if app.trayAvailable {
				hideToTray(hwnd)
				return 0
			}
			break
		}
		layout(hwnd, true)
		return 0
	case WM_PAINT:
		var ps paintStruct
		hdc, _, _ := procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		drawCustomScrollBar(hwnd, hdc)
		procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
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
	case WM_LBUTTONDOWN:
		if beginCustomScrollDrag(hwnd, lparam) {
			return 0
		}
	case WM_MOUSEMOVE:
		if app != nil && app.scrollDrag {
			handleCustomScrollDrag(hwnd, lparam)
			return 0
		}
	case WM_LBUTTONUP:
		if app != nil && app.scrollDrag {
			app.scrollDrag = false
			procReleaseCapture.Call()
			return 0
		}
	case WM_CANCELMODE, WM_CAPTURECHANGED:
		if app != nil {
			app.scrollDrag = false
		}
	case WM_TIMER:
		requestStatusRefresh()
		flushAutoSave()
		return 0
	case WM_STATUS_READY:
		applyStatusSnapshot()
		return 0
	case WM_AUTOSAVE_ERR:
		app.mu.Lock()
		errText := app.autoSaveErr
		app.mu.Unlock()
		if errText != "" {
			message(hwnd, "自动保存失败，将在后续修改后重试。\n\n"+errText, "程序启动管理器")
		}
		return 0
	case WM_OPERATION_DONE:
		showNextOperationResult(hwnd)
		return 0
	case WM_COMMAND:
		if app != nil && app.updatingUI {
			return 0
		}
		handleCommand(hwnd, loword(uint32(wparam)), hiword(uint32(wparam)), lparam)
		return 0
	case WM_DRAWITEM:
		if drawOwnerButton(lparam) {
			return 1
		}
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
		color := rgb(232, 237, 243)
		back := rgb(30, 34, 40)
		if isStatusControl(lparam) {
			txt := getText(lparam)
			if txt != "" && txt != "未运行" && txt != "未绑定" {
				color = rgb(86, 225, 128)
			} else {
				color = rgb(148, 163, 184)
			}
		}
		if msg == WM_CTLCOLORBTN {
			color = rgb(226, 232, 240)
			back = rgb(30, 34, 40)
		}
		procSetTextColor.Call(wparam, uintptr(color))
		procSetBkColor.Call(wparam, uintptr(back))
		if msg == WM_CTLCOLORBTN {
			return app.editBrush
		}
		return app.editBrush
	case WM_CLOSE:
		if !confirmClose(hwnd) {
			return 0
		}
		closeActiveProcessDialog()
		removeTray(hwnd)
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		if app.timerActive {
			procKillTimer.Call(hwnd, ID_TIMER)
		}
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
	y := int32(46)
	button(hwnd, "一键开启", 18, y, 92, 32, ID_START)
	button(hwnd, "一键关闭", 122, y, 92, 32, ID_STOP)
	button(hwnd, "添加程序", 226, y, 92, 32, ID_ADD)
	button(hwnd, "全选", 330, y, 60, 32, ID_ALL)
	button(hwnd, "全不选", 402, y, 76, 32, ID_NONE)
	button(hwnd, "保存配置", 490, y, 92, 32, ID_SAVE)
}

func rebuildRows(hwnd uintptr) {
	for _, r := range app.rows {
		for _, h := range []uintptr{r.check, r.path, r.browse, r.uwp, r.uwpLabel, r.status, r.selectProcess, r.remove} {
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
		rowY := y + 2
		move(r.check, 18, rowY+9, 18, 18)
		right := width - 32
		removeW := int32(52)
		selectW := int32(96)
		statusW := int32(154)
		uwpW := int32(58)
		browseW := int32(82)
		gap := int32(8)
		xRemove := right - removeW
		xSelect := xRemove - gap - selectW
		xStatus := xSelect - gap - statusW
		xUWP := xStatus - gap - uwpW
		xBrowse := xUWP - gap - browseW
		pathX := int32(44)
		pathW := xBrowse - gap - pathX
		if pathW < 280 {
			pathW = 280
		}
		move(r.path, pathX, rowY, pathW, 32)
		move(r.browse, xBrowse, rowY, browseW, 32)
		move(r.uwp, xUWP, rowY+9, 18, 18)
		move(r.uwpLabel, xUWP+22, rowY+6, uwpW-22, 22)
		move(r.status, xStatus, rowY, statusW, 32)
		move(r.selectProcess, xSelect, rowY, selectW, 32)
		move(r.remove, xRemove, rowY, removeW, 32)
		setRowVisible(r, visible)
	}
	updateScrollBar(hwnd, rc.Bottom-rc.Top, erase)
	if erase {
		redraw(hwnd, true)
	}
}

func ensureRowPool(hwnd uintptr, count int) {
	for len(app.rows) < count {
		slot := len(app.rows)
		id := rowBase + slot*rowStep
		app.rows = append(app.rows, rowControls{
			check:         checkbox(hwnd, 0, 0, 18, 18, id+1),
			path:          edit(hwnd, "", 0, 0, 260, 30, id+2, false),
			browse:        button(hwnd, "浏览", 0, 0, 78, 30, id+3),
			uwp:           checkbox(hwnd, 0, 0, 18, 18, id+7),
			uwpLabel:      staticText(hwnd, "UWP", 0, 0, 34, 22, 0),
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
		updateScrollBar(hwnd, rc.Bottom-rc.Top, false)
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

func updateScrollBar(hwnd uintptr, clientH int32, erase bool) {
	invalidateCustomScrollBar(hwnd, erase)
}

func drawCustomScrollBar(hwnd uintptr, hdc uintptr) {
	if hdc == 0 {
		return
	}
	clearCustomScrollBarArea(hwnd, hdc)
	track, thumb, ok := customScrollRects(hwnd)
	if !ok {
		return
	}
	trackBrush := createSolidBrush(rgb(24, 28, 34))
	thumbBrush := createSolidBrush(rgb(96, 110, 132))
	defer procDeleteObject.Call(trackBrush)
	defer procDeleteObject.Call(thumbBrush)
	procFillRect.Call(hdc, uintptr(unsafe.Pointer(&track)), trackBrush)
	procFillRect.Call(hdc, uintptr(unsafe.Pointer(&thumb)), thumbBrush)
}

func clearCustomScrollBarArea(hwnd uintptr, hdc uintptr) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	clear := rect{
		Left:   rc.Right - scrollbarMargin - scrollbarWidth - 12,
		Top:    rowStartY,
		Right:  rc.Right,
		Bottom: rc.Bottom,
	}
	procFillRect.Call(hdc, uintptr(unsafe.Pointer(&clear)), app.bgBrush)
}

func beginCustomScrollDrag(hwnd uintptr, lparam uintptr) bool {
	track, thumb, ok := customScrollRects(hwnd)
	if !ok {
		return false
	}
	x, y := pointFromLParam(lparam)
	if x < track.Left || x >= track.Right || y < track.Top || y >= track.Bottom {
		return false
	}
	if y >= thumb.Top && y < thumb.Bottom {
		app.scrollDrag = true
		app.scrollDragOffset = y - thumb.Top
		procSetCapture.Call(hwnd)
		return true
	}
	if y < thumb.Top {
		setScroll(hwnd, app.scrollY-scrollPage(listViewportHeight(currentClientHeight(hwnd))))
	} else {
		setScroll(hwnd, app.scrollY+scrollPage(listViewportHeight(currentClientHeight(hwnd))))
	}
	return true
}

func handleCustomScrollDrag(hwnd uintptr, lparam uintptr) {
	track, thumb, ok := customScrollRects(hwnd)
	if !ok {
		return
	}
	_, y := pointFromLParam(lparam)
	trackTravel := (track.Bottom - track.Top) - (thumb.Bottom - thumb.Top)
	if trackTravel <= 0 {
		return
	}
	thumbTop := y - app.scrollDragOffset
	if thumbTop < track.Top {
		thumbTop = track.Top
	}
	maxTop := track.Bottom - (thumb.Bottom - thumb.Top)
	if thumbTop > maxTop {
		thumbTop = maxTop
	}
	max := maxScroll(listViewportHeight(currentClientHeight(hwnd)))
	setScroll(hwnd, (thumbTop-track.Top)*max/trackTravel)
}

func customScrollRects(hwnd uintptr) (rect, rect, bool) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	clientH := rc.Bottom - rc.Top
	listH := listViewportHeight(clientH)
	max := maxScroll(listH)
	if max <= 0 || contentHeight() <= listH {
		return rect{}, rect{}, false
	}
	track := rect{
		Left:   rc.Right - scrollbarMargin - scrollbarWidth,
		Top:    rowStartY + 6,
		Right:  rc.Right - scrollbarMargin,
		Bottom: rc.Bottom - 12,
	}
	if track.Bottom <= track.Top {
		return rect{}, rect{}, false
	}
	trackH := track.Bottom - track.Top
	thumbH := trackH * listH / contentHeight()
	if thumbH < scrollbarMinThumb {
		thumbH = scrollbarMinThumb
	}
	if thumbH > trackH {
		thumbH = trackH
	}
	travel := trackH - thumbH
	thumbTop := track.Top
	if travel > 0 {
		thumbTop += app.scrollY * travel / max
	}
	thumb := rect{
		Left:   track.Left,
		Top:    thumbTop,
		Right:  track.Right,
		Bottom: thumbTop + thumbH,
	}
	return track, thumb, true
}

func invalidateCustomScrollBar(hwnd uintptr, erase bool) {
	var rc rect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	dirty := rect{
		Left:   rc.Right - scrollbarMargin - scrollbarWidth - 12,
		Top:    rowStartY,
		Right:  rc.Right,
		Bottom: rc.Bottom,
	}
	eraseFlag := uintptr(0)
	if erase {
		eraseFlag = 1
	}
	procInvalidateRect.Call(hwnd, uintptr(unsafe.Pointer(&dirty)), eraseFlag)
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
	for _, h := range []uintptr{r.check, r.path, r.browse, r.uwp, r.uwpLabel, r.status, r.selectProcess, r.remove} {
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
		launchProgramsAsync(hwnd, cloneConfig(app.cfg).Programs)
		needsRefresh = true
	case ID_STOP:
		syncFromUI()
		markAutoSaveDirty()
		closeProgramsAsync(hwnd, cloneConfig(app.cfg).Programs)
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
		if err := saveConfig(cloneConfig(app.cfg)); err != nil {
			message(hwnd, "配置保存失败：\n\n"+err.Error(), "程序启动管理器")
		} else {
			message(hwnd, "配置已保存。", "程序启动管理器")
		}
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

func saveFromUI() error {
	syncFromUI()
	return saveConfig(cloneConfig(app.cfg))
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
	seq, done := enqueueSave(cfg)
	err := <-done
	app.recordSaveResult(seq, err, false)
	return err
}

func enqueueSave(cfg config.File) (uint64, <-chan error) {
	app.mu.Lock()
	app.nextSaveSeq++
	seq := app.nextSaveSeq
	app.mu.Unlock()
	done := make(chan error, 1)
	app.saveQueue <- saveRequest{seq: seq, cfg: cfg, done: done}
	return seq, done
}

func runSaveWorker(queue <-chan saveRequest) {
	for req := range queue {
		req.done <- config.Save(req.cfg)
	}
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

	seq, done := enqueueSave(cfg)
	hwnd := app.hwnd
	go func() {
		err := <-done
		_, notify := app.recordSaveResult(seq, err, true)
		app.mu.Lock()
		app.autoSaveBusy = false
		app.mu.Unlock()
		if notify {
			procPostMessage.Call(hwnd, WM_AUTOSAVE_ERR, 0, 0)
		}
	}()
}

func (a *appState) recordSaveResult(seq uint64, err error, auto bool) (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if seq <= a.lastSaveSeq {
		return a.autoSaveErr, false
	}
	a.lastSaveSeq = seq
	if err == nil {
		a.autoSaveErr = ""
		a.autoSaveErrSeen = false
		return "", false
	}
	if !auto {
		return a.autoSaveErr, false
	}
	a.autoSaveErr = err.Error()
	if a.autoSaveErrSeen {
		return a.autoSaveErr, false
	}
	a.autoSaveErrSeen = true
	return a.autoSaveErr, true
}

func launchProgramsAsync(hwnd uintptr, items []config.ProgramItem) {
	app.beginLaunchBatch()
	go func() {
		failures := launchEnabledPrograms(items, process.Launch)
		failures, completed := app.finishLaunchBatch(failures)
		if !completed || len(failures) == 0 {
			return
		}
		app.queueOperationResult("以下程序启动失败：", failures)
		procPostMessage.Call(hwnd, WM_OPERATION_DONE, 0, 0)
	}()
}

const maxConcurrentLaunches = 6

var launchSlots = make(chan struct{}, maxConcurrentLaunches)

func launchEnabledPrograms(items []config.ProgramItem, launch func(string, bool) error) []string {
	var enabled []config.ProgramItem
	for _, item := range items {
		if item.Enabled {
			enabled = append(enabled, item)
		}
	}
	if len(enabled) == 0 {
		return nil
	}

	errs := make([]error, len(enabled))
	jobs := make(chan int)
	workerCount := min(maxConcurrentLaunches, len(enabled))
	var workers sync.WaitGroup
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for i := range jobs {
				func() {
					launchSlots <- struct{}{}
					defer func() { <-launchSlots }()
					errs[i] = launch(enabled[i].Path, enabled[i].IsUWP)
				}()
			}
		}()
	}
	for i := range enabled {
		jobs <- i
	}
	close(jobs)
	workers.Wait()

	var failures []string
	for i, err := range errs {
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", enabled[i].Path, err))
		}
	}
	return failures
}

func closeProgramsAsync(hwnd uintptr, items []config.ProgramItem) {
	names := uniqueEnabledProcessNames(items)
	go func() {
		var failures []string
		if err := process.CloseByNames(names); err != nil {
			failures = append(failures, err.Error())
		}
		if len(failures) == 0 {
			return
		}
		app.queueOperationResult("以下程序关闭失败：", failures)
		procPostMessage.Call(hwnd, WM_OPERATION_DONE, 0, 0)
	}()
}

func uniqueEnabledProcessNames(items []config.ProgramItem) []string {
	seen := make(map[string]bool)
	var names []string
	for _, item := range items {
		if !item.Enabled {
			continue
		}
		name := strings.TrimSpace(effectiveProcess(item))
		key := strings.ToLower(name)
		if name != "" && !seen[key] {
			seen[key] = true
			names = append(names, name)
		}
	}
	return names
}

func showNextOperationResult(hwnd uintptr) {
	result, ok := app.popOperationResult()
	if !ok {
		return
	}
	message(hwnd, result.heading+"\n\n"+strings.Join(result.failures, "\n"), "程序启动管理器")
}

func (a *appState) queueOperationResult(heading string, failures []string) {
	a.operationMu.Lock()
	a.operationResults = append(a.operationResults, operationResult{heading: heading, failures: append([]string(nil), failures...)})
	a.operationMu.Unlock()
}

func (a *appState) beginLaunchBatch() {
	a.operationMu.Lock()
	a.launchBatches++
	a.operationMu.Unlock()
}

func (a *appState) finishLaunchBatch(failures []string) ([]string, bool) {
	a.operationMu.Lock()
	defer a.operationMu.Unlock()
	a.launchFailures = append(a.launchFailures, failures...)
	a.launchBatches--
	if a.launchBatches > 0 {
		return nil, false
	}
	all := append([]string(nil), a.launchFailures...)
	a.launchFailures = nil
	return all, true
}

func (a *appState) popOperationResult() (operationResult, bool) {
	a.operationMu.Lock()
	defer a.operationMu.Unlock()
	if len(a.operationResults) == 0 {
		return operationResult{}, false
	}
	result := a.operationResults[0]
	a.operationResults = a.operationResults[1:]
	return result, true
}

func confirmClose(hwnd uintptr) bool {
	for {
		if err := saveFromUI(); err != nil {
			choice := messageChoice(hwnd, "配置保存失败。\n\n"+err.Error()+"\n\n重试：再次保存\n中止：返回程序\n忽略：放弃修改并退出", "程序启动管理器", MB_ABORTRETRYIGNORE)
			switch choice {
			case IDRETRY:
				continue
			case IDIGNORE:
				return true
			default:
				return false
			}
		}
		return true
	}
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
		updateProgramPath(item, getText(r.path))
		item.Enabled = procSendMessageUint(r.check, BM_GETCHECK, 0, 0) == BST_CHECKED
		item.IsUWP = procSendMessageUint(r.uwp, BM_GETCHECK, 0, 0) == BST_CHECKED || strings.HasPrefix(strings.ToLower(item.Path), "shell:")
	}
	if isWindowStateSavable(app.hwnd) {
		var wr rect
		procGetWindowRect.Call(app.hwnd, uintptr(unsafe.Pointer(&wr)))
		if wr.Right > wr.Left && wr.Bottom > wr.Top {
			app.cfg.Window = config.WindowState{X: wr.Left, Y: wr.Top, W: wr.Right - wr.Left, H: wr.Bottom - wr.Top}
		}
	}
}

func updateProgramPath(item *config.ProgramItem, path string) {
	changed := !strings.EqualFold(strings.TrimSpace(item.Path), strings.TrimSpace(path))
	item.Path = path
	if changed && item.SelectedProcess == "" {
		item.ProcessName = process.ProcessNameForPath(path)
	}
}

func requestStatusRefresh() {
	cfg := cloneConfig(app.cfg)

	app.statusMu.Lock()
	if app.statusBusy {
		app.statusMu.Unlock()
		return
	}
	app.statusBusy = true
	app.statusMu.Unlock()

	hwnd := app.hwnd
	go func() {
		runningNames := map[string]bool{}
		procs, err := process.List()
		if err != nil {
			app.statusMu.Lock()
			app.statusBusy = false
			app.statusMu.Unlock()
			return
		}
		for _, p := range procs {
			runningNames[strings.ToLower(p.Name)] = true
		}
		status := configuredProcessStatus(cfg, runningNames)

		app.statusMu.Lock()
		changed := !sameBoolMap(app.runningSnapshot, status)
		app.runningSnapshot = status
		app.statusBusy = false
		app.statusMu.Unlock()

		if changed {
			procPostMessage.Call(hwnd, WM_STATUS_READY, 0, 0)
		}
	}()
}

func configuredProcessStatus(cfg config.File, runningNames map[string]bool) map[string]bool {
	status := map[string]bool{}
	for _, item := range cfg.Programs {
		name := strings.ToLower(strings.TrimSpace(effectiveProcess(item)))
		if name == "" {
			continue
		}
		status[name] = runningNames[name]
	}
	return status
}

func sameBoolMap(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
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
	filterData := dialogFilter("程序 (*.exe;*.lnk)", "*.exe;*.lnk", "所有文件 (*.*)", "*.*")
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
		code, _, _ := procCommDlgExtError.Call()
		if code != 0 {
			message(hwnd, fmt.Sprintf("打开文件选择器失败，错误码：0x%08X", uint32(code)), "程序启动管理器")
		}
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func dialogFilter(parts ...string) []uint16 {
	return utf16.Encode([]rune(strings.Join(parts, "\x00") + "\x00\x00"))
}

func chooseProcess(parent uintptr) string {
	procs, err := process.List()
	if err != nil {
		message(parent, "枚举进程失败：\n\n"+err.Error(), "程序启动管理器")
		return ""
	}
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
	if err := registerWindowClass(&wc); err != nil {
		message(d.parent, "注册进程选择窗口失败：\n\n"+err.Error(), "程序启动管理器")
		return ""
	}
	x, y := centeredDialogPos(d.parent, 620, 680)
	app.controlCreateErr = nil
	var err error
	d.hwnd, err = createWindowChecked(WS_EX_TOOLWINDOW, className, utf16Ptr("选择进程"), WS_CAPTION|WS_SYSMENU, x, y, 620, 680, d.parent, 0, hinst, 0)
	if err != nil {
		message(d.parent, "创建进程选择窗口失败：\n\n"+err.Error(), "程序启动管理器")
		return ""
	}
	applyDarkTitleBar(d.hwnd)
	procShowWindow.Call(d.hwnd, SW_SHOWNORMAL)
	procUpdateWindow.Call(d.hwnd)
	procEnableWindow.Call(d.parent, 0)
	var m msg
	for !d.done {
		status, err := nextMessage(&m)
		if err != nil {
			message(d.hwnd, "读取窗口消息失败：\n\n"+err.Error(), "程序启动管理器")
			d.done = true
			procDestroyWindow.Call(d.hwnd)
			break
		}
		if status == 0 {
			d.done = true
			if d.hwnd != 0 {
				procDestroyWindow.Call(d.hwnd)
			}
			procPostQuitMessage.Call(m.WParam)
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
		dlg.titleFont = createFontWeight("Segoe UI", 23, 600)
		dlg.smallFont = createFontWeight("Segoe UI", 13, 400)
		dlg.listFont = createFontWeight("Segoe UI", 15, 400)
		dlg.title = staticText(hwnd, "选择进程", 26, 24, 260, 32, 300)
		dlg.subtitle = staticText(hwnd, "绑定后用于运行状态检测和一键关闭", 26, 60, 420, 24, 304)
		dlg.countText = staticText(hwnd, fmt.Sprintf("当前进程  %d", len(dlg.procs)), 468, 62, 126, 22, 305)
		procSendMessage.Call(dlg.title, WM_SETFONT, dlg.titleFont, 1)
		procSendMessage.Call(dlg.subtitle, WM_SETFONT, dlg.smallFont, 1)
		procSendMessage.Call(dlg.countText, WM_SETFONT, dlg.smallFont, 1)
		dlg.list = createChild(WS_EX_CLIENTEDGE, "LISTBOX", "", WS_CHILD|WS_VISIBLE|WS_TABSTOP|WS_VSCROLL|WS_BORDER|LBS_NOTIFY, 26, 108, 568, 478, hwnd, 301)
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
		button(hwnd, "确定", 404, 606, 90, 36, 302)
		button(hwnd, "取消", 504, 606, 90, 36, 303)
		if app.controlCreateErr != nil {
			return ^uintptr(0)
		}
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
	case WM_DRAWITEM:
		if drawOwnerButton(lparam) {
			return 1
		}
	case WM_CLOSE:
		dlg.done = true
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		dlg.done = true
		dlg.hwnd = 0
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
		color := rgb(232, 237, 243)
		back := rgb(18, 20, 24)
		if lparam == dlg.title {
			color = rgb(248, 250, 252)
		}
		if lparam == dlg.subtitle || lparam == dlg.countText {
			color = rgb(148, 163, 184)
		}
		if msg == WM_CTLCOLORLISTBOX {
			procSetBkColor.Call(wparam, uintptr(rgb(30, 34, 40)))
			return app.panelBrush
		}
		if msg == WM_CTLCOLORBTN {
			color = rgb(226, 232, 240)
		}
		procSetTextColor.Call(wparam, uintptr(color))
		procSetBkColor.Call(wparam, uintptr(back))
		return app.bgBrush
	}
	r, _, _ := procDefWindowProc.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}

func closeActiveProcessDialog() {
	if dlg == nil || dlg.done || dlg.hwnd == 0 {
		return
	}
	dlg.done = true
	procDestroyWindow.Call(dlg.hwnd)
}

func addTray(hwnd uintptr) bool {
	var data notifyIconData
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.HWnd = hwnd
	data.UID = TRAY_ID
	data.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	data.UCallbackMessage = WM_TRAY_ICON
	data.HIcon = app.icon
	copy(data.SzTip[:], syscall.StringToUTF16(appWindowTitle()))
	r, _, _ := procShellNotifyIcon.Call(NIM_ADD, uintptr(unsafe.Pointer(&data)))
	return r != 0
}

func removeTray(hwnd uintptr) {
	if !app.trayAvailable {
		return
	}
	var data notifyIconData
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.HWnd = hwnd
	data.UID = TRAY_ID
	procShellNotifyIcon.Call(NIM_DELETE, uintptr(unsafe.Pointer(&data)))
	app.trayAvailable = false
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
		procSendMessage.Call(hwnd, WM_CLOSE, 0, 0)
	}
}

func button(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return control(parent, "BUTTON", text, WS_CHILD|WS_VISIBLE|WS_TABSTOP|BS_OWNERDRAW|BS_FLAT, x, y, w, h, id)
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
	ret, err := createWindowChecked(exStyle, class, text, style, x, y, w, h, parent, menu, inst, param)
	if err != nil && app != nil && app.controlCreateErr == nil {
		app.controlCreateErr = err
	}
	return ret
}

func createWindowChecked(exStyle uint32, class, text *uint16, style uint32, x, y, w, h int32, parent, menu, inst, param uintptr) (uintptr, error) {
	ret, _, callErr := procCreateWindowEx.Call(uintptr(exStyle), uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(text)), uintptr(style), uintptr(x), uintptr(y), uintptr(w), uintptr(h), parent, menu, inst, param)
	if ret == 0 {
		return 0, winUIError(callErr)
	}
	return ret, nil
}

func registerWindowClass(wc *wndclassex) error {
	r, _, callErr := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wc)))
	if r != 0 || callErr == syscall.Errno(ERROR_CLASS_ALREADY_EXISTS) {
		return nil
	}
	return winUIError(callErr)
}

func nextMessage(m *msg) (int32, error) {
	r, _, callErr := procGetMessage.Call(uintptr(unsafe.Pointer(m)), 0, 0, 0)
	status := int32(r)
	if status == -1 {
		return -1, winUIError(callErr)
	}
	return status, nil
}

func winUIError(err error) error {
	if err == nil || err == syscall.Errno(0) {
		return syscall.EINVAL
	}
	return err
}
func move(hwnd uintptr, x, y, w, h int32) {
	procMoveWindow.Call(hwnd, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 0)
}
func drawOwnerButton(lparam uintptr) bool {
	if lparam == 0 {
		return false
	}
	var item drawItemStruct
	procRtlMoveMemory.Call(uintptr(unsafe.Pointer(&item)), lparam, unsafe.Sizeof(item))
	if item.HDC == 0 || item.HwndItem == 0 {
		return false
	}
	selected := item.ItemState&ODS_SELECTED != 0
	bg := buttonFillColor(int(item.CtlID), selected)
	border := buttonBorderColor(int(item.CtlID), selected)
	bgBrush := createSolidBrush(bg)
	borderBrush := createSolidBrush(border)
	defer procDeleteObject.Call(bgBrush)
	defer procDeleteObject.Call(borderBrush)

	procFillRect.Call(item.HDC, uintptr(unsafe.Pointer(&item.RcItem)), bgBrush)
	procFrameRect.Call(item.HDC, uintptr(unsafe.Pointer(&item.RcItem)), borderBrush)
	procSetBkMode.Call(item.HDC, TRANSPARENT)
	procSetTextColor.Call(item.HDC, uintptr(rgb(255, 255, 255)))

	text := getText(item.HwndItem)
	textRect := item.RcItem
	textRect.Left += 4
	textRect.Right -= 4
	procDrawText.Call(item.HDC, uintptr(unsafe.Pointer(utf16Ptr(text))), ^uintptr(0), uintptr(unsafe.Pointer(&textRect)), DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	return true
}

func buttonFillColor(id int, selected bool) uint32 {
	r, g, b := buttonRGB(id)
	if selected {
		r = shadeByte(r, -24)
		g = shadeByte(g, -24)
		b = shadeByte(b, -24)
	}
	return rgb(r, g, b)
}

func buttonBorderColor(id int, selected bool) uint32 {
	r, g, b := buttonRGB(id)
	if selected {
		return rgb(shadeByte(r, 18), shadeByte(g, 18), shadeByte(b, 18))
	}
	return rgb(shadeByte(r, 38), shadeByte(g, 38), shadeByte(b, 38))
}

func buttonRGB(id int) (byte, byte, byte) {
	part := 0
	if id >= rowBase {
		part = (id - rowBase) % rowStep
	}
	switch {
	case id == ID_START:
		return 37, 156, 105
	case id == ID_STOP:
		return 220, 80, 92
	case id == ID_ADD:
		return 64, 126, 231
	case id == ID_ALL:
		return 20, 184, 166
	case id == ID_NONE:
		return 234, 151, 54
	case id == ID_SAVE:
		return 139, 92, 246
	case id == 302:
		return 59, 130, 246
	case id == 303:
		return 100, 116, 139
	case part == 3:
		return 56, 189, 248
	case part == 5:
		return 45, 212, 191
	case part == 6:
		return 244, 89, 114
	default:
		return 71, 85, 105
	}
}

func shadeByte(v byte, delta int) byte {
	n := int(v) + delta
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return byte(n)
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
func messageChoice(hwnd uintptr, text, title string, flags uintptr) uintptr {
	r, _, _ := procMessageBox.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr(title))), flags)
	return r
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
func pointFromLParam(v uintptr) (int32, int32) {
	raw := uint32(v)
	return int32(int16(raw & 0xffff)), int32(int16((raw >> 16) & 0xffff))
}

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

func applyDarkTitleBar(hwnd uintptr) {
	if hwnd == 0 {
		return
	}
	dark := int32(1)
	procDwmSetWindowAttr.Call(hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE, uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))
	procDwmSetWindowAttr.Call(hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE_OLD, uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))

	caption := rgb(18, 20, 24)
	text := rgb(248, 250, 252)
	procDwmSetWindowAttr.Call(hwnd, DWMWA_CAPTION_COLOR, uintptr(unsafe.Pointer(&caption)), unsafe.Sizeof(caption))
	procDwmSetWindowAttr.Call(hwnd, DWMWA_TEXT_COLOR, uintptr(unsafe.Pointer(&text)), unsafe.Sizeof(text))
}

func applyMinWindowSize(hwnd uintptr, lparam uintptr) {
	if lparam == 0 {
		return
	}
	minW := modernMinWindowWidth(monitorWidth(hwnd))
	writeInt32(lparam+24, minW)
	writeInt32(lparam+28, 360)
}

func startupMinWindowWidth() int32 {
	if w := getSystemMetric(SM_CXSCREEN); w > 0 {
		return modernMinWindowWidth(w)
	}
	return modernMinWidth
}

func modernMinWindowWidth(screenW int32) int32 {
	minW := screenW / 2
	if minW < modernMinWidth {
		minW = modernMinWidth
	}
	return minW
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
