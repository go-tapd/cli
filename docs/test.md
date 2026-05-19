# Test Commands

Use `tapd test-case` and `tapd test-plan` commands to manage TAPD test cases,
test plans, execution progress, and result relations.

All commands use the configured TAPD credentials from `tapd login` or the
`TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET` environment
variables.

Most commands require a workspace:

```bash
tapd test-case list --workspace-id 123456
```

The short form is also supported:

```bash
tapd test-case list -w 123456
```

## Output Formats

Commands default to table output:

```bash
tapd test-plan list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd test-plan result -w 123456 --id 2222223333001000001 --format json
```

## Test Cases

### Create a Test Case

```bash
tapd test-case create -w 123456 --name "Login regression"
```

Useful optional flags:

```bash
tapd test-case create \
  -w 123456 \
  --name "Login regression" \
  --steps "Open login page, submit valid credentials" \
  --precondition "User exists" \
  --expectation "Dashboard opens" \
  --category-id 10001 \
  --status normal \
  --type functional \
  --priority high \
  --creator alice
```

Supported status values are `normal`, `updating`, and `abandon`.

### List Test Cases

```bash
tapd test-case list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd test-case list -w 123456 --ids 1111112222001000001,1111112222001000002
tapd test-case list -w 123456 --name login
tapd test-case list -w 123456 --status normal
tapd test-case list -w 123456 --category-id 10001
tapd test-case list -w 123456 --creator alice
```

Request specific fields from TAPD:

```bash
tapd test-case list -w 123456 --fields id,name,status,priority,modified
```

### Count Test Cases

```bash
tapd test-case count -w 123456
tapd test-case count -w 123456 --status normal
tapd test-case count -w 123456 --test-plan-id 2222223333001000001
```

### Update a Test Case

```bash
tapd test-case update 1111112222001000001 -w 123456 --status updating
```

Update several common fields:

```bash
tapd test-case update 1111112222001000001 \
  -w 123456 \
  --name "Login regression - updated" \
  --steps "Open login page, submit valid credentials, check dashboard" \
  --expectation "Dashboard opens without error" \
  --priority high
```

### Test Case Categories

```bash
tapd test-case categories -w 123456 --limit 20 --page 1
```

Filter by category metadata:

```bash
tapd test-case categories -w 123456 --name regression
tapd test-case categories -w 123456 --parent-id 10001
tapd test-case categories -w 123456 --creator alice
```

### Test Case Fields

```bash
tapd test-case fields -w 123456
```

Use JSON output to inspect full field options:

```bash
tapd test-case fields -w 123456 --format json
```

### Test Case Results

```bash
tapd test-case results \
  -w 123456 \
  --test-plan-id 2222223333001000001 \
  --test-case-id 1111112222001000001
```

The table output includes the execution status, executor, execution time, and
related bug IDs. Use JSON output when you need the full result payload.

## Test Plans

### Create a Test Plan

```bash
tapd test-plan create -w 123456 --name "Release regression"
```

Useful optional flags:

```bash
tapd test-plan create \
  -w 123456 \
  --name "Release regression" \
  --description "Regression scope for v1.2.0" \
  --owner alice \
  --creator bob \
  --start-date 2026-06-01 \
  --end-date 2026-06-14 \
  --iteration-id 2222223333001000001 \
  --version v1.2.0 \
  --status open
```

### List Test Plans

```bash
tapd test-plan list -w 123456 --limit 20 --page 1
```

Filter by common fields:

```bash
tapd test-plan list -w 123456 --ids 2222223333001000001,2222223333001000002
tapd test-plan list -w 123456 --name regression
tapd test-plan list -w 123456 --owner alice
tapd test-plan list -w 123456 --iteration-id 3333334444001000001
tapd test-plan list -w 123456 --status open
```

### Count Test Plans

```bash
tapd test-plan count -w 123456
tapd test-plan count -w 123456 --status open
tapd test-plan count -w 123456 --owner alice
```

### Update a Test Plan

```bash
tapd test-plan update 2222223333001000001 -w 123456 --status close
```

Update several common fields:

```bash
tapd test-plan update 2222223333001000001 \
  -w 123456 \
  --name "Release regression - final" \
  --modifier alice \
  --owner bob \
  --end-date 2026-06-20
```

### Test Plan Progress

```bash
tapd test-plan progress -w 123456 --id 2222223333001000001
```

The table output shows total story count, test case count, execution rate, and
result counters for pass, no pass, blocked, and unexecuted cases.

### Test Plan Result

```bash
tapd test-plan result -w 123456 --id 2222223333001000001
```

Filter by executor or include repeated execution data:

```bash
tapd test-plan result -w 123456 --id 2222223333001000001 --last-executor alice
tapd test-plan result -w 123456 --id 2222223333001000001 --include-repeat
```

### Related Bugs

```bash
tapd test-plan related-bugs -w 123456 --id 2222223333001000001
```

Use JSON output when you need each test result and linked bug detail:

```bash
tapd test-plan related-bugs -w 123456 --id 2222223333001000001 --format json
```

### Related Stories

```bash
tapd test-plan related-stories -w 123456 --test-plan-id 2222223333001000001
```

This returns the story IDs related to the selected test plan.

## Implementation Notes

The test command implementation currently lives in:

```text
internal/cmd/test.go
```

When adding or renaming test commands:

- Update `features.md`.
- Update the command list and examples in `README.md`.
- Keep this document aligned with the supported flags and table output.
