# AGENTS.md

## Project Overview

Shared Go library + JSON Schema defining the tip/suggestion content catalog consumed by the
`bai` CLI (`bluefunda-ai`), the cai-iOS app, and editor plugins. Single source of truth so all
three surfaces render the same tip content instead of maintaining separate copies. This is a
library, not a service.

Module: `github.com/bluefunda/tipcatalog` | Go 1.26+ | License: Apache 2.0

## Build and Test

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
```

All five must pass before any change is considered complete.

## Architecture

| File | Purpose |
|------|---------|
| `tip.go` | `Tip` struct — the schema, per Phase 0 of the Contextual Tip Engine spec |
| `topics.go` | `Topics` (the shared taxonomy), `EmbeddingFromDomainScope` — derives embeddings from `domain_scope`, no hand-authored floats |
| `validate.go` | `Validate([]Tip) error` — required fields, per-surface render keys, embedding dimensionality, recognized `domain_scope` topics, duplicate IDs |
| `load.go` | `LoadDir(dir) ([]Tip, error)`, `Compile([]Tip) ([]byte, error)` — derives each tip's embedding on load |
| `embed.go` | `//go:embed tips/*.json` + `Embedded() ([]Tip, error)` — offline fallback for consumers |
| `sign.go` | `Sign`/`Verify` — Ed25519 signing of the compiled catalog, stdlib only |
| `pubkey.go` | Checked-in Ed25519 public key used to verify a fetched manifest |
| `cmd/catalogtool` | Compiles + signs `tips/*.json` into `catalog.json`, used by CI on release |
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

Public, Apache 2.0, tagged via release-please — pkg.go.dev indexes it automatically once a
release is tagged and something does the first `go get`, same as `llmrouter`. No separate
publish step needed beyond that.

## Conventions

- Commits: conventional format (`feat:`, `fix:`, `chore:`)
- Branches: `<type>/<short-description>`
- PRs: squash-merged to `main`, title enforced as conventional-commit format by CI
- Releases: release-please + semver tags

## Adding or editing a tip

See [CONTENT_GUIDE.md](CONTENT_GUIDE.md) for the full content-authoring walkthrough. Tips are
tagged with `domain_scope` topics from `Topics` (`topics.go`) — never hand-author `embedding`,
it's derived automatically at load time. If the `Tip` struct's shape changes, update
`schema/tip.schema.json` to match — `schema_test.go` checks the schema stays in sync with both
the Go validator's required fields and `Topics`.
