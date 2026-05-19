# Source Commands

Use `tapd source` commands to add and query source commit information linked to
TAPD work items.

All source commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most source commands require a workspace:

```bash
tapd source commit list --workspace-id 123456 --entity-type story --object-id 1111112222001000001
```

The short form is also supported:

```bash
tapd source commit list -w 123456 --entity-type story --object-id 1111112222001000001
```

## Output Formats

Source commands default to table output:

```bash
tapd source commit list -w 123456 --entity-type story --object-id 1111112222001000001
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd source commit list -w 123456 --entity-type story --object-id 1111112222001000001 --format json
```

## Add Commit Info

```bash
tapd source commit add \
  -w 123456 \
  --commit-id abc123 \
  --author alice \
  --message "Implement login" \
  --files "internal/auth.go,README.md" \
  --repo go-tapd-cli \
  --repo-id go-tapd-cli \
  --commit-time "2026-06-01 10:00:00"
```

Useful optional flags:

```bash
tapd source commit add \
  -w 123456 \
  --commit-id abc123 \
  --author alice \
  --message "Implement login" \
  --files "internal/auth.go,README.md" \
  --repo go-tapd-cli \
  --repo-id go-tapd-cli \
  --commit-time "2026-06-01 10:00:00" \
  --git-env github \
  --repo-url "https://github.com/go-tapd/cli" \
  --commit-url "https://github.com/go-tapd/cli/commit/abc123"
```

## List Commit Info

```bash
tapd source commit list \
  -w 123456 \
  --entity-type story \
  --object-id 1111112222001000001
```

Supported entity types are `story`, `bug`, and `task`.

Useful optional flags:

```bash
tapd source commit list \
  -w 123456 \
  --entity-type bug \
  --object-id 1111112222001000001 \
  --related-type source_code \
  --commit-time "2026-01-01~2026-01-31"
```

## Commit Objects

List TAPD objects related to commits:

```bash
tapd source commit objects \
  -w 123456 \
  --commit-ids abc123,def456 \
  --entity-type story
```

Useful optional flags:

```bash
tapd source commit objects \
  -w 123456 \
  --commit-ids abc123 \
  --entity-type task \
  --scm-type github \
  --fields id,name,status
```

## Developer Notes

The source command implementation currently lives in:

```text
internal/cmd/source.go
```

When adding or renaming source commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
