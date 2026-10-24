package main

// Auto-update from GitHub Releases.
//
// Check asks the releases API for the latest release and compares versions. Install downloads
// this platform's archive and its .sig, verifies the ed25519 signature against the public key
// built into the app, and swaps the new binary in. Nothing runs unless the signature checks out.
//
// Linux and Windows update in place. macOS gets the notice and a link: its build is made in CI,
// which doesn't hold the signing key (README "Release").

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// updatePublicKey verifies release signatures. scripts/release-sign makes the key pair:
// the private key stays on the release machine, this public half ships in every build.
const updatePublicKey = "hYPppbg2VdhFRp4DpOvGP9Y+VB3x4PVLMm5KEHKoFxo="

const (
	maxDownload   = 500 << 20 // refuse anything bigger than 500 MB
	checkInterval = 24 * time.Hour
)

// UpdateInfo is what the frontend shows about an available release.
type UpdateInfo struct {
	Current    string `json:"current"`
	Latest     string `json:"latest"`
	Available  bool   `json:"available"`
	CanInstall bool   `json:"canInstall"` // false: show the link instead (macOS, or no write access)
	Why        string `json:"why"`        // why it can't install here, if it can't
	Notes      string `json:"notes"`
	PageURL    string `json:"pageUrl"`
	asset      release
}

type release struct {
	URL, SigURL, Name string
	Size              int64
}

type Updater struct {
	Repo    string // owner/name on GitHub
	Name    string // archive and binary name, e.g. "genai-writer"
	Version string
	API     string // https://api.github.com, swapped in tests
	Client  *http.Client
	Key     ed25519.PublicKey
	Exe     func() (string, error)
	GOOS    string
}

func NewUpdater(repo, name, version string) *Updater {
	key, _ := base64.StdEncoding.DecodeString(updatePublicKey)
	api := "https://api.github.com"
	if v := os.Getenv("UPDATE_API_BASE"); v != "" {
		api = v // for testing against a local server. Downloads still have to pass the signature check.
	}
	return &Updater{
		Repo: repo, Name: name, Version: version, API: api,
		Client: &http.Client{Timeout: 60 * time.Second}, Key: ed25519.PublicKey(key),
		Exe: os.Executable, GOOS: runtime.GOOS,
	}
}

// archiveName follows scripts/dist.sh and the macOS workflow.
func (u *Updater) archiveName(version string) string {
	switch u.GOOS {
	case "windows":
		return fmt.Sprintf("%s-%s-windows-amd64.zip", u.Name, version)
	case "darwin":
		return fmt.Sprintf("%s-%s-macos-universal.zip", u.Name, version)
	default:
		return fmt.Sprintf("%s-%s-linux-amd64.tar.gz", u.Name, version)
	}
}

// Newer reports whether version a is later than b. Both look like 1.2.3, with an optional leading v.
func Newer(a, b string) bool {
	pa, pb := parseVersion(a), parseVersion(b)
	for i := range 3 {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

// parseVersion reads 1.2.3 (and v1.2.3, 1.2.3-beta). Parts that aren't numbers count as 0.
func parseVersion(v string) [3]int {
	var out [3]int
	core := strings.TrimPrefix(strings.TrimSpace(v), "v")
	if i := strings.IndexAny(core, "-+"); i >= 0 {
		core = core[:i]
	}
	for i, p := range strings.SplitN(core, ".", 3) {
		out[i], _ = strconv.Atoi(p)
	}
	return out
}

func (u *Updater) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", u.Name+"/"+u.Version)
	res, err := u.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("%s: %s", url, res.Status)
	}
	return res, nil
}

