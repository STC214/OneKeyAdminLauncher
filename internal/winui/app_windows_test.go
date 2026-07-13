package winui

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

func TestAppWindowTitleIncludesBuildVersion(t *testing.T) {
	previous := appVersion
	t.Cleanup(func() { appVersion = previous })
	appVersion = "20260711_1430"
	if got, want := appWindowTitle(), "程序启动管理器 20260711_1430"; got != want {
		t.Fatalf("appWindowTitle() = %q, want %q", got, want)
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

func TestOverlappingLaunchBatchesReturnOneCombinedFailureSet(t *testing.T) {
	state := &appState{}
	state.beginLaunchBatch()
	state.beginLaunchBatch()
	if failures, completed := state.finishLaunchBatch([]string{"first.exe: failed"}); completed || failures != nil {
		t.Fatalf("first completion = (%v, %v), want (nil, false)", failures, completed)
	}
	failures, completed := state.finishLaunchBatch([]string{"second.exe: failed", "third.exe: failed"})
	want := []string{"first.exe: failed", "second.exe: failed", "third.exe: failed"}
	if !completed || !reflect.DeepEqual(failures, want) {
		t.Fatalf("combined completion = (%v, %v), want (%v, true)", failures, completed, want)
	}
	if state.launchBatches != 0 || len(state.launchFailures) != 0 {
		t.Fatalf("launch aggregate not reset: batches=%d failures=%v", state.launchBatches, state.launchFailures)
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

func TestLaunchEnabledProgramsLimitsConcurrencyAndPreservesFailureOrder(t *testing.T) {
	items := make([]config.ProgramItem, 13)
	for i := range 12 {
		items[i] = config.ProgramItem{Enabled: true, Path: fmt.Sprintf("tool-%02d.exe", i)}
	}
	items[12] = config.ProgramItem{Enabled: false, Path: "disabled.exe"}

	var active atomic.Int32
	var maximum atomic.Int32
	var calls atomic.Int32
	failures := launchEnabledPrograms(items, func(path string, _ bool) error {
		calls.Add(1)
		current := active.Add(1)
		for current > maximum.Load() && !maximum.CompareAndSwap(maximum.Load(), current) {
		}
		time.Sleep(20 * time.Millisecond)
		active.Add(-1)
		if path == "tool-02.exe" || path == "tool-09.exe" {
			return errors.New("launch failed")
		}
		return nil
	})

	if calls.Load() != 12 {
		t.Fatalf("launch calls = %d, want 12", calls.Load())
	}
	if maximum.Load() != maxConcurrentLaunches {
		t.Fatalf("maximum concurrency = %d, want %d", maximum.Load(), maxConcurrentLaunches)
	}
	want := []string{"tool-02.exe: launch failed", "tool-09.exe: launch failed"}
	if !reflect.DeepEqual(failures, want) {
		t.Fatalf("failures = %v, want %v", failures, want)
	}
}

func TestLaunchEnabledProgramsLimitsConcurrencyAcrossOverlappingBatches(t *testing.T) {
	items := make([]config.ProgramItem, 8)
	for i := range items {
		items[i] = config.ProgramItem{Enabled: true, Path: fmt.Sprintf("batch-tool-%02d.exe", i)}
	}

	var active atomic.Int32
	var maximum atomic.Int32
	launch := func(string, bool) error {
		current := active.Add(1)
		for {
			previous := maximum.Load()
			if current <= previous || maximum.CompareAndSwap(previous, current) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		active.Add(-1)
		return nil
	}

	var batches sync.WaitGroup
	batches.Add(2)
	for range 2 {
		go func() {
			defer batches.Done()
			if failures := launchEnabledPrograms(items, launch); len(failures) != 0 {
				t.Errorf("unexpected failures: %v", failures)
			}
		}()
	}
	batches.Wait()
	if maximum.Load() != maxConcurrentLaunches {
		t.Fatalf("maximum concurrency across batches = %d, want %d", maximum.Load(), maxConcurrentLaunches)
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
