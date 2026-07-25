# CLAUDE.md — tip-catalog

## What is this?

Shared Go library + JSON Schema defining the tip/suggestion content catalog consumed by the
`bai` CLI (`bluefunda-ai`), the cai-iOS app, and editor plugins. Single source of truth so all
three surfaces render the same tip content instead of maintaining separate copies.

Module: `github.com/bluefunda/tip-catalog` | Go 1.25+ | License: Apache 2.0

## Build & Verify

```bash
go build ./...
go test -race ./...
go vet ./...
golangci-lint run
```

All four must pass before any change is considered complete.

## Architecture

| File | Purpose |
|------|---------|
| `tip.go` | `Tip` struct — the schema, per Phase 0 of the Contextual Tip Engine spec |
| `validate.go` | `Validate([]Tip) error` — required fields, per-surface render keys, embedding dimensionality, duplicate IDs |
| `load.go` | `LoadDir(dir) ([]Tip, error)`, `Compile([]Tip) ([]byte, error)` |
| `embed.go` | `//go:embed tips/*.json` + `Embedded() ([]Tip, error)` — offline fallback for consumers |
| `sign.go` | `Sign`/`Verify` — Ed25519 signing of the compiled catalog, stdlib only |
| `pubkey.go` | Checked-in Ed25519 public key used to verify a fetched manifest |
| `tips/*.json` | One JSON file per tip, filename matches `id` |
| `schema/tip.schema.json` | JSON Schema (2020-12) mirroring `Tip`, for the Swift-side codegen |

No third-party dependencies — everything uses stdlib (`encoding/json`, `crypto/ed25519`, `embed`).
Keep it that way unless there's a clear reason to add one.

## Distribution

CI (`.github/workflows/publish-manifest.yml`) compiles `tips/*.json` into a single `catalog.json`
on every GitHub Release, signs it with an Ed25519 key held in the `TIP_CATALOG_SIGNING_KEY` repo
secret, and uploads both `catalog.json` and `catalog.json.sig` as release assets. Consumers fetch
the latest release's assets, verify the signature with the public key in `pubkey.go`, and fall
back to `Embedded()` on any failure (missing network, bad signature, etc.).

## Conventions

- Commits: conventional format (`feat:`, `fix:`, `chore:`)
- Branches: `<type>/<short-description>`
- PRs: squash-merged to `main`
- Releases: release-please + semver tags

## Adding or editing a tip

1. Add/edit a JSON file under `tips/`
2. Run `go test ./...` — validates required fields, per-surface render copy, embedding
   dimensionality (`EmbeddingDim` in `validate.go`), and duplicate IDs
3. If the `Tip` struct's shape changes, update `schema/tip.schema.json` to match —
   `schema_test.go` checks the schema's `required` list stays in sync with the Go validator
