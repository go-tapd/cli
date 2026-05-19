# Story Commands

Use `tapd story` commands to create, inspect, list, update, and query TAPD
stories.

All story commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most story commands require a workspace:

```bash
tapd story list --workspace-id 123456
```

The short form is also supported:

```bash
tapd story list -w 123456
```

## Output Formats

Story commands default to table output:

```bash
tapd story list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd story list -w 123456 --format json
```

## Common Workflows

### Create a Story

```bash
tapd story create -w 123456 --name "Implement login validation" --owner alice
```

Useful optional flags:

```bash
tapd story create \
  -w 123456 \
  --name "Implement login validation" \
  --description "Validate TAPD credentials before saving config" \
  --owner alice \
  --creator alice \
  --priority-label High \
  --iteration-id 12345 \
  --label "cli|auth"
```

### View a Story

```bash
tapd story view 1111112222001000001 -w 123456
```

This command queries the story list API with the requested ID and returns the
first matching story.

### List Stories

```bash
tapd story list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd story list -w 123456 --owner alice
tapd story list -w 123456 --creator bob
tapd story list -w 123456 --status developing
tapd story list -w 123456 --ids 1111112222001000001,1111112222001000002
```

Request specific fields from TAPD:

```bash
tapd story list -w 123456 --fields id,name,status,owner,modified
```

### Count Stories

```bash
tapd story count -w 123456
tapd story count -w 123456 --owner alice
tapd story count -w 123456 --status developing
```

### Update a Story

```bash
tapd story update 1111112222001000001 -w 123456 --name "Updated title"
```

Update several common fields:

```bash
tapd story update 1111112222001000001 \
  -w 123456 \
  --current-user alice \
  --status developing \
  --owner bob \
  --priority-label High \
  --due 2026-06-30
```

`tapd story update` exposes only the common fields currently mapped by the CLI.
Use `tapd story batch-update` when you need fields not exposed as flags.

### Batch Update Stories

Batch updates are read from a JSON file:

```bash
tapd story batch-update -w 123456 --file stories.json
```

The file can contain an array of update objects:

```json
[
  {
    "id": 1111112222001000001,
    "name": "Updated story title",
    "owner": "alice"
  },
  {
    "id": 1111112222001000002,
    "status": "developing",
    "current_user": "bob"
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
      "name": "Updated story title"
    }
  ]
}
```

The top-level `workspace_id` is always set from `--workspace-id`.

## Metadata Commands

### Categories

List story categories:

```bash
tapd story categories -w 123456
```

Count story categories:

```bash
tapd story categories count -w 123456
```

Filter categories:

```bash
tapd story categories -w 123456 --name "Backend"
tapd story categories -w 123456 --parent-id 1001
```

### Fields

List field definitions and option counts:

```bash
tapd story fields -w 123456
```

List field labels:

```bash
tapd story field-labels -w 123456
```

Use JSON output to inspect full field options:

```bash
tapd story fields -w 123456 --format json
```

### Templates

List story templates:

```bash
tapd story templates -w 123456
```

Filter by workitem type:

```bash
tapd story templates -w 123456 --workitem-type-id 10001
```

List template fields:

```bash
tapd story template-fields -w 123456 --template-id 20001
```

## History and Relations

### Changes

List story changes:

```bash
tapd story changes -w 123456 --story-ids 1111112222001000001
```

Count story changes:

```bash
tapd story changes count -w 123456 --story-ids 1111112222001000001
```

Useful filters:

```bash
tapd story changes -w 123456 --creator alice
tapd story changes -w 123456 --change-type manual_update
tapd story changes -w 123456 --change-field status
```

### Related Bugs

```bash
tapd story related-bugs -w 123456 --story-ids 1111112222001000001
```

Multiple story IDs are comma-separated:

```bash
tapd story related-bugs -w 123456 --story-ids 1111112222001000001,1111112222001000002
```

### Related Test Cases

```bash
tapd story related-test-cases -w 123456 --story-id 1111112222001000001
```

Disable test plan relation data:

```bash
tapd story related-test-cases -w 123456 --story-id 1111112222001000001 --include-test-plan=false
```

## View and Query Helpers

### Stories by View

```bash
tapd story by-view -w 123456 --view-conf-id 30001
```

Personal views can require the current user:

```bash
tapd story by-view -w 123456 --view-conf-id 30001 --current-user alice
```

### Convert IDs to Query Token

```bash
tapd story convert-ids -w 123456 --ids 1111112222001000001,1111112222001000002
```

This returns the TAPD list `queryToken` and link for the selected stories.

### Removed Stories

```bash
tapd story removed -w 123456
tapd story removed -w 123456 --creator alice
tapd story removed -w 123456 --archived
```

## Developer Notes

The story command implementation lives in:

```text
internal/cmd/story.go
```

When adding or renaming story commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
