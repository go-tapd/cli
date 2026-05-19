# TAPD CLI

`tapd` is a command line client for TAPD, built on top of the typed
[`github.com/go-tapd/tapd`](https://github.com/go-tapd/tapd) SDK.

It provides commands for common TAPD resources such as stories, bugs, tasks,
iterations, workspaces, wiki pages, test cases, releases, comments, attachments,
webhooks, and TAPD Lite comments.

## Install

Requires Go 1.25 or later.

```bash
go install github.com/go-tapd/cli/cmd/tapd@latest
```

This installs a binary named `tapd`.

## Quick Start

```bash
tapd login --auth-method pat
tapd auth status
tapd story list --workspace-id 123456 --limit 20
```

Use `tapd --help` and `tapd <command> --help` to inspect available commands and
flags.

## Authentication

The recommended interactive setup is `tapd login`.

Personal access token:

```bash
tapd login --auth-method pat
```

Validate credentials during login:

```bash
tapd login --auth-method pat --workspace-id 123456
```

Basic authentication:

```bash
tapd login --auth-method basic
```

You can also pass credentials directly:

```bash
tapd login --token "$TAPD_ACCESS_TOKEN"
tapd login --client-id "$TAPD_CLIENT_ID" --client-secret "$TAPD_CLIENT_SECRET"
```

The login command stores credentials in:

```text
~/.tapd/config.json
```

For CI or temporary shell sessions, you can skip local config and use
environment variables:

```bash
export TAPD_ACCESS_TOKEN=...
```

or:

```bash
export TAPD_CLIENT_ID=...
export TAPD_CLIENT_SECRET=...
```

## Common Examples

```bash
tapd workspace view --workspace-id 123456
tapd workspace users --workspace-id 123456

tapd story list --workspace-id 123456 --limit 20
tapd story view 1111112222001000001 --workspace-id 123456
tapd story create --workspace-id 123456 --name "New story" --owner alice

tapd bug list --workspace-id 123456 --owner alice
tapd bug create --workspace-id 123456 --title "Login fails" --owner alice

tapd task list --workspace-id 123456 --creator bob --format json
tapd webhook inspect --file payload.json
tapd webhook serve --addr 127.0.0.1:8080 --path /webhook
```

## Command Areas

- Authentication: `tapd login`, `tapd auth status`, `tapd auth logout`
- Workspaces: `tapd workspace ...`
- Stories: `tapd story ...`
- Bugs: `tapd bug ...`
- Tasks: `tapd task ...`
- Iterations: `tapd iteration ...`
- Tests: `tapd test-case ...`, `tapd test-plan ...`
- Releases: `tapd release ...`, `tapd launch-form ...`
- Wiki: `tapd wiki ...`
- Reports and metrics: `tapd report ...`, `tapd measure ...`
- Timesheets: `tapd timesheet ...`
- Comments: `tapd comment ...`
- Attachments: `tapd attachment ...`
- Boards: `tapd board ...`
- Workflow and settings: `tapd workflow ...`, `tapd setting ...`
- Labels and users: `tapd label ...`, `tapd user ...`
- Source integrations: `tapd source ...`
- Webhooks: `tapd webhook ...`
- TAPD Lite: `tapd lite comment ...`

## Documentation

- [Shell completion](docs/shell-completion.md)
- [Workspace commands](docs/workspace.md)
- [Story commands](docs/story.md)
- [Bug commands](docs/bug.md)
- [Task commands](docs/task.md)
- [Iteration commands](docs/iteration.md)
- [Test commands](docs/test.md)
- [Release commands](docs/release.md)
- [Wiki commands](docs/wiki.md)
- [Report commands](docs/report.md)
- [Attachment commands](docs/attachment.md)
- [Measure commands](docs/measure.md)
- [Timesheet commands](docs/timesheet.md)
- [Comment commands](docs/comment.md)
- [Board commands](docs/board.md)
- [Workflow commands](docs/workflow.md)
- [Setting commands](docs/setting.md)
- [Label commands](docs/label.md)
- [User commands](docs/user.md)
- [Source commands](docs/source.md)
- [Webhook commands](docs/webhook.md)
- [Lite commands](docs/lite.md)
- [Feature coverage](features.md)

## Development

```bash
make lint
make test
go build ./...
```

When adding a command, prefer the typed services, requests, and responses from
`github.com/go-tapd/tapd` rather than writing ad hoc HTTP calls.

## License

This project is licensed under the [MIT License](LICENSE).
