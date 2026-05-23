# AGENTS.md

## Cursor Cloud specific instructions

### Overview

`moov-go` is the official Go SDK client for the Moov payments API. It is a **library** (no runnable server).

### Build / Lint / Test

| Action | Command |
|--------|---------|
| Build | `go build ./...` |
| Lint | `make check` (with `SKIP_TESTS=yes` to lint only) |
| Test (full, requires credentials) | `go test ./...` |
| Unit tests (no credentials needed) | `go test ./pkg/mhooks/... ./pkg/mv2607/...` |
| Vet | `go vet ./...` |

- `make check` downloads `lint-project.sh` from moov-io/infra and runs golangci-lint + tests.
- Set `SKIP_TESTS=yes` to run only linters, or `SKIP_LINTERS=yes` to run only tests.
- Almost all tests are integration tests that require live Moov API credentials (`MOOV_PUBLIC_KEY`, `MOOV_SECRET_KEY`). Without them, only `pkg/mhooks` and `pkg/mv2607` tests pass.

### Credentials

Set in environment or in `secrets.env` (auto-loaded by Makefile):
- `MOOV_PUBLIC_KEY` / `MOOV_SECRET_KEY` — required for integration tests
- `PLAID_CLIENT_ID` / `PLAID_SECRET` — only for Plaid-specific examples

### Important notes

- Go 1.26+ is required (toolchain go1.26.3). The VM's default Go is too old; `/usr/local/go/bin` must be on PATH.
- `moov-money-dev` depends on this repo via a `replace` directive; both repos must remain sibling directories.
