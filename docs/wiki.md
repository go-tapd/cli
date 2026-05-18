# Wiki Commands

Use `tapd wiki` commands to create, list, count, update, and inspect TAPD Wiki
data.

All wiki commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd wiki list --workspace-id 123456
```

The short form is also supported:

```bash
tapd wiki list -w 123456
```

## Output Formats

Wiki commands default to table output:

```bash
tapd wiki list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd wiki list -w 123456 --format json
```

## Wiki Pages

Create a wiki page:

```bash
tapd wiki create -w 123456 --name "Release notes" --creator alice
```

Create a child page with Markdown content:

```bash
tapd wiki create \
  -w 123456 \
  --name "Release notes" \
  --creator alice \
  --parent-wiki-id 10001 \
  --markdown-description "## v1.2.0"
```

List and count pages:

```bash
tapd wiki list -w 123456 --limit 20 --page 1
tapd wiki list -w 123456 --creator alice
tapd wiki count -w 123456 --name "Release"
```

Update a wiki page:

```bash
tapd wiki update 10001 -w 123456 --name "Updated release notes"
tapd wiki update 10001 -w 123456 --markdown-description "## Updated"
```

At least one wiki field flag is required for `update`.

## Drawio

Get drawio data:

```bash
tapd wiki drawio -w 123456 --id 20001
```

If TAPD requires a token:

```bash
tapd wiki drawio -w 123456 --id 20001 --token abc123
```

Table output shows the data ID and XML byte length. Use JSON output for the raw
XML payload:

```bash
tapd wiki drawio -w 123456 --id 20001 --format json
```

## Followers

List followers:

```bash
tapd wiki followers -w 123456 --wiki-id 10001
tapd wiki followers -w 123456 --user alice
```

Count followers:

```bash
tapd wiki followers count -w 123456 --wiki-id 10001
```

## Permissions

List wiki entity permissions:

```bash
tapd wiki permissions -w 123456 --wiki-id 10001
```

Filter by target:

```bash
tapd wiki permissions -w 123456 --wiki-id 10001 --target-type user --target-id alice
```

## Tags

List tags:

```bash
tapd wiki tags -w 123456
tapd wiki tags -w 123456 --wiki-id 10001
tapd wiki tags -w 123456 --tag release
```

Count tags:

```bash
tapd wiki tags count -w 123456 --wiki-id 10001
```

## Attachments

Count wiki attachments:

```bash
tapd wiki attachments count -w 123456 --wiki-id 10001
tapd wiki attachments count -w 123456 --filename spec.pdf
```

## Developer Notes

The wiki command implementation currently lives in:

```text
internal/cmd/wiki.go
```

When adding or renaming wiki commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
