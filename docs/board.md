# Board Commands

Use `tapd board` commands to create, list, and update TAPD board cards, and to
list board columns.

All board commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd board card list --workspace-id 123456
```

The short form is also supported:

```bash
tapd board card list -w 123456
```

## Output Formats

Board commands default to table output:

```bash
tapd board card list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd board card list -w 123456 --format json
```

## Board Cards

### Create a Card

Create a card in a board column:

```bash
tapd board card create \
  -w 123456 \
  --board-id 20001 \
  --column-id 30001 \
  --name "Investigate login issue"
```

Set optional owner, dates, label, and description:

```bash
tapd board card create \
  -w 123456 \
  --board-id 20001 \
  --column-id 30001 \
  --name "Investigate login issue" \
  --owner alice \
  --status open \
  --begin 2026-06-01 \
  --due 2026-06-07 \
  --label 10001 \
  --description "Reproduce and document the failure"
```

### List Cards

List cards:

```bash
tapd board card list -w 123456 --limit 20 --page 1
```

Filter by board, column, owner, or status:

```bash
tapd board card list -w 123456 --board-id 20001
tapd board card list -w 123456 --column-id 30001
tapd board card list -w 123456 --owner alice
tapd board card list -w 123456 --status open
```

Filter by IDs or label:

```bash
tapd board card list -w 123456 --ids 10001,10002
tapd board card list -w 123456 --label 90001
```

Request specific fields:

```bash
tapd board card list -w 123456 --fields id,name,b_board_id,b_column_id,status,owner
```

Date filters accept the expressions supported by TAPD:

```bash
tapd board card list -w 123456 --created "2026-06-01~2026-06-30"
tapd board card list -w 123456 --begin "2026-06-01~2026-06-30"
tapd board card list -w 123456 --due "2026-06-01~2026-06-30"
```

### Update a Card

Move a card to another column:

```bash
tapd board card update 10001 -w 123456 --column-id 30002
```

Update card details:

```bash
tapd board card update 10001 \
  -w 123456 \
  --name "Investigate login issue" \
  --owner bob \
  --status done \
  --due 2026-06-10
```

At least one card field flag is required.

## Board Columns

List columns:

```bash
tapd board columns -w 123456 --limit 20 --page 1
```

Filter by board, status, creator, or IDs:

```bash
tapd board columns -w 123456 --board-id 20001
tapd board columns -w 123456 --status open
tapd board columns -w 123456 --creator alice
tapd board columns -w 123456 --ids 30001,30002
```

Request specific fields:

```bash
tapd board columns -w 123456 --fields id,name,board_id,status,creator,created
```

## Developer Notes

The board command implementation currently lives in:

```text
internal/cmd/board.go
```

When adding or renaming board commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
