# Attachment Commands

Use `tapd attachment` commands to upload TAPD attachments, list existing
attachments, and fetch download URLs for attachments, images, and documents.

All attachment commands use the configured TAPD credentials from `tapd login`
or the `TAPD_ACCESS_TOKEN` / `TAPD_CLIENT_ID` / `TAPD_CLIENT_SECRET`
environment variables.

Attachment commands require a workspace:

```bash
tapd attachment list --workspace-id 123456
```

The short form is also supported:

```bash
tapd attachment list -w 123456
```

## Output Formats

Attachment commands default to table output:

```bash
tapd attachment list -w 123456
```

Use JSON when piping to another program or when you need the full SDK response
shape:

```bash
tapd attachment list -w 123456 --format json
```

## List Attachments

```bash
tapd attachment list -w 123456
```

Paginate large attachment lists:

```bash
tapd attachment list -w 123456 --limit 100 --page 2
```

Filter by common fields:

```bash
tapd attachment list -w 123456 --id 10001
tapd attachment list -w 123456 --entry-id 1111112222001000001
tapd attachment list -w 123456 --type story
tapd attachment list -w 123456 --filename spec
tapd attachment list -w 123456 --owner alice
```

## Upload Attachments

Upload a file to a custom field on a TAPD work item:

```bash
tapd attachment upload \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --custom-field custom_field_1 \
  --file ./spec.pdf
```

The upload type defaults to `story_custom_field`, which is the attachment type
documented by TAPD for custom-field uploads. Override it only when TAPD documents
another value for the target endpoint:

```bash
tapd attachment upload \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --type story_custom_field \
  --custom-field custom_field_1 \
  --owner alice \
  --filename spec-v2.pdf \
  --file ./spec.pdf
```

Upload an image by reading a local file and encoding it as base64:

```bash
tapd attachment upload-image-base64 \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --custom-field custom_field_1 \
  --image-file ./screenshot.png
```

Or pass base64 data directly:

```bash
tapd attachment upload-image-base64 \
  -w 123456 \
  --entry-id 1111112222001000001 \
  --custom-field custom_field_1 \
  --base64-data "$IMAGE_BASE64"
```

## Download URLs

Get a normal attachment download URL:

```bash
tapd attachment download-url -w 123456 --id 10001
```

Get an image download URL from an image path or URL:

```bash
tapd attachment image-url -w 123456 --image-path "/path/to/image.png"
```

Get a document download URL:

```bash
tapd attachment document-url -w 123456 --id 20001
```

## Developer Notes

The attachment command implementation currently lives in:

```text
internal/cmd/attachment.go
```

When adding or renaming attachment commands:

1. Reuse the typed SDK methods from `github.com/go-tapd/tapd`.
2. Keep table output compact and use `--format json` for full response data.
3. Update `features.md`.
4. Update this document.
5. Regenerate shell completion files if they have been installed locally.
