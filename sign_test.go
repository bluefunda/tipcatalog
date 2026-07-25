package tipcatalog

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestSignVerify_RoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	data := []byte(`{"catalog_version":"1"}`)

	sig := Sign(data, priv)
	if !Verify(data, sig, pub) {
		t.Fatal("expected valid signature to verify")
	}
}

func TestVerify_RejectsTamperedData(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	data := []byte(`{"catalog_version":"1"}`)
	sig := Sign(data, priv)

	tampered := append([]byte(nil), data...)
	tampered[0] = 'X'

	if Verify(tampered, sig, pub) {
		t.Fatal("expected tampered data to fail verification")
	}
}

func TestVerify_RejectsWrongKey(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	otherPub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	data := []byte(`{"catalog_version":"1"}`)
	sig := Sign(data, priv)

	if Verify(data, sig, otherPub) {
		t.Fatal("expected signature from a different key to fail verification")
	}
}

func TestPublicKey_Decodes(t *testing.T) {
	if len(PublicKey) != ed25519.PublicKeySize {
		t.Fatalf("PublicKey has length %d, want %d", len(PublicKey), ed25519.PublicKeySize)
	}
}
