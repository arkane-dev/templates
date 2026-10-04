package main

import (
	"path/filepath"
	"testing"
)

func TestSettingsRoundTrip(t *testing.T) {
	s := &SettingsStore{path: filepath.Join(t.TempDir(), "x", "settings.json")}
	if got := s.Load(); got != defaultSettings() {
		t.Fatalf("missing file should load defaults, got %+v", got)
	}
	want := Settings{Accent: "jade", GlowSize: 0.5, Scanlines: true}
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	if got := s.Load(); got != want {
		t.Fatalf("round trip: got %+v want %+v", got, want)
	}
}
