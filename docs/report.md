# Report Commands

Use `tapd report` commands to query TAPD workspace reports.

All report commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Report commands require a workspace:

```bash
tapd report list --workspace-id 123456
```

The short form is also supported:

```bash
tapd report list -w 123456
```

## Output Formats

Report commands default to table output:

```bash
tapd report list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd report list -w 123456 --format json
```

## List Reports

```bash
tapd report list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd report list -w 123456 --id 10001
tapd report list -w 123456 --title "Weekly"
tapd report list -w 123456 --author alice
tapd report list -w 123456 --created "2026-01-01~2026-01-31"
```

Request specific fields from TAPD:

```bash
tapd report list -w 123456 --fields id,title,report_type,status,author,created
```

## Developer Notes

The report command implementation currently lives in:

```text
internal/cmd/report.go
```

When adding or renaming report commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
