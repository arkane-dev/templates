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
0. Once per machine: `go run ./scripts/release-sign -keygen`. It writes the update signing key to
   `~/.config/cyberpunk-apps/release-signing.key` and prints the public key, which goes in
   `updatePublicKey` in `updater.go`. All the apps share one key. **Back the private key up.** Without
   it, installed apps can't verify new releases, and users have to download the next one by hand.
1. `make dist VERSION=x.y.z` sets the version in `app.go` and `wails.json`, builds Linux and Windows,
   signs each archive (`*.sig`), and writes them and `SHA256SUMS` to `build/dist/`.
2. Commit, tag `vx.y.z` and push. Then `gh release create vx.y.z -F notes.md build/dist/*`.
3. Publishing the release starts `.github/workflows/macos.yml`. It builds a universal macOS app
   on a GitHub Mac runner and attaches the zip and its checksum to the release. To build for
   a release that already exists: `gh workflow run macos.yml -f tag=vx.y.z`.

The Mac app is ad-hoc signed only (no Apple Developer ID, not notarised). Tell users how to
get past Gatekeeper: *System Settings → Privacy & Security → Open Anyway*, or
`xattr -dr com.apple.quarantine /Applications/<app>.app`. Needs macOS 12 or later.

## Auto-update
`updater.go` checks the latest GitHub release of `updateRepo` (set in `main.go`).

- **When:** once a day at startup, if Settings → Updates allows it, which is the default. "Check now" checks any time.
- **What it does:** shows a banner with Update, What's new and Later. Update downloads this system's archive and its `.sig`,
  checks the ed25519 signature against `updatePublicKey`, and swaps the new binary in. Restart runs it.
  An archive whose signature fails is never installed.
- **Linux and Windows** update in place. Windows renames the running `.exe` aside first, and the next start deletes it.
- **macOS** gets the banner and a download link. Its build happens in CI, which doesn't hold the signing key.
- **No write access** to the app's folder (say `/usr/bin` or Program Files): the banner offers the download instead.
- **Testing:** `UPDATE_API_BASE=http://127.0.0.1:8799 make dev` points it at a local stand-in for the GitHub API.
  Downloads from it still have to pass the signature check.

## Notes
- `main.go` sets a large max window size. Without it, Wails caps the window at the size of the
  monitor it opened on, and xfwm4 won't maximise on a monitor with a different shape.
- The top bar is the title bar (`--wails-draggable: drag`). Buttons and links inside opt out.
- NEONDECK is copied in (`install-links=true`). After changing it in `sharable_assets`, rebuild it there,
  then run `cd frontend && npm run update:neondeck` here.
