# Task Commands

Use `tapd task` commands to create, inspect, list, update, and query TAPD
tasks.

All task commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most task commands require a workspace:

```bash
tapd task list --workspace-id 123456
```

The short form is also supported:

```bash
tapd task list -w 123456
```

## Output Formats

Task commands default to table output:

```bash
tapd task list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd task list -w 123456 --format json
```

## Common Workflows

### Create a Task

```bash
tapd task create -w 123456 --name "Implement login" --owner alice
```

Useful optional flags:

```bash
tapd task create \
  -w 123456 \
  --name "Implement login" \
  --description "Add TAPD credential validation" \
  --owner alice \
  --creator bob \
  --status open \
  --story-ids 1111112222001000001 \
  --iteration-id 2222223333001000001 \
  --priority-label High \
  --begin 2026-06-01 \
  --due 2026-06-14
```

Set task custom fields at creation time with repeatable `--field` flags:

```bash
tapd task create \
  -w 123456 \
  --name "Implement login" \
  --story-ids 1111112222001000001 \
  --field custom_field_one=开发阶段
```

Only task custom fields supported by the typed SDK request are accepted, such as
`custom_field_one` through `custom_field_50`.

### View a Task

```bash
tapd task view 1111112222001000001 -w 123456
```

This command queries the task list API with the requested ID and returns the
first matching task.

### List Tasks

```bash
tapd task list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd task list -w 123456 --owner alice
tapd task list -w 123456 --creator bob
tapd task list -w 123456 --status progressing
tapd task list -w 123456 --story-ids 1111112222001000001
tapd task list -w 123456 --iteration-id 2222223333001000001
tapd task list -w 123456 --ids 1111112222001000001,1111112222001000002
```

Request specific fields from TAPD:

```bash
tapd task list -w 123456 --fields id,name,status,owner,creator,progress
```

### Count Tasks

```bash
tapd task count -w 123456
tapd task count -w 123456 --owner alice
tapd task count -w 123456 --status done
```

### Update a Task

```bash
tapd task update 1111112222001000001 -w 123456 --name "Updated task title"
```

Update several common fields:

```bash
tapd task update 1111112222001000001 \
  -w 123456 \
  --current-user alice \
  --status progressing \
  --owner bob \
  --progress 50 \
  --due 2026-06-14
```

`tapd task update` also supports repeatable `--field key=value` flags for task
custom fields.

### Batch Update Tasks

Batch updates are read from a JSON file:

```bash
tapd task batch-update -w 123456 --file tasks.json
```

The file can contain an array of update objects:

```json
[
  {
    "id": 1111112222001000001,
    "name": "Updated task title",
    "owner": "alice"
  },
  {
    "id": 1111112222001000002,
    "status": "done",
    "progress": 100
  }
]
```

The CLI adds `workspace_id` to each item when it is missing.

The file can also contain the SDK request object:

```json
{
  "workitems": [
    {
      "id": 1111112222001000001,
      "name": "Updated task title"
    }
  ]
}
```

The top-level `workspace_id` is always set from `--workspace-id`.

## Metadata Commands

### Fields

List task fields:

```bash
tapd task fields -w 123456
```

Use JSON output to inspect full field options:

```bash
tapd task fields -w 123456 --format json
```

## History

### Changes

List task changes:

```bash
tapd task changes -w 123456 --task-id 1111112222001000001
```

Count task changes:

```bash
tapd task changes count -w 123456 --task-id 1111112222001000001
```

Useful filters:

```bash
tapd task changes -w 123456 --creator alice
tapd task changes -w 123456 --change-summary "status"
tapd task changes -w 123456 --need-parse-changes=false
```

### Removed Tasks

```bash
tapd task removed -w 123456
tapd task removed -w 123456 --creator alice
tapd task removed -w 123456 --archived
```

## Unsupported Commands

`tapd task by-view` remains unavailable because the current SDK does not expose
a typed task-by-view API.

## Developer Notes

The task command implementation currently lives in:

```text
internal/cmd/task.go
```

When adding or renaming task commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
