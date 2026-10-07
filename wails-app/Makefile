# Wails v2 on Linux needs the webkit2_41 tag (Manjaro/Arch ship WebKitGTK 4.1, not 4.0).
# Always go through these targets so the tag is never forgotten.
WAILS ?= $(shell command -v wails 2>/dev/null || echo $(HOME)/go/bin/wails)
TAGS  := $(if $(filter Linux,$(shell uname -s)),-tags webkit2_41,)

.PHONY: dev build test bindings frontend windows windows-installer dist

dev:            ## live-reload desktop app (Go + Svelte)
	$(WAILS) dev $(TAGS)

build:          ## release binary in build/bin/
	$(WAILS) build $(TAGS) -clean

frontend:       ## build the SvelteKit frontend (Go embeds frontend/build)
	cd frontend && npm run build

bindings: frontend  ## regenerate frontend/src/lib/wailsjs after changing Go methods
	$(WAILS) generate module $(TAGS)

test: frontend  ## Go tests (need frontend/build to exist for go:embed)
	go vet ./... && go test ./...

# Cross-compiles from Linux: WebView2 is pure Go, no Windows toolchain needed.
# Targets need Windows 10/11 (WebView2 runtime ships with them).
windows:        ## build/bin/<app>.exe
	$(WAILS) build -platform windows/amd64

windows-installer:  ## NSIS installer in build/bin/ (needs makensis: sudo pacman -S nsis)
	$(WAILS) build -platform windows/amd64 -nsis

# Release archives (Linux + Windows) and SHA256SUMS in build/dist/. Sets the version first.
dist:           ## make dist VERSION=x.y.z
	@test -n "$(VERSION)" || { echo "usage: make dist VERSION=x.y.z"; exit 1; }
	scripts/dist.sh $(VERSION)

# macOS: Wails cannot cross-compile to Mac (needs Apple's SDK + toolchain).
# .github/workflows/macos.yml builds it on a Mac runner when a release is published.

