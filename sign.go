package tipcatalog

import "crypto/ed25519"

// Sign returns an Ed25519 signature over data using priv. Used by the
// publish-manifest CI workflow to sign the compiled catalog.json.
func Sign(data []byte, priv ed25519.PrivateKey) []byte {
	return ed25519.Sign(priv, data)
}

// Verify reports whether sig is a valid Ed25519 signature over data for
// pub. Consumers should call this (with PublicKey) before trusting a
// fetched manifest.
func Verify(data, sig []byte, pub ed25519.PublicKey) bool {
	return ed25519.Verify(pub, data, sig)
}
