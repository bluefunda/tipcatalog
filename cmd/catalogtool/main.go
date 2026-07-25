// Command catalogtool compiles tips/*.json into a single catalog.json and,
// optionally, signs it with the Ed25519 key in TIP_CATALOG_SIGNING_KEY.
// Used by .github/workflows/publish-manifest.yml on every GitHub Release.
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	tipcatalog "github.com/bluefunda/tip-catalog"
)

func main() {
	tipsDir := flag.String("tips", "tips", "directory of tip JSON files")
	out := flag.String("out", "catalog.json", "output compiled catalog path")
	sigOut := flag.String("sig-out", "", "output signature path; if set, signs with TIP_CATALOG_SIGNING_KEY")
	flag.Parse()

	if err := run(*tipsDir, *out, *sigOut); err != nil {
		fmt.Fprintln(os.Stderr, "catalogtool:", err)
		os.Exit(1)
	}
}

func run(tipsDir, out, sigOut string) error {
	tips, err := tipcatalog.LoadDir(tipsDir)
	if err != nil {
		return fmt.Errorf("load tips: %w", err)
	}

	compiled, err := tipcatalog.Compile(tips)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	if err := os.WriteFile(out, compiled, 0o644); err != nil {
		return fmt.Errorf("write catalog: %w", err)
	}

	if sigOut == "" {
		return nil
	}

	seedB64 := os.Getenv("TIP_CATALOG_SIGNING_KEY")
	if seedB64 == "" {
		return fmt.Errorf("-sig-out requires TIP_CATALOG_SIGNING_KEY to be set")
	}
	seed, err := base64.StdEncoding.DecodeString(seedB64)
	if err != nil {
		return fmt.Errorf("decode TIP_CATALOG_SIGNING_KEY: %w", err)
	}
	if len(seed) != ed25519.SeedSize {
		return fmt.Errorf("TIP_CATALOG_SIGNING_KEY must decode to %d bytes, got %d", ed25519.SeedSize, len(seed))
	}

	priv := ed25519.NewKeyFromSeed(seed)
	sig := tipcatalog.Sign(compiled, priv)
	sigB64 := base64.StdEncoding.EncodeToString(sig)
	if err := os.WriteFile(sigOut, []byte(sigB64), 0o644); err != nil {
		return fmt.Errorf("write signature: %w", err)
	}
	return nil
}
