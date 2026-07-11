package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMalformedConfigIsNotOverwritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "launcher_config.json")
	original := []byte(`{"programs":[`)
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadFromPaths(path, filepath.Join(dir, "legacy.json")); err == nil {
		t.Fatal("loadFromPaths() should reject malformed JSON")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, original) {
		t.Fatalf("malformed config was changed: %q", got)
	}
}

func TestSaveToPathRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "launcher_config.json")
	want := File{Programs: []ProgramItem{{Path: `C:\\Tools\\tool.exe`, ProcessName: "tool.exe", Enabled: true}}, Window: DefaultWindow()}
	if err := saveToPath(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := readAny(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readAny() = %#v, want %#v", got, want)
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".launcher_config-*.tmp"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("temporary files remain: %v, err=%v", matches, err)
	}
}

func TestLoadFromPathsMigratesLegacyConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "launcher_config.json")
	legacy := filepath.Join(dir, "data", "launcher_config.json")
	if err := os.MkdirAll(filepath.Dir(legacy), 0755); err != nil {
		t.Fatal(err)
	}
	want := File{Programs: []ProgramItem{{Path: `C:\Tools\legacy.exe`, Enabled: true}}, Window: DefaultWindow()}
	data := []byte(`[{"path":"C:\\Tools\\legacy.exe","enabled":true}]`)
	if err := os.WriteFile(legacy, data, 0644); err != nil {
		t.Fatal(err)
	}
	got, err := loadFromPaths(path, legacy)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loadFromPaths() = %#v, want %#v", got, want)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("legacy config still exists: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("migrated config missing: %v", err)
	}
}