// Check asks GitHub for the latest release.
func (u *Updater) Check(ctx context.Context) (UpdateInfo, error) {
	info := UpdateInfo{Current: u.Version}
	res, err := u.get(ctx, fmt.Sprintf("%s/repos/%s/releases/latest", u.API, u.Repo))
	if err != nil {
		return info, fmt.Errorf("could not reach GitHub: %w", err)
	}
	defer res.Body.Close()
	var rel struct {
		Tag     string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Body    string `json:"body"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
			Size int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&rel); err != nil {
		return info, fmt.Errorf("could not read the release: %w", err)
	}
	info.Latest = strings.TrimPrefix(rel.Tag, "v")
	info.PageURL = rel.HTMLURL
	info.Notes = rel.Body
	info.Available = Newer(info.Latest, u.Version)
	if !info.Available {
		return info, nil
	}
	want := u.archiveName(info.Latest)
	for _, a := range rel.Assets {
		switch a.Name {
		case want:
			info.asset.URL, info.asset.Size, info.asset.Name = a.URL, a.Size, a.Name
		case want + ".sig":
			info.asset.SigURL = a.URL
		}
	}
	info.CanInstall, info.Why = u.canInstall(info.asset)
	return info, nil
}

func (u *Updater) canInstall(a release) (bool, string) {
	switch {
	case u.GOOS == "darwin":
		return false, "Updates on macOS are a download for now."
	case len(u.Key) != ed25519.PublicKeySize:
		return false, "This build has no update key."
	case a.URL == "" || a.SigURL == "":
		return false, "This release has no signed download for your system."
	case a.Size > maxDownload:
		return false, "The download is too large."
	}
	exe, err := u.Exe()
	if err != nil {
		return false, "Could not find the running app."
	}
	if f, err := os.CreateTemp(filepath.Dir(exe), ".update-*"); err != nil {
		return false, fmt.Sprintf("No write access to %s. Download the new version instead.", filepath.Dir(exe))
	} else {
		f.Close()
		os.Remove(f.Name())
	}
	return true, ""
}

func (u *Updater) download(ctx context.Context, url string, limit int64) ([]byte, error) {
	res, err := u.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, errors.New("download is larger than expected")
	}
	return b, nil
}

// Install downloads, verifies and swaps in the new binary. The app keeps running the old
// one until it restarts.
func (u *Updater) Install(ctx context.Context, info UpdateInfo) error {
	if ok, why := u.canInstall(info.asset); !ok {
		return errors.New(why)
	}
	archive, err := u.download(ctx, info.asset.URL, maxDownload)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	sigText, err := u.download(ctx, info.asset.SigURL, 4096)
	if err != nil {
		return fmt.Errorf("signature download failed: %w", err)
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(sigText)))
	if err != nil || !ed25519.Verify(u.Key, archive, sig) {
		return errors.New("the download's signature is not valid, so it was not installed")
	}
	bin, err := u.extract(archive)
	if err != nil {
		return err
	}
	exe, err := u.Exe()
	if err != nil {
		return err
	}
	return replaceExecutable(exe, bin, u.GOOS)
}

// extract pulls the app binary out of the release archive.
func (u *Updater) extract(archive []byte) ([]byte, error) {
	want := u.Name
	if u.GOOS == "windows" {
		want += ".exe"
		zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
		if err != nil {
			return nil, err
		}
		for _, f := range zr.File {
			if filepath.Base(f.Name) == want {
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				return io.ReadAll(io.LimitReader(rc, maxDownload))
			}
		}
		return nil, fmt.Errorf("%s not found in the download", want)
	}
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("%s not found in the download", want)
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag == tar.TypeReg && filepath.Base(h.Name) == want {
			return io.ReadAll(io.LimitReader(tr, maxDownload))
		}
	}
}

// replaceExecutable writes the new binary next to the old one, then renames it into place.
// Windows won't overwrite a running .exe, but it will rename one, so the old file moves aside first.
func replaceExecutable(exe string, bin []byte, goos string) error {
	dir := filepath.Dir(exe)
	tmp, err := os.CreateTemp(dir, ".update-*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	tmp.Close()
	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if goos == "windows" {
		old := exe + ".old"
		os.Remove(old)
		if err := os.Rename(exe, old); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		if err := os.Rename(tmp.Name(), exe); err != nil {
			os.Rename(old, exe) // put the old one back
			return err
		}
		return nil
	}
	return os.Rename(tmp.Name(), exe)
}

// CleanupOldBinary removes the .old file a Windows update leaves behind.
func CleanupOldBinary() {
	if exe, err := os.Executable(); err == nil {
		os.Remove(exe + ".old")
	}
}

// Relaunch starts the (new) binary with the same arguments. The caller then quits.
func Relaunch() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Start()
}

// DueForCheck reports whether a day has passed since the last automatic check.
func DueForCheck(last string, now time.Time) bool {
	t, err := time.Parse(time.RFC3339, last)
	return err != nil || now.Sub(t) >= checkInterval
}
