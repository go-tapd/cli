# Setting Commands

Use `tapd setting` commands to query TAPD workspace settings.

All setting commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Setting commands require a workspace:

```bash
tapd setting workspace --workspace-id 123456
```

The short form is also supported:

```bash
tapd setting workspace -w 123456
```

## Output Formats

Setting commands default to table output:

```bash
tapd setting workspace -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd setting workspace -w 123456 --format json
```

## Workspace Settings

Show workspace settings:

```bash
tapd setting workspace -w 123456
```

Query one setting by name:

```bash
tapd setting workspace -w 123456 --type is_enabled_story_category
tapd setting workspace -w 123456 --type workspace_metrology
```

Known setting names exposed by the SDK are:

- `is_enabled_story_category`
- `workspace_metrology`

## Developer Notes

The setting command implementation currently lives in:

```text
internal/cmd/setting.go
```

When adding or renaming setting commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
