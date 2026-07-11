package process

import (
	"os"
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

func TestCloseByNameIgnoresEmptyAndMissingProcess(t *testing.T) {
	if err := CloseByName(""); err != nil {
		t.Fatalf("CloseByName(empty) = %v", err)
	}
	if err := CloseByName("__one_key_admin_open_missing_process__.exe"); err != nil {
		t.Fatalf("CloseByName(missing) = %v", err)
	}
}

func TestTargetsStillRunningFindsCurrentProcess(t *testing.T) {
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
	alive, err := targetsStillRunning([]uint32{pid}, name)
	if err != nil {
		t.Fatal(err)
	}
	if len(alive) != 1 || alive[0] != pid {
		t.Fatalf("targetsStillRunning() = %v, want [%d]", alive, pid)
	}
}
