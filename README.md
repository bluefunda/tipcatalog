# tipcatalog

[![Go Reference](https://pkg.go.dev/badge/github.com/bluefunda/tipcatalog.svg)](https://pkg.go.dev/github.com/bluefunda/tipcatalog)
[![CI](https://github.com/bluefunda/tipcatalog/actions/workflows/ci.yml/badge.svg)](https://github.com/bluefunda/tipcatalog/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluefunda/tipcatalog)](https://goreportcard.com/report/github.com/bluefunda/tipcatalog)

Shared tip/suggestion content catalog for the bluefunda CLI (`bai`), the cai-iOS app, and editor
plugins. One schema, one set of tip content, rendered differently per surface.

## Why

Each client (CLI, iOS, editor) needs contextual tips, but authoring and maintaining separate copy
per surface drifts fast. `tipcatalog` is the single source of truth: a `Tip` schema with
per-surface render copy, a validator, and a signed, versioned distribution format so clients can
fetch updates without redeploying.

## Installation

```bash
go get github.com/bluefunda/tipcatalog
```

## Usage

```go
import tipcatalog "github.com/bluefunda/tipcatalog"

// Offline fallback baked into the binary via go:embed.
tips, err := tipcatalog.Embedded()

// Or load from a directory of tip JSON files (e.g. this repo's tips/ during CI).
tips, err := tipcatalog.LoadDir("tips")

// Verify a fetched manifest before trusting it.
ok := tipcatalog.Verify(manifestBytes, sig, tipcatalog.PublicKey)
```

## Schema

Each tip is one JSON file under `tips/`, validated against the fields documented in
[`schema/tip.schema.json`](schema/tip.schema.json) (used by the Swift side to codegen matching
types). See [`tip.go`](tip.go) for the canonical Go definition. To add or edit tip content, see
[CONTENT_GUIDE.md](CONTENT_GUIDE.md).

## Distribution

On every GitHub Release, CI compiles `tips/*.json` into a single `catalog.json`, signs it with
Ed25519, and attaches both `catalog.json` and `catalog.json.sig` to the release. Clients fetch the
latest release's assets, verify the signature against `PublicKey` (in [`pubkey.go`](pubkey.go)),
and fall back to the embedded copy on any failure.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache 2.0 — see [LICENSE](LICENSE).
