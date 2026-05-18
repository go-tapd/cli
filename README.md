# tapd CLI

`tapd` is a TAPD command line client built on top of [`github.com/go-tapd/tapd`](https://github.com/go-tapd/tapd).

## Install

```bash
go install github.com/go-tapd/cli/cmd/tapd@latest
```

This installs a binary named `tapd`.

## Current Commands

- `tapd login`
- `tapd auth status`
- `tapd auth logout`
- `tapd workspace view`
- `tapd workspace users`
- `tapd story create`
- `tapd story view`
- `tapd story list`
- `tapd story count`
- `tapd story update`
- `tapd story batch-update`
- `tapd story categories`
- `tapd story categories count`
- `tapd story changes`
- `tapd story changes count`
- `tapd story fields`
- `tapd story field-labels`
- `tapd story templates`
- `tapd story template-fields`
- `tapd story removed`
- `tapd story related-bugs`
- `tapd story related-test-cases`
- `tapd story by-view`
- `tapd story convert-ids`
- `tapd bug create`
- `tapd bug view`
- `tapd bug list`
- `tapd bug count`
- `tapd bug update`
- `tapd bug batch-update`
- `tapd bug changes`
- `tapd bug changes count`
- `tapd bug fields`
- `tapd bug field-labels`
- `tapd bug templates`
- `tapd bug template-fields`
- `tapd bug removed`
- `tapd bug related-stories`
- `tapd bug link`
- `tapd bug unlink`
- `tapd bug by-view`
- `tapd bug convert-ids`
- `tapd iteration create`
- `tapd iteration view`
- `tapd iteration list`
- `tapd iteration count`
- `tapd iteration update`
- `tapd iteration changes`
- `tapd iteration fields`
- `tapd iteration workitem-types`
- `tapd iteration templates`
- `tapd iteration template-fields`
- `tapd iteration lock`
- `tapd iteration unlock`
- `tapd task create`
- `tapd task view`
- `tapd task list`
- `tapd task count`
- `tapd task update`
- `tapd task batch-update`
- `tapd task changes`
- `tapd task changes count`
- `tapd task fields`
- `tapd task removed`
- `tapd report list`
- `tapd attachment list`
- `tapd attachment download-url`
- `tapd attachment image-url`
- `tapd attachment document-url`
- `tapd measure life-times`
- `tapd workflow last-steps`
- `tapd setting workspace`

## Login

PAT:

```bash
go run ./cmd/tapd login --auth-method pat
```

Validate at login time with a workspace:

```bash
go run ./cmd/tapd login --auth-method pat --workspace-id 123456
```

PAT with flag:

```bash
go run ./cmd/tapd login --token "$TAPD_ACCESS_TOKEN"
```

Basic Authentication:

```bash
go run ./cmd/tapd login --auth-method basic
```

Basic Authentication with flags:

```bash
go run ./cmd/tapd login --client-id "$TAPD_CLIENT_ID" --client-secret "$TAPD_CLIENT_SECRET"
```

The login command validates credentials before saving them to:

```text
~/.tapd/config.json
```

By default, `login` stores the credentials directly. If you want an API-side validation during login, pass `--workspace-id`.

You can also skip local config and use environment variables directly:

```bash
export TAPD_ACCESS_TOKEN=...
```

or

```bash
export TAPD_CLIENT_ID=...
export TAPD_CLIENT_SECRET=...
```

## Examples

```bash
go run ./cmd/tapd auth status
go run ./cmd/tapd workspace view --workspace-id 123456
go run ./cmd/tapd workspace users --workspace-id 123456
go run ./cmd/tapd story list --workspace-id 123456 --limit 20
go run ./cmd/tapd story view 1111112222001000001 --workspace-id 123456
go run ./cmd/tapd story create --workspace-id 123456 --name "New story" --owner alice
go run ./cmd/tapd story changes --workspace-id 123456 --story-ids 1111112222001000001
go run ./cmd/tapd story batch-update --workspace-id 123456 --file stories.json
go run ./cmd/tapd bug list --workspace-id 123456 --owner alice
go run ./cmd/tapd bug view 1111112222001000001 --workspace-id 123456
go run ./cmd/tapd bug create --workspace-id 123456 --title "Login fails" --owner alice
go run ./cmd/tapd bug changes --workspace-id 123456 --bug-ids 1111112222001000001
go run ./cmd/tapd bug link --workspace-id 123456 --bug-id 1111112222001000001 --relate-bug-ids 1111112222001000002
go run ./cmd/tapd iteration list --workspace-id 123456 --status open
go run ./cmd/tapd iteration create --workspace-id 123456 --name "Sprint 1" --description "Sprint goal" --start-date 2026-06-01 --end-date 2026-06-14 --creator alice
go run ./cmd/tapd iteration changes --workspace-id 123456 --iteration-id 1111112222001000001
go run ./cmd/tapd task create --workspace-id 123456 --name "Implement login" --owner alice
go run ./cmd/tapd task view 1111112222001000001 --workspace-id 123456
go run ./cmd/tapd task changes --workspace-id 123456 --task-id 1111112222001000001
go run ./cmd/tapd task list --workspace-id 123456 --creator bob --format json
go run ./cmd/tapd report list --workspace-id 123456 --author alice
go run ./cmd/tapd attachment list --workspace-id 123456 --entry-id 1111112222001000001
go run ./cmd/tapd attachment download-url --workspace-id 123456 --id 10001
go run ./cmd/tapd measure life-times --workspace-id 123456 --entity-type story --entity-id 1111112222001000001
go run ./cmd/tapd workflow last-steps --workspace-id 123456 --group-key workitem_type_id
go run ./cmd/tapd setting workspace --workspace-id 123456
```

## Documentation

- [Shell completion](docs/shell-completion.md)
- [Story commands](docs/story.md)
- [Bug commands](docs/bug.md)
- [Iteration commands](docs/iteration.md)
- [Task commands](docs/task.md)
- [Report commands](docs/report.md)
- [Attachment commands](docs/attachment.md)
- [Measure commands](docs/measure.md)
- [Workflow commands](docs/workflow.md)
- [Setting commands](docs/setting.md)
