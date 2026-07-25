# Contributing to tip-catalog

## Prerequisites

- Go 1.25+
- `golangci-lint` ([install guide](https://golangci-lint.run/welcome/install/))

## Getting started

```bash
git clone https://github.com/bluefunda/tip-catalog
cd tip-catalog
go build ./...
go test -race ./...
```

## Making changes

1. Fork the repo and create a branch: `git checkout -b feat/my-change`
2. Make your changes
3. Ensure all checks pass:
   ```bash
   go build ./...
   go test -race ./...
   go vet ./...
   golangci-lint run
   ```
4. Commit using [Conventional Commits](https://www.conventionalcommits.org/): `fix:`, `feat:`, `chore:`, etc.
5. Open a pull request against `main`

## Adding or editing a tip

1. Add or edit a JSON file under `tips/` (one tip per file, filename matches the tip `id`)
2. Run `go test ./...` — the loader validates every tip in `tips/` against required fields,
   per-surface render copy, embedding dimensionality, and duplicate IDs
3. If you change the `Tip` struct's shape, update `schema/tip.schema.json` to match — `schema_test.go`
   checks the schema's `required` list stays in sync with the Go validator

## License

By contributing you agree that your contributions will be licensed under the Apache 2.0 License.
