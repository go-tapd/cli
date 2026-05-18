# Comment Commands

Use `tapd comment` commands to create, list, count, and update TAPD comments.
These commands target the classic TAPD comment API, not the Lite comment API.

All comment commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd comment list --workspace-id 123456
```

The short form is also supported:

```bash
tapd comment list -w 123456
```

## Output Formats

Comment commands default to table output:

```bash
tapd comment list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd comment list -w 123456 --format json
```

## Entry Types

The SDK supports these entry types:

- `bug`
- `bug_remark`
- `stories`
- `tasks`
- `wiki`
- `mini_items`

## Create Comments

Create a comment on a story:

```bash
tapd comment create \
  -w 123456 \
  --entry-type stories \
  --entry-id 1111112222001000001 \
  --author alice \
  --description "Looks good"
```

Create a reply:

```bash
tapd comment create \
  -w 123456 \
  --entry-type bug \
  --entry-id 1111112222001000002 \
  --author alice \
  --description "I will verify this" \
  --root-id 10001 \
  --reply-id 10002
```

`--title` is optional and maps to TAPD's comment title field.

## List Comments

List recent comments:

```bash
tapd comment list -w 123456 --limit 20 --page 1
```

Filter by entry:

```bash
tapd comment list -w 123456 --entry-type bug --entry-id 1111112222001000002
tapd comment list -w 123456 --entry-type stories --entry-id 1111112222001000001
```

Filter by IDs, author, or reply chain:

```bash
tapd comment list -w 123456 --ids 10001,10002
tapd comment list -w 123456 --author alice
tapd comment list -w 123456 --root-id 10001
tapd comment list -w 123456 --reply-id 10002
```

Request specific fields:

```bash
tapd comment list -w 123456 --fields id,entry_type,entry_id,author,created
```

Date filters accept the expressions supported by TAPD:

```bash
tapd comment list -w 123456 --created "2026-06-01~2026-06-30"
tapd comment list -w 123456 --modified "2026-06-01~2026-06-30"
```

## Count Comments

Count all matching comments:

```bash
tapd comment count -w 123456
```

Apply the same filters used by `list`:

```bash
tapd comment count -w 123456 --entry-type bug --entry-id 1111112222001000002
tapd comment count -w 123456 --author alice
```

## Update Comments

Update a comment body:

```bash
tapd comment update 10001 -w 123456 --description "Updated comment"
```

Set the change creator when TAPD requires it:

```bash
tapd comment update 10001 \
  -w 123456 \
  --description "Updated comment" \
  --change-creator alice
```

## Developer Notes

The comment command implementation currently lives in:

```text
internal/cmd/comment.go
```

When adding or renaming comment commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
