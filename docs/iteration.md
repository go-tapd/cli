# Iteration Commands

Use `tapd iteration` commands to create, inspect, list, update, and query TAPD
iterations.

All iteration commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most iteration commands require a workspace:

```bash
tapd iteration list --workspace-id 123456
```

The short form is also supported:

```bash
tapd iteration list -w 123456
```

## Output Formats

Iteration commands default to table output:

```bash
tapd iteration list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd iteration list -w 123456 --format json
```

## Common Workflows

### Create an Iteration

```bash
tapd iteration create \
  -w 123456 \
  --name "Sprint 1" \
  --description "Sprint goal" \
  --start-date 2026-06-01 \
  --end-date 2026-06-14 \
  --creator alice
```

Useful optional flags:

```bash
tapd iteration create \
  -w 123456 \
  --name "Sprint 1" \
  --description "Sprint goal" \
  --start-date 2026-06-01 \
  --end-date 2026-06-14 \
  --creator alice \
  --status open \
  --label "backend,release" \
  --workitem-type-id 10001 \
  --plan-app-id 20001
```

### View an Iteration

```bash
tapd iteration view 1111112222001000001 -w 123456
```

This command queries the iteration list API with the requested ID and returns
the first matching iteration.

### List Iterations

```bash
tapd iteration list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd iteration list -w 123456 --status open
tapd iteration list -w 123456 --creator alice
tapd iteration list -w 123456 --name "Sprint"
tapd iteration list -w 123456 --ids 1111112222001000001,1111112222001000002
```

Request specific fields from TAPD:

```bash
tapd iteration list -w 123456 --fields id,name,status,startdate,enddate,creator,modified
```

### Count Iterations

```bash
tapd iteration count -w 123456
tapd iteration count -w 123456 --status open
tapd iteration count -w 123456 --creator alice
```

### Update an Iteration

`current-user` is required by the TAPD update API:

```bash
tapd iteration update 1111112222001000001 \
  -w 123456 \
  --current-user alice \
  --name "Sprint 1.1"
```

Update several common fields:

```bash
tapd iteration update 1111112222001000001 \
  -w 123456 \
  --current-user alice \
  --status done \
  --end-date 2026-06-21 \
  --label "backend,release"
```

## Metadata Commands

### Fields

List iteration custom field settings:

```bash
tapd iteration fields -w 123456
```

Use JSON output to inspect full custom field options:

```bash
tapd iteration fields -w 123456 --format json
```

### Workitem Types

```bash
tapd iteration workitem-types -w 123456
```

### Templates

List iteration templates:

```bash
tapd iteration templates -w 123456
```

List fields for a specific iteration template:

```bash
tapd iteration template-fields -w 123456 --template-id 20001
```

List default fields for a workitem type:

```bash
tapd iteration template-fields -w 123456 --workitem-type-id 10001
```

`template-fields` requires exactly one of `--template-id` or
`--workitem-type-id`.

## History and Locks

### Changes

List iteration changes:

```bash
tapd iteration changes -w 123456 --iteration-id 1111112222001000001
```

Useful filters:

```bash
tapd iteration changes -w 123456 --iteration-id 1111112222001000001 --author alice
tapd iteration changes -w 123456 --iteration-id 1111112222001000001 --field status
```

### Lock and Unlock

Lock iteration scopes:

```bash
tapd iteration lock \
  -w 123456 \
  --iteration-id 1111112222001000001 \
  --lock-types story,bug
```

Unlock iteration scopes:

```bash
tapd iteration unlock \
  -w 123456 \
  --iteration-id 1111112222001000001 \
  --lock-types story,bug
```

## Developer Notes

The iteration command implementation currently lives in:

```text
internal/cmd/iteration.go
```

When adding or renaming iteration commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
