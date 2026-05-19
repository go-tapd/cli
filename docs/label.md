# Label Commands

Use `tapd label` commands to create, list, count, and update TAPD labels.

All label commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most label commands require a workspace:

```bash
tapd label list --workspace-id 123456
```

The short form is also supported:

```bash
tapd label list -w 123456
```

## Output Formats

Label commands default to table output:

```bash
tapd label list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd label list -w 123456 --format json
```

## Create a Label

```bash
tapd label create -w 123456 --name backend
```

Useful optional flags:

```bash
tapd label create -w 123456 --name backend --color 1 --creator alice
```

TAPD label colors use numeric values `1`, `2`, `3`, and `4`.

## List Labels

```bash
tapd label list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd label list -w 123456 --ids 10001,10002
tapd label list -w 123456 --name backend
tapd label list -w 123456 --creator alice
tapd label list -w 123456 --created "2026-01-01~2026-01-31"
```

## Count Labels

```bash
tapd label count -w 123456
tapd label count -w 123456 --creator alice
tapd label count -w 123456 --name backend
```

## Update a Label

```bash
tapd label update 10001 -w 123456 --color 2
```

Set the modifier when updating:

```bash
tapd label update 10001 -w 123456 --color 2 --modifier alice
```

`tapd label update` requires at least one of `--color` or `--modifier`.

## Developer Notes

The label command implementation currently lives in:

```text
internal/cmd/label.go
```

When adding or renaming label commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
