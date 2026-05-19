# Timesheet Commands

Use `tapd timesheet` commands to create, list, count, update, and delete TAPD
timesheet records.

All timesheet commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd timesheet list --workspace-id 123456
```

The short form is also supported:

```bash
tapd timesheet list -w 123456
```

## Output Formats

Timesheet commands default to table output:

```bash
tapd timesheet list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd timesheet list -w 123456 --format json
```

## Create Timesheets

Create a timesheet for a story, task, or bug:

```bash
tapd timesheet create \
  -w 123456 \
  --entity-type story \
  --entity-id 1111112222001000001 \
  --timespent 2h \
  --owner alice
```

Add remaining time, spent date, and memo:

```bash
tapd timesheet create \
  -w 123456 \
  --entity-type task \
  --entity-id 1111112222001000002 \
  --timespent 1.5h \
  --timeremain 3h \
  --spentdate 2026-06-01 \
  --owner alice \
  --memo "API implementation"
```

Supported entity types are `story`, `task`, and `bug`.

## List Timesheets

List recent timesheets:

```bash
tapd timesheet list -w 123456 --limit 20 --page 1
```

Filter by owner or entity:

```bash
tapd timesheet list -w 123456 --owner alice
tapd timesheet list -w 123456 --entity-type bug --entity-id 1111112222001000003
```

Filter by IDs:

```bash
tapd timesheet list -w 123456 --ids 10001,10002
```

Request specific fields:

```bash
tapd timesheet list -w 123456 --fields id,entity_type,entity_id,timespent,owner
```

Include deleted records or exclude parent story timesheets:

```bash
tapd timesheet list -w 123456 --include-deleted
tapd timesheet list -w 123456 --exclude-parent-story
```

Date filters accept the expressions supported by TAPD:

```bash
tapd timesheet list -w 123456 --spentdate "2026-06-01~2026-06-30"
tapd timesheet list -w 123456 --created "2026-06-01~2026-06-30"
tapd timesheet list -w 123456 --modified "2026-06-01~2026-06-30"
```

## Count Timesheets

Count all matching timesheets:

```bash
tapd timesheet count -w 123456
```

Apply the same filters used by `list`:

```bash
tapd timesheet count -w 123456 --owner alice
tapd timesheet count -w 123456 --entity-type story --entity-id 1111112222001000001
tapd timesheet count -w 123456 --include-deleted
```

## Update Timesheets

Update spent time, remaining time, or memo:

```bash
tapd timesheet update 10001 -w 123456 --timespent 3h
tapd timesheet update 10001 -w 123456 --timeremain 2h --memo "Adjusted estimate"
```

At least one of `--timespent`, `--timeremain`, or `--memo` is required.

## Delete Timesheets

Delete one or more timesheet cost IDs from an entity:

```bash
tapd timesheet delete \
  -w 123456 \
  --entity-type story \
  --entity-id 1111112222001000001 \
  --cost-ids 10001,10002
```

The table output shows the TAPD message, the count of successfully deleted cost
IDs, and the count of failed result groups. Use JSON output to inspect each
failed item in detail:

```bash
tapd timesheet delete \
  -w 123456 \
  --entity-type story \
  --entity-id 1111112222001000001 \
  --cost-ids 10001,10002 \
  --format json
```

## Developer Notes

The timesheet command implementation currently lives in:

```text
internal/cmd/timesheet.go
```

When adding or renaming timesheet commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
