---
name: implement-api
description: Implement or extend a TAPD API command in the go-tapd/cli repository. Use when adding commands like story list, bug list, task list, workspace view, or when wiring new TAPD API capabilities from the go-tapd/tapd SDK into the CLI.
---

# Implement API

This skill is for implementing a TAPD API command in the CLI repository:

- CLI repo: `https://github.com/go-tapd/cli`
- SDK repo: `https://github.com/go-tapd/tapd`

Follow this workflow.

## 1. Define the command contract

Before editing code, decide:

- command path
- required flags
- optional filters
- default output columns
- whether `--format json` should be supported

Prefer CLI behavior that matches existing commands in the repository.

## 2. Inspect the SDK first

Do not write request logic from memory.

Find the matching TAPD SDK service, request, and response in `go-tapd/tapd`.

Check:

- which service method to call
- which request struct to use
- which fields are required
- which response fields are suitable for table output

If the SDK already supports the API, reuse it directly.
If the SDK is missing support for the API, stop and note that the SDK must be extended first.

## 3. Reuse existing CLI patterns

Follow the structure already used in `go-tapd/cli`.

Look for:

- root command registration
- shared flag helpers
- output helpers
- the closest existing command of the same shape

Use the existing style for:

- workspace flag handling
- request construction
- table output
- json output
- command help text

## 4. Implement with minimal abstraction

Add only what the new command needs:

- flags
- SDK request mapping
- service call
- table row mapping
- command registration

Do not introduce new abstractions unless there is already repeated logic.

## 5. Keep output consistent

For list commands:

- support `table` and `json`
- choose a small default column set
- prefer stable identifiers and operator-relevant fields

For single-resource commands:

- return the raw object for `json`
- use a compact table or summary view for default output

## 6. Validate before finishing

Run:

- `go build ./...`
- `go test ./...`
- `go run . <command> --help`

If the repo is in an offline or sandboxed environment, adapt the validation command accordingly, but still verify compile success and command help output.

## 7. Update docs only when useful

If the command is user-facing and materially expands CLI capability, update the repository documentation with a short usage example.

## Decision rules

- Prefer extending an existing command file over creating a new one when the resource already exists.
- Prefer adding shared helpers only after the second real reuse.
- Prefer SDK-backed typed requests over ad hoc HTTP logic.
- Prefer narrow, incremental commands over large multi-command batches in one change.

## Expected implementation checklist

- command contract is clear
- SDK method is confirmed
- flags map correctly to SDK request fields
- root registration is complete
- `table|json` output works
- build passes
- tests pass
- help output is checked
