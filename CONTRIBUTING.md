# Contributing

This document covers maintenance tasks that are easy to forget between
releases. For regular code changes, run the local checks before opening a pull
request:

```bash
make lint
make test
```

## Release Secrets

The release workflow publishes three outputs when a `v*` tag is pushed:

- GitHub Release assets, using the repository `GITHUB_TOKEN`
- Homebrew formula updates, using `TAP_GITHUB_TOKEN`
- npm package updates, using `NPM_TOKEN`

`GITHUB_TOKEN` is created automatically by GitHub Actions. The other two tokens
must be created manually and saved as repository secrets in `go-tapd/cli`.

### Required Secrets

| Secret | Used for | Minimum access |
| --- | --- | --- |
| `TAP_GITHUB_TOKEN` | Commit the generated Homebrew formula to `go-tapd/homebrew-tap` | GitHub fine-grained personal access token with `Contents: Read and write` on `go-tapd/homebrew-tap` |
| `NPM_TOKEN` | Publish `@go-tapd/tapd` to npm | npm granular access token with `Read and write` access to the `@go-tapd` package scope |

### Configure `TAP_GITHUB_TOKEN`

Create this token from GitHub:

1. Open `https://github.com/settings/personal-access-tokens`.
2. Click `Generate new token`.
3. Choose `Fine-grained personal access token`.
4. Set a descriptive name, such as `go-tapd-homebrew-release`.
5. Set an expiration date and note it somewhere you will check before the next
   release.
6. Set `Resource owner` to `go-tapd`.
7. Set repository access to `Only select repositories`.
8. Select `go-tapd/homebrew-tap`.
9. Under `Repository permissions`, set `Contents` to `Read and write`.
10. Leave unrelated permissions as `No access`.
11. Generate the token and copy it immediately.

Save it in the CLI repository:

1. Open `https://github.com/go-tapd/cli/settings/secrets/actions`.
2. Click `New repository secret`.
3. Set `Name` to `TAP_GITHUB_TOKEN`.
4. Paste the token into `Secret`.
5. Save the secret.

Why this permission is needed: GoReleaser writes the generated
`Formula/tapd.rb` file into `go-tapd/homebrew-tap`. That operation is a content
write to the tap repository. The token does not need access to issues,
pull requests, actions, administration, or any repository other than
`go-tapd/homebrew-tap`.

### Configure `NPM_TOKEN`

Create this token from npm:

1. Open `https://www.npmjs.com/settings/flc1125/tokens`.
2. Click `Generate New Token`.
3. Use a descriptive name, such as `go-tapd-cli-release`.
4. Add a description, such as `Publish @go-tapd/tapd from GitHub Actions`.
5. If npm asks for token type, choose a granular access token.
6. In `Packages and scopes`, set permission to `Read and write`.
7. Select the `@go-tapd` scope. If the scope cannot be selected before the
   first package exists, temporarily select all packages, publish once, then
   replace the token with a narrower `@go-tapd` token.
8. Do not rely on organization permissions alone. Organization access controls
   teams and settings, but package publishing requires package or scope access.
9. Set an expiration date and note it somewhere you will check before the next
   release.
10. If the account or package requires two-factor authentication for publish
    actions, enable `Bypass two-factor authentication` for this CI token.
11. Generate the token and copy it immediately.

Save it in the CLI repository:

1. Open `https://github.com/go-tapd/cli/settings/secrets/actions`.
2. Click `New repository secret`.
3. Set `Name` to `NPM_TOKEN`.
4. Paste the token into `Secret`.
5. Save the secret.

Why this permission is needed: the release workflow runs
`npm publish --access public --provenance` from the `npm/` package directory.
The token must be able to publish `@go-tapd/tapd`. It does not need access to
private packages or unrelated npm scopes.

### Verify Secret Configuration

Check that both GitHub repository secrets exist:

```bash
gh secret list
```

Expected names:

```text
NPM_TOKEN
TAP_GITHUB_TOKEN
```

Check local npm ownership and package dry-run behavior:

```bash
cd npm
npm whoami
npm org ls go-tapd
npm publish --access public --dry-run
```

Expected results:

- `npm whoami` prints the npm user that owns or can publish under `go-tapd`.
- `npm org ls go-tapd` lists that user as an owner or member with publish
  rights.
- `npm publish --access public --dry-run` prints the package contents and ends
  with `+ @go-tapd/tapd@...`.

The dry-run does not prove that `NPM_TOKEN` itself has the correct access. It
only proves the local npm account and package metadata are valid. The token is
exercised by the GitHub Actions release workflow.

### Rotate Expired Tokens

When a token is close to expiration:

1. Create a replacement token using the same permissions listed above.
2. Update the matching repository secret in
   `https://github.com/go-tapd/cli/settings/secrets/actions`.
3. Re-run the verification commands.
4. Revoke the old token from GitHub or npm after the new secret is saved.

Do not commit tokens, `.npmrc` files containing tokens, or screenshots that
show token values.

### References

- GitHub Docs: [Managing your personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
- npm Docs: [Creating and viewing access tokens](https://docs.npmjs.com/creating-and-viewing-access-tokens/)
- npm Docs: [Requiring 2FA for package publishing and settings modification](https://docs.npmjs.com/requiring-2fa-for-package-publishing-and-settings-modification/)
- npm Docs: [Generating provenance statements](https://docs.npmjs.com/generating-provenance-statements/)
