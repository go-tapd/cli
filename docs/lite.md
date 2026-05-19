# Lite Commands

Use `tapd lite` commands for TAPD Lite resources that are supported by the
current SDK.

The current SDK exposes Lite comments through the generic comment API with
`entry_type=mini_items`. The CLI wraps that as `tapd lite comment ...` so you do
not need to pass the entry type manually.

## Output Formats

Lite commands default to table output:

```bash
tapd lite comment list -w 123456 --entry-id 1111112222001000001
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd lite comment list -w 123456 --entry-id 1111112222001000001 --format json
```

## Lite Comments

### Create a Comment

```bash
tapd lite comment create \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --author alice \
  --description "Looks good"
```

Optional reply fields:

```bash
tapd lite comment create \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --author alice \
  --description "Reply content" \
  --root-id 10001 \
  --reply-id 10002
```

### List Comments

```bash
tapd lite comment list -w 123456 --entry-id 1111112222001000001
```

Filter by common fields:

```bash
tapd lite comment list -w 123456 --ids 10001,10002
tapd lite comment list -w 123456 --entry-id 1111112222001000001 --author alice
tapd lite comment list -w 123456 --created "2026-01-01~2026-01-31"
tapd lite comment list -w 123456 --root-id 10001
```

### Count Comments

```bash
tapd lite comment count -w 123456 --entry-id 1111112222001000001
tapd lite comment count -w 123456 --author alice
```

### Update a Comment

```bash
tapd lite comment update 10001 -w 123456 --description "Updated content"
```

Set the change creator when updating:

```bash
tapd lite comment update 10001 \
  -w 123456 \
  --description "Updated content" \
  --change-creator alice
```

## SDK Coverage Notes

Other Lite sections in `features.md`, such as workitems, spaces, attachments,
and Lite Git integrations, are still waiting on typed SDK support. They should
be implemented under `tapd lite ...` once the SDK exposes stable request and
response types for those APIs.

## Implementation Notes

The Lite command implementation currently lives in:

```text
internal/cmd/lite.go
```

When adding or renaming Lite commands:

- Update `features.md`.
- Update the command list and examples in `README.md`.
- Keep this document aligned with the supported flags and table output.
