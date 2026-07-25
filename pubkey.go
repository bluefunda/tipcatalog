package tipcatalog

import (
	"crypto/ed25519"
	"encoding/base64"
)

// publicKeyB64 is the Ed25519 public key used to verify a fetched manifest's
// signature. Its matching private key is held only in the tipcatalog repo's
// TIP_CATALOG_SIGNING_KEY secret and used by the publish-manifest workflow —
// it never appears in source control.
const publicKeyB64 = "M1U752fzmcXdY+7L+NlNdsHBvVE4o/S41CuDiHPfyBA="

// PublicKey is the Ed25519 public key that Verify checks manifest signatures
// against. It is decoded once at package init from publicKeyB64.
var PublicKey ed25519.PublicKey = mustDecodeKey(publicKeyB64)

func mustDecodeKey(b64 string) ed25519.PublicKey {
	key, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic("tipcatalog: invalid embedded public key: " + err.Error())
	}
	return ed25519.PublicKey(key)
}
