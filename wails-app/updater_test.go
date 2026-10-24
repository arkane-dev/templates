package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewer(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"0.1.2", "0.1.1", true}, {"v1.0.0", "0.9.9", true}, {"0.1.1", "0.1.1", false},
		{"0.1.0", "0.1.1", false}, {"0.10.0", "0.9.0", true}, {"1.2.3-beta", "1.2.2", true}, {"", "0.0.1", false},
	}
	for _, c := range cases {
		if got := Newer(c.a, c.b); got != c.want {
			t.Errorf("Newer(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestDueForCheck(t *testing.T) {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	if !DueForCheck("", now) || !DueForCheck("garbage", now) {
		t.Error("no or bad last check should be due")
	}
	if DueForCheck(now.Add(-time.Hour).Format(time.RFC3339), now) {
		t.Error("an hour ago should not be due")
	}
	if !DueForCheck(now.Add(-25*time.Hour).Format(time.RFC3339), now) {
		t.Error("a day ago should be due")
	}
}

func tarGz(t *testing.T, name string, body []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(body)), Typeflag: tar.TypeReg})
	tw.Write(body)
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

func zipOf(t *testing.T, name string, body []byte) []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create(name)
	w.Write(body)
	zw.Close()
	return buf.Bytes()
}

// fakeGitHub serves a latest release with one archive and its signature.
func fakeGitHub(t *testing.T, tag string, files map[string][]byte) *httptest.Server {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/me/app/releases/latest" {
			var assets []string
			for name, b := range files {
				assets = append(assets, fmt.Sprintf(`{"name":%q,"browser_download_url":"%s/dl/%s","size":%d}`, name, srv.URL, name, len(b)))
			}
			fmt.Fprintf(w, `{"tag_name":%q,"html_url":"https://example.com/r","body":"notes","assets":[%s]}`, tag, strings.Join(assets, ","))
			return
		}
		if b, ok := files[strings.TrimPrefix(r.URL.Path, "/dl/")]; ok {
			w.Write(b)
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setup(t *testing.T, goos string, archive []byte, sign ed25519.PrivateKey) (*Updater, string) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	if sign == nil {
		sign = priv
	}
	name := map[string]string{"linux": "app-0.2.0-linux-amd64.tar.gz", "windows": "app-0.2.0-windows-amd64.zip", "darwin": "app-0.2.0-macos-universal.zip"}[goos]
	sig := base64.StdEncoding.EncodeToString(ed25519.Sign(sign, archive))
	srv := fakeGitHub(t, "v0.2.0", map[string][]byte{name: archive, name + ".sig": []byte(sig)})
	dir := t.TempDir()
	exe := filepath.Join(dir, "app")
	if goos == "windows" {
		exe += ".exe"
	}
	os.WriteFile(exe, []byte("old binary"), 0o755)
	u := &Updater{Repo: "me/app", Name: "app", Version: "0.1.0", API: srv.URL, Client: srv.Client(), Key: pub,
		Exe: func() (string, error) { return exe, nil }, GOOS: goos}
	return u, exe
}

func TestCheckAndInstallLinux(t *testing.T) {
	u, exe := setup(t, "linux", tarGz(t, "app", []byte("new binary")), nil)
	info, err := u.Check(context.Background())
	if err != nil || !info.Available || !info.CanInstall || info.Latest != "0.2.0" || info.Notes != "notes" {
		t.Fatalf("check: %+v %v", info, err)
	}
	if err := u.Install(context.Background(), info); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(exe); string(b) != "new binary" {
		t.Fatalf("binary not replaced: %q", b)
	}
	if st, _ := os.Stat(exe); st.Mode().Perm()&0o100 == 0 {
		t.Error("new binary is not executable")
	}
}

func TestInstallWindowsMovesTheRunningExeAside(t *testing.T) {
	u, exe := setup(t, "windows", zipOf(t, "app.exe", []byte("new exe")), nil)
	info, _ := u.Check(context.Background())
	if err := u.Install(context.Background(), info); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(exe); string(b) != "new exe" {
		t.Fatalf("exe not replaced: %q", b)
	}
	if b, _ := os.ReadFile(exe + ".old"); string(b) != "old binary" {
		t.Fatalf("old exe not kept aside: %q", b)
	}
}

func TestBadSignatureIsRefused(t *testing.T) {
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	u, exe := setup(t, "linux", tarGz(t, "app", []byte("evil binary")), other)
	info, _ := u.Check(context.Background())
	err := u.Install(context.Background(), info)
	if err == nil || !strings.Contains(err.Error(), "signature") {
		t.Fatalf("want a signature error, got %v", err)
	}
	if b, _ := os.ReadFile(exe); string(b) != "old binary" {
		t.Fatal("binary was replaced despite a bad signature")
	}
}

func TestMacOSGetsALinkNotAnInstall(t *testing.T) {
	u, _ := setup(t, "darwin", zipOf(t, "app.app/x", []byte("x")), nil)
	info, err := u.Check(context.Background())
	if err != nil || !info.Available || info.CanInstall || info.PageURL == "" {
		t.Fatalf("%+v %v", info, err)
	}
}

func TestUpToDate(t *testing.T) {
	u, _ := setup(t, "linux", tarGz(t, "app", []byte("x")), nil)
	u.Version = "0.2.0"
	if info, err := u.Check(context.Background()); err != nil || info.Available {
		t.Fatalf("%+v %v", info, err)
	}
}

func TestUnsignedReleaseCannotInstall(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	srv := fakeGitHub(t, "v0.2.0", map[string][]byte{"app-0.2.0-linux-amd64.tar.gz": []byte("x")})
	u := &Updater{Repo: "me/app", Name: "app", Version: "0.1.0", API: srv.URL, Client: srv.Client(), Key: pub,
		Exe: func() (string, error) { return filepath.Join(t.TempDir(), "app"), nil }, GOOS: "linux"}
	info, _ := u.Check(context.Background())
	if !info.Available || info.CanInstall || !strings.Contains(info.Why, "signed") {
		t.Fatalf("%+v", info)
	}
}

func TestShippedKeyIsValid(t *testing.T) {
	if len(NewUpdater("a/b", "b", "0").Key) != ed25519.PublicKeySize {
		t.Fatal("updatePublicKey is missing or malformed")
	}
}
