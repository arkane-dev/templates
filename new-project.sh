#!/usr/bin/env bash
# Create a new cyberpunk_apps project from a template.
#
#   templates/new-project.sh <name> <site|app>
#
#   site  SvelteKit website (static, NEONDECK)
#   app   Wails v2 desktop app (Go + SvelteKit, NEONDECK)
#
# The project lands in cyberpunk_apps/<name>/ as its own git repo, with the standard
# environment (uv venv + nodeenv), NEONDECK copied in, a passing build and a first commit.
set -euo pipefail

die() { echo "error: $*" >&2; exit 1; }
step() { printf '\n\033[95m▶ %s\033[0m\n' "$*"; }

[[ $# -eq 2 ]] || die "usage: $0 <name> <site|app>"
name=$1
kind=$2
[[ $name =~ ^[a-z][a-z0-9-]*[a-z0-9]$ ]] || die "name must be kebab-case (a-z, 0-9, -), e.g. genai-calculator"

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
root=$(dirname "$here")                       # cyberpunk_apps/
target="$root/$name"
neondeck="$root/sharable_assets/neondeck"
brand="$(echo "$name" | tr 'a-z-' 'A-Z_')_"   # genai-calculator -> GENAI_CALCULATOR_
title="$(echo "$name" | sed -E 's/(^|-)([a-z])/\1\u\2/g; s/-/ /g')"   # -> Genai Calculator

case $kind in
	site) src="$here/site" ;;
	app) src="$here/wails-app" ;;
	*) die "kind must be 'site' or 'app'" ;;
esac
[[ -e $target ]] && die "$target already exists"
[[ -f $neondeck/dist/index.js ]] || die "NEONDECK is not built: cd $neondeck && npm run build"
command -v uv >/dev/null || die "uv not found"
if [[ $kind == app ]]; then
	command -v go >/dev/null || die "go not found (sudo pacman -S go)"
	wails=$(command -v wails 2>/dev/null || echo "$HOME/go/bin/wails")
	[[ -x $wails ]] || die "wails not found: go install github.com/wailsapp/wails/v2/cmd/wails@v2.14.0"
	tags=""
	[[ $(uname -s) == Linux ]] && tags="-tags webkit2_41"
fi

step "Copying $kind template → $target"
mkdir -p "$target"
# Generated output never travels. For apps, build/ itself stays: it holds the icon and platform manifests.
excludes=(--exclude=node_modules --exclude=.svelte-kit --exclude=.venv)
if [[ $kind == site ]]; then excludes+=(--exclude=./build); else excludes+=(--exclude=./build/bin --exclude=./frontend/build); fi
tar -C "$src" "${excludes[@]}" -cf - . | tar -C "$target" -xf -
cd "$target"

step "Renaming template → $name"
if [[ $kind == site ]]; then
	sed -i "s/\"name\": \"template\"/\"name\": \"$name\"/" package.json package-lock.json
	sed -i "s|file:../../sharable_assets/neondeck|file:../sharable_assets/neondeck|" package.json package-lock.json
	sed -i -e "s/name: 'template'/name: '$name'/" \
		-e "s/brand: 'TEMPLATE_'/brand: '$brand'/" \
		-e "s/title: 'Template'/title: '$title'/" src/lib/config.ts
	sed -i "s/^# site$/# $name/" README.md
else
	# "wails-app" appears in go.mod, wails.json, main.go, config.ts, backend.ts, package files, README.
	grep -rlI --exclude-dir=wailsjs "wails-app" . | xargs sed -i "s/wails-app/$name/g"
	sed -i "s/WAILS_APP_/$brand/" frontend/src/lib/config.ts
	sed -i "s/\"productName\": \"Wails App\"/\"productName\": \"$title\"/" wails.json
	sed -i "s|file:../../../sharable_assets/neondeck|file:../../sharable_assets/neondeck|" frontend/package.json frontend/package-lock.json
fi

step "Environment: uv venv + nodeenv"
uv venv -q .venv
uv pip install -q --python .venv nodeenv
.venv/bin/nodeenv -p --node=lts -q >/dev/null
# shellcheck disable=SC1091
source .venv/bin/activate

if [[ $kind == site ]]; then
	step "npm install + check + build"
	npm install --silent
	npm run check
	npm run build >/dev/null
else
	step "Frontend: npm install + check + build"
	(cd frontend && npm install --silent && npm run check && npm run build >/dev/null)
	step "Go: bindings, tidy, vet, test, native build"
	# shellcheck disable=SC2086
	"$wails" generate module $tags >/dev/null
	go mod tidy
	go vet ./... && go test ./...
	# shellcheck disable=SC2086
	"$wails" build $tags >/dev/null
	echo "binary: build/bin/$name"
fi

step "git init"
git init -q -b main
git add -A
git commit -q -m "Initial commit from the $kind template"
git log --oneline -1

cat <<EOF

✔ $name is ready at $target

  cd $name && source .venv/bin/activate
EOF
if [[ $kind == site ]]; then
	echo "  npm run dev"
else
	echo "  make dev          # (make build / make test / make bindings)"
fi
echo "  Edit identity in $( [[ $kind == site ]] && echo src/lib/config.ts || echo frontend/src/lib/config.ts )"
