# Contributing to tipcatalog

## Prerequisites

- Go 1.26+
- `golangci-lint` ([install guide](https://golangci-lint.run/welcome/install/))

## Getting started

```bash
git clone https://github.com/bluefunda/tipcatalog
cd tipcatalog
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
5. Open a pull request against `main` — the PR title must also follow conventional commit format (enforced by CI)

## Adding or editing a tip

See [CONTENT_GUIDE.md](CONTENT_GUIDE.md) for the full walkthrough (topic tagging, validation,
what happens on merge). If you change the `Tip` struct's shape itself, also update
`schema/tip.schema.json` to match — `schema_test.go` checks the two stay in sync.

## License

By contributing you agree that your contributions will be licensed under the Apache 2.0 License.
