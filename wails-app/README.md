# wails-app

Desktop app: Wails v2 (Go) + SvelteKit (Svelte 5) + NEONDECK.

## Setup
```bash
source .venv/bin/activate          # node/npm live in the venv (uv + nodeenv)
make dev                           # live reload
make build                         # binary in build/bin/
make test                          # go vet + go test
make bindings                      # after adding/changing exported Go methods
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

## Notes
- The top bar is the title bar (`--wails-draggable: drag`). Buttons and links inside opt out.
- NEONDECK is copied in (`install-links=true`). After changing it in `sharable_assets`, rebuild it there,
  then run `cd frontend && npm run update:neondeck` here.
