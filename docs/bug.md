# Bug Commands

Use `tapd bug` commands to create, inspect, list, update, and query TAPD bugs.

All bug commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most bug commands require a workspace:

```bash
tapd bug list --workspace-id 123456
```

The short form is also supported:

```bash
tapd bug list -w 123456
```

## Output Formats

Bug commands default to table output:

```bash
tapd bug list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd bug list -w 123456 --format json
```

## Common Workflows

### Create a Bug

```bash
tapd bug create -w 123456 --title "Login fails" --owner alice
```

Useful optional flags:

```bash
tapd bug create \
  -w 123456 \
  --title "Login fails" \
  --description "User cannot log in with a valid token" \
  --owner alice \
  --reporter bob \
  --priority-label High \
  --severity serious \
  --module auth \
  --iteration-id 12345 \
  --label "cli|auth"
```

### View a Bug

```bash
tapd bug view 1111112222001000001 -w 123456
```

This command queries the bug list API with the requested ID and returns the
first matching bug.

### List Bugs

```bash
tapd bug list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd bug list -w 123456 --owner alice
tapd bug list -w 123456 --creator bob
tapd bug list -w 123456 --ids 1111112222001000001,1111112222001000002
```

Request specific fields from TAPD:

```bash
tapd bug list -w 123456 --fields id,title,status,current_owner,reporter,modified
```

### Count Bugs

```bash
tapd bug count -w 123456
tapd bug count -w 123456 --owner alice
tapd bug count -w 123456 --creator bob
```

### Update a Bug

```bash
tapd bug update 1111112222001000001 -w 123456 --title "Updated bug title"
```

Update several common fields:

```bash
tapd bug update 1111112222001000001 \
  -w 123456 \
  --status resolved \
  --owner bob \
  --priority-label High \
  --severity serious \
  --resolution fixed \
  --deadline 2026-06-30
```

`tapd bug update` exposes common fields currently mapped by the CLI. Use
`tapd bug batch-update` when you need fields not exposed as flags.

### Batch Update Bugs

Batch updates are read from a JSON file:

```bash
tapd bug batch-update -w 123456 --file bugs.json
```

The file can contain an array of update objects:

```json
[
  {
    "id": 1111112222001000001,
    "title": "Updated bug title",
    "current_owner": "alice"
  },
  {
    "id": 1111112222001000002,
    "status": ["resolved"],
    "resolution": ["fixed"]
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
      "title": "Updated bug title"
    }
  ]
}
```

The top-level `project_id` is always set from `--workspace-id`.

## Metadata Commands

### Fields

List field definitions and option counts:

```bash
tapd bug fields -w 123456
```

Include closed field options:

```bash
tapd bug fields -w 123456 --all-options
```

List field labels:

```bash
tapd bug field-labels -w 123456
```

Use JSON output to inspect full field options:

```bash
tapd bug fields -w 123456 --format json
```

### Templates

List bug templates:

```bash
tapd bug templates -w 123456
```

List template fields:

```bash
tapd bug template-fields -w 123456 --template-id 20001
```

Use `priority_label` instead of `priority` in template fields:

```bash
tapd bug template-fields -w 123456 --template-id 20001 --use-priority-label
```

## History and Relations

### Changes

List bug changes:

```bash
tapd bug changes -w 123456 --bug-ids 1111112222001000001
```

Count bug changes:

```bash
tapd bug changes count -w 123456 --bug-ids 1111112222001000001
```

Useful filters:

```bash
tapd bug changes -w 123456 --author alice
tapd bug changes -w 123456 --field status
tapd bug changes -w 123456 --include-add-bug
```

### Related Stories

```bash
tapd bug related-stories -w 123456 --bug-ids 1111112222001000001
```

Multiple bug IDs are comma-separated:

```bash
tapd bug related-stories -w 123456 --bug-ids 1111112222001000001,1111112222001000002
```

### Link and Unlink Bugs

Create bug links:

```bash
tapd bug link \
  -w 123456 \
  --bug-id 1111112222001000001 \
  --relate-bug-ids 1111112222001000002,1111112222001000003
```

Remove bug links by relation IDs:

```bash
tapd bug unlink \
  -w 123456 \
  --bug-id 1111112222001000001 \
  --link-ids 90001,90002
```

## View and Query Helpers

### Bugs by View

```bash
tapd bug by-view -w 123456 --view-conf-id 30001
```

Personal views can require the current user:

```bash
tapd bug by-view -w 123456 --view-conf-id 30001 --current-user alice
```

### Convert IDs to Query Token

```bash
tapd bug convert-ids -w 123456 --ids 1111112222001000001,1111112222001000002
```

This returns the TAPD list `queryToken` and link for the selected bugs.

### Removed Bugs

```bash
tapd bug removed -w 123456
tapd bug removed -w 123456 --creator alice
tapd bug removed -w 123456 --include-all
```

## Developer Notes

The bug command implementation currently lives in:

```text
internal/cmd/bug.go
```

When adding or renaming bug commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
