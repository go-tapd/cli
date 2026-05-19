# Measure Commands

Use `tapd measure` commands to query TAPD measurement data.

All measure commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Measure commands require a workspace:

```bash
tapd measure life-times --workspace-id 123456 --entity-type story --entity-id 1111112222001000001
```

The short form is also supported:

```bash
tapd measure life-times -w 123456 --entity-type story --entity-id 1111112222001000001
```

## Output Formats

Measure commands default to table output:

```bash
tapd measure life-times -w 123456 --entity-type story --entity-id 1111112222001000001
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd measure life-times -w 123456 --entity-type story --entity-id 1111112222001000001 --format json
```

## Life Times

Query status transition life times for a story, task, or bug:

```bash
tapd measure life-times \
  -w 123456 \
  --entity-type story \
  --entity-id 1111112222001000001
```

Supported entity types are `story`, `task`, and `bug`.

Useful optional flags:

```bash
tapd measure life-times \
  -w 123456 \
  --entity-type bug \
  --entity-id 1111112222001000001 \
  --created "2026-01-01~2026-01-31" \
  --limit 50 \
  --page 1
```

Request specific fields from TAPD:

```bash
tapd measure life-times -w 123456 --entity-type task --entity-id 1111112222001000001 --fields id,status,time_cost,created
```

## Developer Notes

The measure command implementation currently lives in:

```text
internal/cmd/measure.go
```

When adding or renaming measure commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
