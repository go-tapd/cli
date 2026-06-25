# Raw API Command

Use `tapd api` for advanced TAPD API calls when a typed `tapd <resource>`
command does not cover an endpoint or field yet.

Prefer typed commands such as `tapd story list`, `tapd task create`, and
`tapd bug update` for regular workflows. `tapd api` is a raw API escape hatch
that returns the TAPD response `data` value as JSON.

## Endpoint

Pass a TAPD API endpoint path relative to the configured TAPD API base URL:

```bash
tapd api tasks/get_fields_info --field workspace_id=123456
```

Absolute URLs are not allowed, so configured TAPD credentials are only sent to
the TAPD API base URL.

## Request Parameters

`tapd api` uses `GET` by default. Parameters are sent as query string values for
`GET` and `DELETE`, and as JSON body values for `POST`, `PUT`, and `PATCH`.
Use `--query` when any method needs explicit query string parameters.

Typed fields:

```bash
tapd api tasks/get_fields_info \
  --method GET \
  --field workspace_id=123456
```

Raw string fields:

```bash
tapd api tasks \
  --method POST \
  --raw-field name="Implement login" \
  --field workspace_id=123456
```

`--field` converts `true`, `false`, `null`, integers, and floats to JSON values.
If a `--field` value starts with `@`, it is treated as a file path and the field
value is replaced with that file's contents as a JSON string. Use `--raw-field`
to send a literal string that starts with `@`.
`--raw-field` always sends a string value.

Explicit query parameters:

```bash
tapd api tasks \
  --method POST \
  --query workspace_id=123456 \
  --field name="Implement login"
```

Repeat `--query` to send multi-valued query parameters with the same name.

## JSON Input

Use `--input` to send a JSON body from a file:

```bash
tapd api tasks --method POST --input task.json
```

Use `-` to read the JSON body from standard input:

```bash
printf '{"workspace_id":123456,"name":"Implement login"}' |
  tapd api tasks --method POST --input -
```

`--input` cannot be mixed with `--field` or `--raw-field`.

## Headers

Add request headers with `--header`:

```bash
tapd api tasks/get_fields_info \
  --field workspace_id=123456 \
  --header "X-Debug: true"
```

Overriding `Authorization` is not allowed. Use `tapd login`, `TAPD_ACCESS_TOKEN`,
or basic auth environment variables for authentication.

## Response Options

Include HTTP status and response headers:

```bash
tapd api tasks/get_fields_info --field workspace_id=123456 --include
```

Suppress the response body:

```bash
tapd api tasks/get_fields_info --field workspace_id=123456 --silent
```
