package process

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLaunchRejectsEmbeddedNUL(t *testing.T) {
	if err := Launch("bad\x00path.exe", false); err == nil {
		t.Fatal("Launch() should reject a path containing NUL")
	}
}

func TestLaunchRejectsEmptyPath(t *testing.T) {
	if err := Launch("   ", false); err == nil {
		t.Fatal("Launch() should reject an empty path")
	}
}

func TestLaunchRejectsMissingFileAndDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.exe")
	if err := Launch(missing, false); err == nil || !strings.Contains(err.Error(), "程序文件不存在") {
		t.Fatalf("Launch(missing) = %v", err)
	}
	if err := Launch(t.TempDir(), false); err == nil || !strings.Contains(err.Error(), "程序路径是目录") {
		t.Fatalf("Launch(directory) = %v", err)
	}
}

func TestCloseByNameIgnoresEmptyAndMissingProcess(t *testing.T) {
	if err := CloseByName(""); err != nil {
		t.Fatalf("CloseByName(empty) = %v", err)
	}
	if err := CloseByName("__one_key_admin_open_missing_process__.exe"); err != nil {
		t.Fatalf("CloseByName(missing) = %v", err)
	}
}

func TestCloseByNamesIgnoresEmptyDuplicateAndMissingProcesses(t *testing.T) {
	err := CloseByNames([]string{"", "  ", "__one_key_admin_open_missing_process__.exe", "__ONE_KEY_ADMIN_OPEN_MISSING_PROCESS__.EXE"})
	if err != nil {
		t.Fatalf("CloseByNames(missing) = %v", err)
	}
}

func TestCloseTargetsStillRunningFindsCurrentProcessAndRejectsPIDNameMismatch(t *testing.T) {
	pid := uint32(os.Getpid())
	entries, err := processEntries()
	if err != nil {
		t.Fatal(err)
	}
	var name string
	for _, entry := range entries {
		if entry.PID == pid {
			name = entry.Name
			break
		}
	}
	if name == "" {
		t.Fatal("current process was not enumerated")
	}
	targets := []closeTarget{{pid: pid, name: strings.ToUpper(name)}, {pid: pid, name: "not-the-current-process.exe"}}
	alive, err := closeTargetsStillRunning(targets)
	if err != nil {
		t.Fatal(err)
	}
	if len(alive) != 1 || alive[0].pid != pid || !strings.EqualFold(alive[0].name, name) {
		t.Fatalf("closeTargetsStillRunning() = %v, want only PID %d named %s", alive, pid, name)
	}
}
