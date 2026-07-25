// Package tipcatalog defines the shared tip/suggestion content catalog
// consumed by the bluefunda CLI, the cai-iOS app, and editor plugins — one
// schema, one set of tip content, rendered differently per surface.
//
// # Installation
//
//	go get github.com/bluefunda/tipcatalog
//
// # Usage
//
//	import tipcatalog "github.com/bluefunda/tipcatalog"
//
//	// Offline fallback baked into the binary via go:embed.
//	tips, err := tipcatalog.Embedded()
//
//	// Or load from a directory of tip JSON files.
//	tips, err := tipcatalog.LoadDir("tips")
//
//	// Verify a fetched manifest before trusting it.
//	ok := tipcatalog.Verify(manifestBytes, sig, tipcatalog.PublicKey)
//
// # Distribution
//
// On every GitHub Release, CI compiles tips/*.json into a single
// catalog.json, signs it with Ed25519, and attaches both catalog.json and
// catalog.json.sig to the release. Consumers fetch the latest release's
// assets, verify the signature against PublicKey, and fall back to the
// embedded copy on any failure.
package tipcatalog
