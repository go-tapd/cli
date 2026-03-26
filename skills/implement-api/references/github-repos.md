# GitHub Repositories

- CLI repository: `https://github.com/go-tapd/cli`
- SDK repository: `https://github.com/go-tapd/tapd`

Use the CLI repository as the implementation target.

Use the SDK repository as the source of truth for:

- API service definitions
- request structs
- response structs
- authentication capabilities

When implementing a new command, first confirm whether the SDK already exposes the needed TAPD API.
