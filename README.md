# tapd CLI

`tapd` is a TAPD command line client built on top of [`github.com/go-tapd/tapd`](https://github.com/go-tapd/tapd).

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
go run . login --auth-method pat
```

PAT with flag:

```bash
go run . login --token "$TAPD_ACCESS_TOKEN"
```

Basic Authentication:

```bash
go run . login --auth-method basic
```

Basic Authentication with flags:

```bash
go run . login --client-id "$TAPD_CLIENT_ID" --client-secret "$TAPD_CLIENT_SECRET"
```

The login command validates credentials before saving them to:

```text
~/.config/tapd/config.json
```

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
go run . auth status
go run . workspace view --workspace-id 123456
go run . workspace users --workspace-id 123456
go run . story list --workspace-id 123456 --limit 20
go run . bug list --workspace-id 123456 --owner alice
go run . task list --workspace-id 123456 --creator bob --format json
```
