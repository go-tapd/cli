# tapd CLI

`tapd` is a TAPD command line client built on top of [`github.com/go-tapd/tapd`](https://github.com/go-tapd/tapd).

## Install

```bash
go install github.com/go-tapd/cli/cmd/tapd@latest
```

This installs a binary named `tapd`.

## Current Commands

- `tapd login`
- `tapd auth status`
- `tapd auth logout`
- `tapd workspace view`
- `tapd workspace users`
- `tapd story list`
- `tapd bug list`
- `tapd task list`

## Login

PAT:

```bash
go run ./cmd/tapd login --auth-method pat
```

Validate at login time with a workspace:

```bash
go run ./cmd/tapd login --auth-method pat --workspace-id 123456
```

PAT with flag:

```bash
go run ./cmd/tapd login --token "$TAPD_ACCESS_TOKEN"
```

Basic Authentication:

```bash
go run ./cmd/tapd login --auth-method basic
```

Basic Authentication with flags:

```bash
go run ./cmd/tapd login --client-id "$TAPD_CLIENT_ID" --client-secret "$TAPD_CLIENT_SECRET"
```

The login command validates credentials before saving them to:

```text
~/.tapd/config.json
```

By default, `login` stores the credentials directly. If you want an API-side validation during login, pass `--workspace-id`.

You can also skip local config and use environment variables directly:

```bash
export TAPD_ACCESS_TOKEN=...
```

or

```bash
export TAPD_CLIENT_ID=...
export TAPD_CLIENT_SECRET=...
```

## Examples

```bash
go run ./cmd/tapd auth status
go run ./cmd/tapd workspace view --workspace-id 123456
go run ./cmd/tapd workspace users --workspace-id 123456
go run ./cmd/tapd story list --workspace-id 123456 --limit 20
go run ./cmd/tapd bug list --workspace-id 123456 --owner alice
go run ./cmd/tapd task list --workspace-id 123456 --creator bob --format json
```
