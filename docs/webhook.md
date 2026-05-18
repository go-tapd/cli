# Webhook Commands

Use `tapd webhook` commands to inspect and validate TAPD webhook payloads, or
to run a small local receiver while developing integrations.

Webhook commands do not call the TAPD REST API. They use the SDK webhook parser
for local payload decoding, so they can run without `tapd login`.

## Output Formats

Webhook commands default to table output:

```bash
tapd webhook inspect --file payload.json
```

Use JSON when piping to another program or when you need the original decoded
payload:

```bash
tapd webhook inspect --file payload.json --format json
```

## Inspect Payloads

Inspect a payload from a file:

```bash
tapd webhook inspect --file payload.json
```

Read from stdin:

```bash
cat payload.json | tapd webhook inspect
```

The table output includes the event type, workspace ID, entity ID, current user,
event ID, created time, and a best-effort subject from common fields such as
`name`, `title`, `old_name`, and `old_title`.

## Validate Payloads

Validate a payload from a file:

```bash
tapd webhook validate --file payload.json
```

Read from stdin:

```bash
cat payload.json | tapd webhook validate
```

The command exits with a non-zero status when the payload is not valid JSON, has
no supported `event` value, or cannot be decoded by the SDK webhook package.

## Serve a Local Receiver

Run a local POST endpoint:

```bash
tapd webhook serve --addr 127.0.0.1:8080 --path /webhook
```

Send a sample payload:

```bash
curl -X POST \
  -H 'Content-Type: application/json' \
  --data @payload.json \
  http://127.0.0.1:8080/webhook
```

The receiver validates each request body with the SDK parser. Valid requests
return a JSON response like:

```json
{
  "event": "story::create",
  "id": "1111112222001071295",
  "ok": true,
  "workspace_id": "11112222"
}
```

The command logs one summary line for each accepted request. Use
`--max-body-bytes` to control the maximum request body size.

## Supported Events

Supported event types are determined by the SDK webhook package. The current SDK
supports story, task, bug, story comment, task comment, bug comment, and
iteration create/update/delete events.

## Implementation Notes

The webhook command implementation currently lives in:

```text
internal/cmd/webhook.go
```

When adding or renaming webhook commands:

- Update `features.md`.
- Update the command list and examples in `README.md`.
- Keep this document aligned with the supported flags and table output.
