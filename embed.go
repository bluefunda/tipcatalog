package tipcatalog

import "embed"

//go:embed tips/*.json
var embeddedTips embed.FS

// Embedded returns the tip set baked into the binary at compile time — the
// offline fallback consumers fall back to when a fetched manifest is
// unavailable or fails signature verification.
func Embedded() ([]Tip, error) {
	return loadFS(embeddedTips, "tips")
}
