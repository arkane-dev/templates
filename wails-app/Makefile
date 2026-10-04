# Wails v2 on Linux needs the webkit2_41 tag (Manjaro/Arch ship WebKitGTK 4.1, not 4.0).
# Always go through these targets so the tag is never forgotten.
WAILS ?= $(shell command -v wails 2>/dev/null || echo $(HOME)/go/bin/wails)
TAGS  := $(if $(filter Linux,$(shell uname -s)),-tags webkit2_41,)

.PHONY: dev build test bindings frontend

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
