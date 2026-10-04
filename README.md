# templates

Starting points for every cyberpunk_apps project. Each new project lives in `cyberpunk_apps/<name>/`
as its own git repo, with the standard environment (uv venv + nodeenv) and NEONDECK.

```bash
templates/new-project.sh <name> site   # SvelteKit website, prerendered static (adapter-static)
templates/new-project.sh <name> app    # Wails v2 desktop app (Go + SvelteKit SPA)
```

The script copies the template, renames it, creates `.venv` with node, installs, type-checks, builds
(apps also run Go tests and produce a native binary), then runs `git init` with a first commit.
It refuses to overwrite an existing folder.

| Template | What you get |
|---|---|
| `site/` | AppShell, hero with MoonScroll, numbered bilingual sections, paper "about" page, "signal lost" 404, `config.ts` for identity |
| `wails-app/` | frameless window with NEONDECK title bar + window controls, sidebar nav, dashboard, Go-persisted settings (accent/glow/scanlines), `backend.ts` with browser stubs, Makefile |

## Requirements
- `uv`, Go ≥ 1.25, Wails CLI v2.14 (`go install github.com/wailsapp/wails/v2/cmd/wails@v2.14.0`).
- Windows builds cross-compile from Linux: `make windows` (.exe) or `make windows-installer` (NSIS, needs `nsis`).
- macOS builds cannot be made on Linux (Wails refuses; needs Apple SDK). Use a Mac or macOS CI.
- Linux: GTK3 + WebKitGTK 4.1. Builds use `-tags webkit2_41` (the Makefile and script add it).
  Ignore `wails doctor` saying "libwebkit not found": it checks for 4.0.
- NEONDECK built: `cd sharable_assets/neondeck && npm run build`.

## Working on the templates themselves
```bash
cd templates && source .venv/bin/activate
cd site && npm run dev                          # or: cd wails-app && make dev
```
After changing a template, test it end to end: `./new-project.sh zz-test site` (and `app`), then delete `../zz-test`.
