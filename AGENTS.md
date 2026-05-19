# AGENTS.md

## Project Overview

This repository contains `tapd`, a Go command line client for TAPD. It is a
single-module Go project built with Cobra and the typed
`github.com/go-tapd/tapd` SDK.

The executable entrypoint is `cmd/tapd/main.go`. Runtime concerns such as
configuration, authentication, prompting, and output helpers live in
`internal/app`. User-facing Cobra commands live in `internal/cmd`, generally one
file per TAPD resource area. User documentation lives in `docs/`, and command
coverage is tracked in `features.md`.

## Setup Commands

- Install dependencies: `go mod download`
- Tidy modules: `make go-mod-tidy`
- Install the CLI locally: `go install ./cmd/tapd`
- Show top-level help: `go run ./cmd/tapd --help`

The module declares Go `1.25.0` in `go.mod`. CI currently uses the latest Go
`1.26.x` series for lint compatibility checks.

## Development Workflow

- Run the CLI from source with `go run ./cmd/tapd <command>`.
- Build all packages with `go build ./...`.
- Use `tapd --help` and `tapd <command> --help` to verify command registration
  and flag help after changing commands.
- The generated user config is `~/.tapd/config.json`; do not commit local TAPD
  credentials or example files containing real tokens.

When adding or changing TAPD API commands:

- Confirm the matching typed service, request, and response already exist in
  `github.com/go-tapd/tapd` before editing the CLI.
- Reuse SDK typed requests and services instead of writing ad hoc HTTP calls.
- Prefer extending the existing resource file in `internal/cmd` over creating a
  new abstraction.
- Register new top-level resources in `internal/cmd/root.go`.
- Commands that require a project context should use `--workspace-id` / `-w`
  via `newWorkspaceFlag`.
- List/detail commands should support table output by default and `--format
  json` through `writeOutput`.
- Update `features.md`, `README.md`, and the relevant file in `docs/` when a
  user-facing command surface changes.

## Testing Instructions

- Run all tests: `make test`
- Run tests without Make: `go test ./... -race`
- Run a package directly: `go test ./internal/cmd`
- Run build verification: `go build ./...`

The project currently has little or no dedicated test coverage. For non-trivial
logic changes, add focused tests near the changed package before relying only on
manual command-help checks.

## Code Style

- Run lint before submitting: `make lint`
- Apply automatic lint fixes where safe: `make lint-fix`
- `make lint` runs `go mod tidy -compat=1.25.0` before `go tool
  golangci-lint run`.
- Keep imports formatted by Go tooling; the linter also enables `goimports` and
  `gofumpt` formatters.
- Keep command implementations narrow: parse flags, build an SDK request, call
  the service, map table rows, and delegate output to shared helpers.
- Avoid adding shared helpers until there is real reuse in multiple commands.
- Use existing helper patterns in `internal/cmd/common.go` for workspace flags,
  pagination, CSV parsing, strict multi-value parsing, and table/JSON output.

## Build and CI

- Local build: `go build ./...`
- Local lint: `make lint`
- Local test: `make test`
- GitHub Actions workflow: `.github/workflows/lint.yml`

The CI workflow currently runs lint and verifies that linting/tidying leaves a
clean working tree with `make check-clean-work`.

## Security Considerations

- Never commit real TAPD access tokens, client secrets, webhook payloads with
  private data, or generated `~/.tapd/config.json` contents.
- Prefer placeholders such as `$TAPD_ACCESS_TOKEN` and `<drawio-token>` in docs.
- Auth resolution supports `TAPD_ACCESS_TOKEN`, `TAPD_CLIENT_ID`,
  `TAPD_CLIENT_SECRET`, and local config. Keep examples clear about which mode is
  intended.

## Pull Request Guidelines

- Keep changes scoped to one resource area or one project-maintenance concern.
- Before opening a PR, run `make lint` and `make test` when the change touches Go
  code. Also run `go build ./...` for command or dependency changes.
- For documentation-only changes, inspect the rendered Markdown intent and check
  links to local docs files.
- Prefer conventional commit-style titles such as `feat(story): add ...`,
  `fix(auth): ...`, `docs(readme): ...`, or `chore(ci): ...`.

## Additional Notes

- `internal/x/skills/implement-api/SKILL.md` documents the local workflow for
  implementing new TAPD API commands.
- The SDK coverage ledger is in
  `https://github.com/go-tapd/tapd/blob/main/features.md`; this repository's
  `features.md` tracks CLI command coverage only.
- The CLI is intended to stay a thin wrapper over the SDK. If the SDK is missing
  an API, extend the SDK first instead of embedding HTTP behavior in this
  repository.
