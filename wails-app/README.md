# wails-app

Desktop app: Wails v2 (Go) + SvelteKit (Svelte 5) + NEONDECK.

## Setup
```bash
source .venv/bin/activate          # node/npm live in the venv (uv + nodeenv)
make dev                           # live reload
make build                         # binary in build/bin/
make test                          # go vet + go test
make bindings                      # after adding/changing exported Go methods
make dist VERSION=0.2.0            # release archives + SHA256SUMS in build/dist/
```
UI-only work without Go: `cd frontend && npm run dev`. Backend calls fall back to stubs
(`src/lib/backend.ts`), and the status bar shows `BROWSER` instead of `NATIVE`.

## Layout
| Path | What |
|---|---|
| `main.go` | window options: frameless, NEONDECK background, `appName` |
| `app.go` | exported methods = frontend API (`AppInfo`, `Ping`, settings) |
| `settings.go` | JSON settings in `~/.config/<app>/settings.json`, atomic writes |
| `env_linux.go` | disables WebKit's DMA-BUF renderer (blank window on NVIDIA) |
| `frontend/src/lib/backend.ts` | typed Go calls + browser stubs. Call Go only through this |
| `frontend/src/lib/wailsjs/` | generated bindings (`make bindings`), committed |
| `frontend/src/lib/config.ts` | app name, brand, sidebar nav |
| `frontend/src/routes/` | pages; hash router (`#/settings`) |

## Release
1. `make dist VERSION=x.y.z` sets the version in `app.go` and `wails.json`, builds Linux and Windows,
   and writes the archives and `SHA256SUMS` to `build/dist/`.
2. Commit, tag `vx.y.z` and push. Then `gh release create vx.y.z -F notes.md build/dist/*`.
3. Publishing the release starts `.github/workflows/macos.yml`. It builds a universal macOS app
   on a GitHub Mac runner and attaches the zip and its checksum to the release. To build for
   a release that already exists: `gh workflow run macos.yml -f tag=vx.y.z`.

The Mac app is ad-hoc signed only (no Apple Developer ID, not notarised). Tell users how to
get past Gatekeeper: *System Settings → Privacy & Security → Open Anyway*, or
`xattr -dr com.apple.quarantine /Applications/<app>.app`. Needs macOS 12 or later.

## Notes
- `main.go` sets a large max window size. Without it, Wails caps the window at the size of the
  monitor it opened on, and xfwm4 won't maximise on a monitor with a different shape.
- The top bar is the title bar (`--wails-draggable: drag`). Buttons and links inside opt out.
- NEONDECK is copied in (`install-links=true`). After changing it in `sharable_assets`, rebuild it there,
  then run `cd frontend && npm run update:neondeck` here.
