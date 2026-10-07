package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

// The SvelteKit frontend builds to frontend/build (adapter-static, SPA mode, hash router).
//
//go:embed all:frontend/build
var assets embed.FS

// appName is rewritten by new-project.sh. It names the window and the settings folder.
const appName = "wails-app"

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     appName,
		Width:     1280,
		Height:    820,
		MinWidth:  960,
		MinHeight: 640,
		// Wails defaults the max size to the monitor the window opens on. With mixed
		// portrait and landscape monitors, xfwm4 then refuses to maximise. A large cap avoids that.
		MaxWidth:  16384,
		MaxHeight: 16384,
		// NEONDECK draws its own title bar (AppShell top bar = drag region + window controls).
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// --nd-bg (#070818): no white flash before the frontend paints.
		BackgroundColour: &options.RGBA{R: 7, G: 8, B: 24, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			ProgramName:         appName,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
