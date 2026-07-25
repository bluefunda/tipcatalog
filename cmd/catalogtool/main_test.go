package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	tipcatalog "github.com/bluefunda/tipcatalog"
)

func TestRun_CompilesAndSigns(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "catalog.json")
	sigOut := filepath.Join(dir, "catalog.json.sig")

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	t.Setenv("TIP_CATALOG_SIGNING_KEY", base64.StdEncoding.EncodeToString(priv.Seed()))

	if err := run("../../tips", out, sigOut); err != nil {
		t.Fatalf("run: %v", err)
	}

	compiled, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read catalog: %v", err)
	}
	var tips []tipcatalog.Tip
	if err := json.Unmarshal(compiled, &tips); err != nil {
		t.Fatalf("unmarshal catalog: %v", err)
	}
	if len(tips) == 0 {
		t.Fatal("expected at least one tip in compiled catalog")
	}

	sigB64, err := os.ReadFile(sigOut)
	if err != nil {
		t.Fatalf("read signature: %v", err)
	}
	sig, err := base64.StdEncoding.DecodeString(string(sigB64))
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	if !tipcatalog.Verify(compiled, sig, priv.Public().(ed25519.PublicKey)) {
		t.Fatal("expected signature to verify against the catalog bytes")
	}
}

func TestRun_RequiresSigningKeyForSigOut(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TIP_CATALOG_SIGNING_KEY", "")
	err := run("../../tips", filepath.Join(dir, "catalog.json"), filepath.Join(dir, "catalog.json.sig"))
	if err == nil {
		t.Fatal("expected error when TIP_CATALOG_SIGNING_KEY is unset")
	}
}
