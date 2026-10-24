package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Version is set by `make dist VERSION=x.y.z` (or by hand). Shown in the status bar.
const Version = "0.1.0"

// App holds backend state. Every exported method is callable from the frontend
// through the generated bindings in frontend/src/lib/wailsjs/go/main/App.
type App struct {
	ctx      context.Context
	settings *SettingsStore
	updater  *Updater

	mu     sync.Mutex
	update UpdateInfo // the last check's result, which InstallUpdate installs
}

func NewApp() *App {
	return &App{settings: NewSettingsStore(appName), updater: NewUpdater(updateRepo, appName, Version)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	CleanupOldBinary()
	go a.autoCheck()
}

// autoCheck runs once at startup: at most once a day, and only if the setting is on.
// A new release raises the "update:available" event. Errors stay quiet: no network is normal.
func (a *App) autoCheck() {
	s := a.settings.Load()
	if !s.AutoUpdate || !DueForCheck(s.LastUpdateCheck, time.Now()) {
		return
	}
	info, err := a.CheckForUpdate()
	s.LastUpdateCheck = time.Now().UTC().Format(time.RFC3339)
	_ = a.settings.Save(s)
	if err == nil && info.Available {
		wruntime.EventsEmit(a.ctx, "update:available", info)
	}
}

// CheckForUpdate asks GitHub for the latest release.
func (a *App) CheckForUpdate() (UpdateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	info, err := a.updater.Check(ctx)
	if err == nil {
		a.mu.Lock()
		a.update = info
		a.mu.Unlock()
	}
	return info, err
}

// InstallUpdate downloads and verifies the release from the last check and swaps it in.
// RestartForUpdate then starts it.
func (a *App) InstallUpdate() error {
	a.mu.Lock()
	info := a.update
	a.mu.Unlock()
	if !info.Available {
		return fmt.Errorf("no update to install; check for updates first")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	return a.updater.Install(ctx, info)
}

// RestartForUpdate starts the new binary and quits this one.
func (a *App) RestartForUpdate() error {
	if err := Relaunch(); err != nil {
		return err
	}
	wruntime.Quit(a.ctx)
	return nil
}

// OpenURL opens a link (release notes, downloads) in the system browser.
func (a *App) OpenURL(url string) { wruntime.BrowserOpenURL(a.ctx, url) }

// AppInfo is what the status bar and dashboard show about the running build.
type AppInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Hostname  string `json:"hostname"`
	CPUs      int    `json:"cpus"`
}

func (a *App) AppInfo() AppInfo {
	host, _ := os.Hostname()
	return AppInfo{
		Name:      appName,
		Version:   Version,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Hostname:  host,
		CPUs:      runtime.NumCPU(),
	}
}

// Ping is the example round trip: frontend calls Go, Go answers. Replace with real work.
func (a *App) Ping(msg string) string {
	return fmt.Sprintf("> ACK %q from %s backend", msg, appName)
}

// GetSettings / SaveSettings persist user preferences (see settings.go).
func (a *App) GetSettings() Settings         { return a.settings.Load() }
func (a *App) SaveSettings(s Settings) error { return a.settings.Save(s) }
func (a *App) SettingsPath() string          { return a.settings.Path() }
