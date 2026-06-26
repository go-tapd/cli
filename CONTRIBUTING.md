# Contributing

This document covers maintenance tasks that are easy to forget between
releases. For regular code changes, run the local checks before opening a pull
request:

```bash
make lint
make test
```

## Release Credentials

The release workflow publishes three outputs when a `v*` tag is pushed:

- GitHub Release assets, using the repository `GITHUB_TOKEN`
- Homebrew formula updates, using `TAP_GITHUB_TOKEN`
- npm package updates, using npm Trusted Publishing with GitHub Actions OIDC

`GITHUB_TOKEN` is created automatically by GitHub Actions. `TAP_GITHUB_TOKEN`
must be created manually and saved as a repository secret in `go-tapd/cli`.
npm publishing does not use a long-lived repository secret; it depends on a
trusted publisher configuration on npmjs.com.

### Required Repository Secrets

| Secret | Used for | Minimum access |
| --- | --- | --- |
| `TAP_GITHUB_TOKEN` | Commit the generated Homebrew formula to `go-tapd/homebrew-tap` | GitHub fine-grained personal access token with `Contents: Read and write` on `go-tapd/homebrew-tap` |

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

### Configure npm Trusted Publishing

Configure trusted publishing from npm:

1. Open `https://www.npmjs.com/package/@go-tapd/tapd`.
2. Open the package settings.
3. Find the Trusted Publishing or Trusted Publisher section.
4. Choose GitHub Actions as the publisher.
5. Set organization or user to `go-tapd`.
6. Set repository to `cli`.
7. Set workflow filename to `release.yml`.
8. Leave environment empty unless the release workflow starts using a GitHub
   environment.
9. Under Allowed actions, select `npm publish`.
10. Click Add publisher to save the trusted publisher configuration.

Why this permission is needed: the release workflow runs `npm publish --access
public` from the `npm/` package directory. npm exchanges the GitHub Actions OIDC
identity for a short-lived publish token, so no long-lived npm repository secret
is required. Trusted publishing automatically generates npm provenance
attestations from GitHub Actions, so the workflow does not need to pass
`--provenance`.

### Verify Release Configuration

Check that the required GitHub repository secret exists:

```bash
gh secret list
```

Expected name:

```text
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

Manually re-open the npm package settings and confirm the trusted publisher
entry still targets `go-tapd/cli`, workflow `release.yml`, and the `npm publish`
allowed action.

The dry-run does not prove that the trusted publisher is configured correctly.
It only proves the local npm account and package metadata are valid. The trusted
publisher is exercised by the GitHub Actions release workflow.

### Rotate TAP_GITHUB_TOKEN

When `TAP_GITHUB_TOKEN` is close to expiration:

1. Create a replacement token using the same permissions listed above.
2. Update the repository secret in
   `https://github.com/go-tapd/cli/settings/secrets/actions`.
3. Re-run the verification commands.
4. Revoke the old token from GitHub after the new secret is saved.

Do not commit tokens, `.npmrc` files containing tokens, or screenshots that
show token values.

### References

- GitHub Docs: [Managing your personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
- npm Docs: [Creating and viewing access tokens](https://docs.npmjs.com/creating-and-viewing-access-tokens/)
- npm Docs: [Requiring 2FA for package publishing and settings modification](https://docs.npmjs.com/requiring-2fa-for-package-publishing-and-settings-modification/)
- npm Docs: [Generating provenance statements](https://docs.npmjs.com/generating-provenance-statements/)
- npm Docs: [Trusted publishing for npm packages](https://docs.npmjs.com/trusted-publishers/)
