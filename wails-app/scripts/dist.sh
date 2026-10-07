#!/usr/bin/env bash
# Package a release: set the version, build Linux and Windows, and write the archives
# plus SHA256SUMS to build/dist/. Run from the project root on Linux: make dist VERSION=x.y.z
#
# Publishing stays a separate, deliberate step (see README "Release"). The macOS build
# runs in CI when the release is published and attaches itself.
set -euo pipefail

version=${1:?usage: scripts/dist.sh x.y.z}
[[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "version must look like 1.2.3" >&2; exit 1; }
name=$(sed -n 's/.*"outputfilename": *"\([^"]*\)".*/\1/p' wails.json)

# One version, three places: the status bar, the app bundle metadata, the archive names.
sed -i "s/^const Version = \".*\"/const Version = \"$version\"/" app.go
sed -i "s/\"productVersion\": *\"[^\"]*\"/\"productVersion\": \"$version\"/" wails.json

make build
make windows

dist=build/dist
rm -rf "$dist" && mkdir -p "$dist"
tar czf "$dist/$name-$version-linux-amd64.tar.gz" -C build/bin "$name"
(cd build/bin && zip -q "../dist/$name-$version-windows-amd64.zip" "$name.exe")
(cd "$dist" && sha256sum -- * > SHA256SUMS)

echo
ls -la "$dist"
cat <<NEXT

Next:
  git commit -am "Release $version" && git tag v$version && git push && git push origin v$version
  gh release create v$version -t "<App name> $version" -F notes.md $dist/*
NEXT
