# User Commands

Use `tapd user` commands to query TAPD user metadata.

All user commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

User commands require a workspace:

```bash
tapd user roles --workspace-id 123456
```

The short form is also supported:

```bash
tapd user roles -w 123456
```

## Output Formats

User commands default to table output:

```bash
tapd user roles -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd user roles -w 123456 --format json
```

## Roles

List role ID/name mappings for a workspace:

```bash
tapd user roles -w 123456
```

## Developer Notes

The user command implementation currently lives in:

```text
internal/cmd/user.go
```

When adding or renaming user commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
