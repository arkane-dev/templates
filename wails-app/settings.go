package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// Settings are user preferences, stored as JSON in the OS config dir
// (~/.config/<app>/settings.json on Linux). Add fields freely; missing fields get defaults.
type Settings struct {
	Accent    string  `json:"accent"`    // NEONDECK accent: magenta | cyan | yellow | red | violet | jade | gold
	GlowSize  float64 `json:"glowSize"`  // multiplies --nd-glow-size; 0 = off
	Scanlines bool    `json:"scanlines"` // CRT overlay on the whole window
}

func defaultSettings() Settings {
	return Settings{Accent: "magenta", GlowSize: 1, Scanlines: false}
}

type SettingsStore struct {
	mu   sync.Mutex
	path string
}

func NewSettingsStore(app string) *SettingsStore {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return &SettingsStore{path: filepath.Join(dir, app, "settings.json")}
}

func (s *SettingsStore) Path() string { return s.path }

// Load returns saved settings, or defaults if none exist or the file is unreadable.
func (s *SettingsStore) Load() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := defaultSettings()
	b, err := os.ReadFile(s.path)
	if err != nil {
		return out
	}
	_ = json.Unmarshal(b, &out) // keep defaults for anything missing or malformed
	return out
}

// Save writes atomically (temp file + rename) so a crash never leaves half a file.
func (s *SettingsStore) Save(v Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, s.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
