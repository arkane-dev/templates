// release-sign makes the update signing key and signs release archives.
//
//	go run ./scripts/release-sign -keygen      one time: writes the private key, prints the public key
//	go run ./scripts/release-sign FILE...      writes FILE.sig next to each file
//
// The private key lives outside the repo, in ~/.config/cyberpunk-apps/release-signing.key
// (or $RELEASE_SIGNING_KEY). Back it up: lose it and existing installs can't verify new
// releases, so users would have to download the next version by hand.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func keyPath() string {
	if p := os.Getenv("RELEASE_SIGNING_KEY"); p != "" {
		return p
	}
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "cyberpunk-apps", "release-signing.key")
}

func main() {
	keygen := flag.Bool("keygen", false, "make a new key pair (refuses to overwrite one)")
	flag.Parse()
	path := keyPath()

	if *keygen {
		if _, err := os.Stat(path); err == nil {
			fail("a key already exists at %s. Delete it first if you mean to replace it.", path)
		}
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		check(err)
		check(os.MkdirAll(filepath.Dir(path), 0o700))
		check(os.WriteFile(path, []byte(base64.StdEncoding.EncodeToString(priv)+"\n"), 0o600))
		fmt.Printf("private key: %s (back it up)\npublic key:  %s\n", path, base64.StdEncoding.EncodeToString(pub))
		return
	}

	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		fail("no signing key at %s. Run with -keygen once.", path)
	}
	check(err)
	priv, err := base64.StdEncoding.DecodeString(string(trim(raw)))
	if err != nil || len(priv) != ed25519.PrivateKeySize {
		fail("%s is not a valid signing key", path)
	}
	if flag.NArg() == 0 {
		fail("usage: release-sign FILE...")
	}
	for _, f := range flag.Args() {
		data, err := os.ReadFile(f)
		check(err)
		sig := ed25519.Sign(ed25519.PrivateKey(priv), data)
		check(os.WriteFile(f+".sig", []byte(base64.StdEncoding.EncodeToString(sig)+"\n"), 0o644))
		fmt.Printf("signed %s\n", f)
	}
}

func trim(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r' || b[len(b)-1] == ' ') {
		b = b[:len(b)-1]
	}
	return b
}

func check(err error) {
	if err != nil {
		fail("%v", err)
	}
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "release-sign: "+format+"\n", a...)
	os.Exit(1)
}
