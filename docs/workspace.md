# Workspace Commands

Use `tapd workspace` commands to inspect and maintain TAPD workspaces, members,
documents, work item ID mappings, activity logs, and work calendars.

All workspace commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd workspace view --workspace-id 123456
```

The short form is also supported:

```bash
tapd workspace users -w 123456
```

## Output Formats

Workspace commands default to table output:

```bash
tapd workspace roles -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd workspace custom-fields -w 123456 --format json
```

## Basic Workspace Data

Show workspace details:

```bash
tapd workspace view -w 123456
```

List users:

```bash
tapd workspace users -w 123456
tapd workspace users -w 123456 --user alice
tapd workspace users -w 123456 --fields user,user_id,role_id,name,email,real_join_time
```

List roles:

```bash
tapd workspace roles -w 123456
```

## Workspace Lists

List sub-workspaces:

```bash
tapd workspace sub-workspaces -w 123456
tapd workspace sub-workspaces -w 123456 --template-id 20001
```

List company workspaces:

```bash
tapd workspace company-workspaces --company-id 90001
tapd workspace company-workspaces --company-id 90001 --category project --with-extends
```

List workspaces a user participates in:

```bash
tapd workspace participant-workspaces --company-id 90001 --nick alice
```

## Members

Add a workspace member:

```bash
tapd workspace add-member -w 123456 --nick alice --role-ids 10001,10002
```

Cloud deployments can require the member company ID:

```bash
tapd workspace add-member -w 123456 --nick alice --company-id 90001 --role-ids 10001
```

## Update Workspace Info

Update one workspace field:

```bash
tapd workspace update -w 123456 --field name --value "New workspace name"
```

Use field names accepted by TAPD's workspace update API.

## Metadata and Documents

List workspace custom field settings:

```bash
tapd workspace custom-fields -w 123456
tapd workspace custom-fields -w 123456 --format json
```

List workspace documents:

```bash
tapd workspace documents -w 123456 --limit 20 --page 1
tapd workspace documents -w 123456 --fields id,name,type,folder_id,creator,modified
```

## Short ID Conversion

Convert work item short IDs to long IDs:

```bash
tapd workspace short-id convert \
  -w 123456 \
  --entity-type story \
  --short-ids "1001;1002"
```

Convert long IDs to short IDs:

```bash
tapd workspace short-id convert \
  -w 123456 \
  --entity-type bug \
  --long-ids "1111112222001000001;1111112222001000002"
```

## Member Activity Logs

List member activity logs:

```bash
tapd workspace member-activity-log -w 123456 --limit 20 --page 1
```

Filter logs:

```bash
tapd workspace member-activity-log -w 123456 --operator alice
tapd workspace member-activity-log -w 123456 --operate-type add --operate-object story
tapd workspace member-activity-log -w 123456 --start-time "2026-06-01 00:00" --end-time "2026-06-30 23:59"
```

For company-level logs, use the company ID as `--workspace-id` and pass
`--company-only`.

## Work Calendars

Set a custom work calendar:

```bash
tapd workspace calendar set-custom \
  -w 123456 \
  --year 2026 \
  --weekdays 1,2,3,4,5 \
  --holidays 2026-10-01,2026-10-02 \
  --workdays 2026-10-10
```

Enable a calendar type:

```bash
tapd workspace calendar enable -w 123456 --type custom
tapd workspace calendar enable -w 123456 --type system
```

View a custom calendar:

```bash
tapd workspace calendar view-custom -w 123456 --year 2026
```

List calendar settings:

```bash
tapd workspace calendar settings -w 123456
```

## Developer Notes

The workspace command implementation currently lives in:

```text
internal/cmd/workspace.go
```

When adding or renaming workspace commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
