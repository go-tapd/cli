# Workflow Commands

Use `tapd workflow` commands to query TAPD workflow metadata.

All workflow commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Workflow commands require a workspace:

```bash
tapd workflow last-steps --workspace-id 123456
```

The short form is also supported:

```bash
tapd workflow last-steps -w 123456
```

## Output Formats

Workflow commands default to table output:

```bash
tapd workflow last-steps -w 123456
```

Use JSON when you need the full status alias/name map:

```bash
tapd workflow last-steps -w 123456 --format json
```

## Last Steps

List workflow final statuses:

```bash
tapd workflow last-steps -w 123456
```

The TAPD API currently supports the `story` workflow system:

```bash
tapd workflow last-steps -w 123456 --system story
```

Group by workflow ID or workitem type ID:

```bash
tapd workflow last-steps -w 123456 --group-key workflow_id
tapd workflow last-steps -w 123456 --group-key workitem_type_id
```

## Developer Notes

The workflow command implementation currently lives in:

```text
internal/cmd/workflow.go
```

When adding or renaming workflow commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
