# Release Commands

Use `tapd release` commands to manage TAPD release plans, and `tapd launch-form`
commands to manage launch forms.

All release commands use the configured TAPD credentials from `tapd login` or
the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd release list --workspace-id 123456
```

The short form is also supported:

```bash
tapd release list -w 123456
```

## Output Formats

Release commands default to table output:

```bash
tapd release list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd launch-form list -w 123456 --format json
```

## Release Plans

### Create a Release

```bash
tapd release create \
  -w 123456 \
  --name "v1.2.0" \
  --start-date 2026-06-01 \
  --end-date 2026-06-30
```

Optional fields:

```bash
tapd release create \
  -w 123456 \
  --name "v1.2.0" \
  --start-date 2026-06-01 \
  --end-date 2026-06-30 \
  --creator alice \
  --description "June release"
```

### View a Release

```bash
tapd release view 10001 -w 123456
```

This command queries the release list API with the requested ID and returns the
first matching release.

### List Releases

```bash
tapd release list -w 123456 --limit 20 --page 1
```

Filter releases:

```bash
tapd release list -w 123456 --name "v1.2"
tapd release list -w 123456 --creator alice
tapd release list -w 123456 --status open
tapd release list -w 123456 --ids 10001,10002
```

Request specific fields:

```bash
tapd release list -w 123456 --fields id,name,status,startdate,enddate,creator
```

Date filters accept the expressions supported by TAPD:

```bash
tapd release list -w 123456 --created "2026-06-01~2026-06-30"
tapd release list -w 123456 --modified "2026-06-01~2026-06-30"
```

### Count Releases

```bash
tapd release count -w 123456
tapd release count -w 123456 --status open
tapd release count -w 123456 --creator alice
```

### Update a Release

```bash
tapd release update 10001 -w 123456 --status closed
tapd release update 10001 -w 123456 --name "v1.2.1" --end-date 2026-07-05
```

At least one release field flag is required.

## Launch Forms

### Create a Launch Form

```bash
tapd launch-form create \
  -w 123456 \
  --creator alice \
  --template-id 20001 \
  --title "v1.2.0 review"
```

Optional metadata:

```bash
tapd launch-form create \
  -w 123456 \
  --creator alice \
  --template-id 20001 \
  --title "v1.2.0 review" \
  --version-type stable \
  --release-type regular \
  --signed-by bob \
  --archived-by carol \
  --cc dave
```

### List Launch Forms

```bash
tapd launch-form list -w 123456 --limit 20 --page 1
```

Filter forms:

```bash
tapd launch-form list -w 123456 --id 30001
tapd launch-form list -w 123456 --creator alice
tapd launch-form list -w 123456 --status open
tapd launch-form list -w 123456 --release-type regular
tapd launch-form list -w 123456 --signed-by bob
```

Request specific fields:

```bash
tapd launch-form list -w 123456 --fields id,title,status,creator,release_id,created
```

### Count Launch Forms

```bash
tapd launch-form count -w 123456
tapd launch-form count -w 123456 --creator alice
tapd launch-form count -w 123456 --status open
```

### Field Settings

List launch form custom field settings:

```bash
tapd launch-form fields -w 123456
```

Use JSON output to inspect full option payloads:

```bash
tapd launch-form fields -w 123456 --format json
```

### Templates

List launch form templates:

```bash
tapd launch-form templates -w 123456
```

### Activity Logs

List launch form activity logs:

```bash
tapd launch-form logs -w 123456 --form-id 30001
```

Use JSON output to inspect raw old and new values:

```bash
tapd launch-form logs -w 123456 --form-id 30001 --format json
```

## Developer Notes

The release and launch form command implementation currently lives in:

```text
internal/cmd/release.go
```

When adding or renaming release commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
