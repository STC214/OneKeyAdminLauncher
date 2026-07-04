package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	legacy, _ := LegacyConfigPath()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, legacyErr := os.Stat(legacy); legacyErr == nil {
			cfg, err := readAny(legacy)
			if err != nil {
				return cfg, err
			}
			if saveErr := Save(cfg); saveErr == nil {
				_ = os.Remove(legacy)
			}
			return cfg, nil
		}
		return DefaultFile(), nil
	}
	return readAny(path)
}

func Save(cfg File) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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
