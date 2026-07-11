package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

type ProgramItem struct {
	Path            string `json:"path"`
	ProcessName     string `json:"process_name"`
	SelectedProcess string `json:"selected_process"`
	IsUWP           bool   `json:"is_uwp"`
	Enabled         bool   `json:"enabled"`
}

type WindowState struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

type File struct {
	Programs []ProgramItem `json:"programs"`
	Window   WindowState   `json:"window"`
}

func DefaultWindow() WindowState {
	return WindowState{X: -1, Y: -1, W: 800, H: 594}
}

func DefaultFile() File {
	return File{Programs: []ProgramItem{}, Window: DefaultWindow()}
}

func ConfigPath() (string, error) {
	root, err := appRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "launcher_config.json"), nil
}

func LegacyConfigPath() (string, error) {
	root, err := appRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "data", "launcher_config.json"), nil
}

func Load() (File, error) {
	path, err := ConfigPath()
	if err != nil {
		return DefaultFile(), err
	}
	legacy, legacyErr := LegacyConfigPath()
	if legacyErr != nil {
		return DefaultFile(), legacyErr
	}
	return loadFromPaths(path, legacy)
}

func loadFromPaths(path, legacy string) (File, error) {
	if _, err := os.Stat(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return DefaultFile(), fmt.Errorf("检查配置 %s: %w", path, err)
		}
		if _, legacyErr := os.Stat(legacy); legacyErr == nil {
			cfg, err := readAny(legacy)
			if err != nil {
				return cfg, err
			}
			if saveErr := saveToPath(path, cfg); saveErr != nil {
				return cfg, fmt.Errorf("迁移配置到 %s: %w", path, saveErr)
			}
			_ = os.Remove(legacy)
			return cfg, nil
		} else if !errors.Is(legacyErr, os.ErrNotExist) {
			return DefaultFile(), fmt.Errorf("检查旧配置 %s: %w", legacy, legacyErr)
		}
		return DefaultFile(), nil
	}
	cfg, err := readAny(path)
	if err != nil {
		return DefaultFile(), fmt.Errorf("读取配置 %s: %w", path, err)
	}
	return cfg, nil
}

func Save(cfg File) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return saveToPath(path, cfg)
}

func saveToPath(path string, cfg File) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(path, data)
}

func writeAtomic(path string, data []byte) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".launcher_config-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()
	if err = tmp.Chmod(0644); err != nil {
		return err
	}
	if _, err = tmp.Write(data); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	from, err := windows.UTF16PtrFromString(tmpPath)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(from, to, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH)
}

func readAny(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultFile(), err
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	var cfg File
	if err := json.Unmarshal(data, &cfg); err == nil && (cfg.Programs != nil || cfg.Window != (WindowState{})) {
		if cfg.Window.W == 0 || cfg.Window.H == 0 {
			cfg.Window = DefaultWindow()
		}
		if cfg.Programs == nil {
			cfg.Programs = []ProgramItem{}
		}
		return cfg, nil
	}

	var legacy []ProgramItem
	if err := json.Unmarshal(data, &legacy); err != nil {
		return DefaultFile(), err
	}
	return File{Programs: legacy, Window: DefaultWindow()}, nil
}

func appRoot() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}
