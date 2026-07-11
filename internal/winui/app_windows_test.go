package winui

import (
	"errors"
	"reflect"
	"testing"
	"unicode/utf16"

	"program-launch-manager/internal/config"
)

func TestDialogFilter(t *testing.T) {
	got := dialogFilter("程序 (*.exe)", "*.exe", "所有文件 (*.*)", "*.*")
	want := utf16.Encode([]rune("程序 (*.exe)\x00*.exe\x00所有文件 (*.*)\x00*.*\x00\x00"))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("dialogFilter() = %v, want %v", got, want)
	}
	if len(got) < 2 || got[len(got)-1] != 0 || got[len(got)-2] != 0 {
		t.Fatal("dialog filter must end with two NUL code units")
	}
}

func TestSaveResultIgnoresStaleAutoFailure(t *testing.T) {
	state := &appState{}
	if _, notify := state.recordSaveResult(2, nil, false); notify {
		t.Fatal("manual success should not notify")
	}
	if _, notify := state.recordSaveResult(1, errors.New("old failure"), true); notify {
		t.Fatal("stale auto-save failure should be ignored")
	}
	if state.autoSaveErr != "" || state.autoSaveErrSeen {
		t.Fatalf("stale failure changed state: err=%q seen=%v", state.autoSaveErr, state.autoSaveErrSeen)
	}
}

func TestManualSuccessResetsAutoSaveNotification(t *testing.T) {
	state := &appState{}
	if _, notify := state.recordSaveResult(1, errors.New("first failure"), true); !notify {
		t.Fatal("first auto-save failure should notify")
	}
	state.recordSaveResult(2, nil, false)
	if state.autoSaveErr != "" || state.autoSaveErrSeen {
		t.Fatalf("manual success did not reset state: err=%q seen=%v", state.autoSaveErr, state.autoSaveErrSeen)
	}
	if _, notify := state.recordSaveResult(3, errors.New("new failure"), true); !notify {
		t.Fatal("a new auto-save failure after recovery should notify")
	}
}

func TestRepeatedAutoSaveFailureNotifiesOnce(t *testing.T) {
	state := &appState{}
	if _, notify := state.recordSaveResult(1, errors.New("failure"), true); !notify {
		t.Fatal("first failure should notify")
	}
	if _, notify := state.recordSaveResult(2, errors.New("failure"), true); notify {
		t.Fatal("repeated failure should not notify again")
	}
}

func TestOperationResultsRemainFIFO(t *testing.T) {
	state := &appState{}
	first := []string{"first"}
	state.queueOperationResult("launch", first)
	state.queueOperationResult("close", []string{"second"})
	first[0] = "mutated"
	got, ok := state.popOperationResult()
	if !ok || got.heading != "launch" || !reflect.DeepEqual(got.failures, []string{"first"}) {
		t.Fatalf("first result = %v, ok=%v", got, ok)
	}
	got, ok = state.popOperationResult()
	if !ok || got.heading != "close" || !reflect.DeepEqual(got.failures, []string{"second"}) {
		t.Fatalf("second result = %v, ok=%v", got, ok)
	}
	if _, ok := state.popOperationResult(); ok {
		t.Fatal("empty queue returned a result")
	}
}

func TestUniqueEnabledProcessNames(t *testing.T) {
	items := []config.ProgramItem{
		{Enabled: true, ProcessName: "Tool.exe"},
		{Enabled: true, SelectedProcess: "tool.EXE"},
		{Enabled: false, ProcessName: "disabled.exe"},
		{Enabled: true, Path: `C:\\Apps\\other.exe`},
	}
	want := []string{"Tool.exe", "other.exe"}
	if got := uniqueEnabledProcessNames(items); !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueEnabledProcessNames() = %v, want %v", got, want)
	}
}

func TestProcessDialogDestroyMarksDone(t *testing.T) {
	previous := dlg
	t.Cleanup(func() { dlg = previous })
	d := &processDialog{hwnd: 123}
	dlg = d
	processDlgProc(123, WM_DESTROY, 0, 0)
	if !d.done || d.hwnd != 0 {
		t.Fatalf("destroyed dialog state = done:%v hwnd:%d", d.done, d.hwnd)
	}
}

func TestUpdateProgramPathRefreshesAutomaticProcessName(t *testing.T) {
	item := config.ProgramItem{Path: `C:\\Tools\\old.exe`, ProcessName: "old.exe"}
	updateProgramPath(&item, `C:\\Tools\\new.exe`)
	if item.ProcessName != "new.exe" {
		t.Fatalf("ProcessName = %q, want new.exe", item.ProcessName)
	}
}

func TestUpdateProgramPathPreservesSelectedProcess(t *testing.T) {
	item := config.ProgramItem{Path: `C:\\Tools\\old.exe`, ProcessName: "old.exe", SelectedProcess: "custom.exe"}
	updateProgramPath(&item, `C:\\Tools\\new.exe`)
	if item.ProcessName != "old.exe" || item.SelectedProcess != "custom.exe" {
		t.Fatalf("manual binding changed: %#v", item)
	}
}
