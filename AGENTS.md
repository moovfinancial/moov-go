# AGENTS.md

Guidance for AI coding agents working in the `moov-go` repo. This is the single source of truth for agent instructions; vendor-specific tools point back here.

## Agent messages

Apply these rules only to messages posted to shared communication systems. Examples include GitHub PR descriptions and GitHub comments.

Do not apply these rules to direct user-agent conversations in a terminal, an IDE, or an agent application.

1. **Start with 🤖** as the first character.
2. **Write in [ASD-STE100 Simplified Technical English](https://www.asd-ste100.org/)** (Issue 9):
   - Use short sentences. Use a maximum of 20 words for instructions and 25 words for descriptions.
   - Write one instruction or one fact in each sentence.
   - Use the active voice. For instructions, use the imperative.
   - Do not use contractions. Write `do not`, not `don't`.
   - Do not omit articles (`a`, `the`) to make a sentence shorter.
   - Use American English spelling.
   - Use one word for one meaning. Do not use synonyms for the same thing.
   - Keep noun clusters to three words or fewer.
   - Use a vertical list when a sentence has many items or steps.
   - Domain terms (`Go`, `SDK`, `Moov`, `ACH`, `RTP`) are technical nouns. You can use them.
3. **Sign the last line** with the product you are running as:

   ```
   — Claude
   ```

   Use `Claude`, `Gemini`, `Copilot`, `Codex`, `Cursor`, `Grok`, or the equivalent name. Do not skip the emoji or the signature.

Commit messages are exempt. They follow the Conventional Commits format in the styleguide.

## What this repo is

`moov-go` is the official Go SDK for the [Moov payments API](https://docs.moov.io/api/). It is a library. It does not run a server.

Consumers import `github.com/moovfinancial/moov-go/pkg/moov` and call `moov.NewClient`.

## Layout

- `pkg/moov` — HTTP client, shared request helpers, and default API models.
- `pkg/mvYYMM` — version-pinned models and wrappers. They send `X-Moov-Version` for that release (for example `pkg/mv2607` uses `moov.Version2026_07`).
- `pkg/mhooks` — webhook signature checks. These tests do not need live API credentials.
- `examples/` — integration examples that hit the live Moov API.
- `internal/testtools` — sandbox account IDs and helpers for integration tests.

## Commands

```bash
make build                          # go build ./...
make check                          # download lint-project.sh, then lint and test
SKIP_TESTS=yes make check           # lint only
SKIP_LINTERS=yes make check         # tests only
go test ./pkg/mhooks/... ./pkg/mv2607/...   # unit tests that do not need credentials
go test ./...                       # full suite, including live API tests
go test -run TestName ./pkg/moov/...        # one test
go vet ./...
```

`make check` downloads `lint-project.sh` from `moov-io/infra`. Cover threshold is 30 percent (`COVER_THRESHOLD=30.0`).

## Credentials

Copy `secrets-template.env` to `secrets.env`. The Makefile loads `secrets.env` when that file exists.

- `MOOV_PUBLIC_KEY` / `MOOV_SECRET_KEY` — required for integration tests and examples
- `MOOV_HOST` — optional. Default is `api.moov.io`
- `PLAID_CLIENT_ID` / `PLAID_SECRET` — only for Plaid examples

Most tests in `pkg/moov` and `examples/` call the live API. Without credentials those tests fail. `pkg/mhooks` and httptest cases in `pkg/mv2607` do not need credentials.

## Versioned APIs

Add version-specific types and methods in `pkg/mvYYMM`, not by changing shared `pkg/moov` models in place.

- Put the versioned client wrapper next to the models for that version.
- Pass the matching `moov.Version` constant so requests set `X-Moov-Version`.
- Prefer httptest unit tests for versioned packages. Add live tests only when the sandbox must prove the call.

Exported `pkg/moov` changes are a public SDK contract. A renamed field or a removed method is a breaking change. Mark the PR title with `!`.

## Styleguide

See [.agents/styleguide.md](.agents/styleguide.md) for the PR title format (Conventional Commits variant) and Go conventions. Gemini Code Assist only reads `.gemini/styleguide.md`, so that path is a symlink to the single `.agents/styleguide.md` source — edit the `.agents/styleguide.md` file, never a separate copy.
