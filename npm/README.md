# TAPD CLI

This npm package installs the `tapd` command line client from the matching
GitHub Release binary.

```bash
npm install -g @go-tapd/tapd
tapd --help
```

The installer downloads the archive for the current operating system and CPU
architecture, verifies it with the release checksum, and exposes the `tapd`
binary through npm's global bin path.
