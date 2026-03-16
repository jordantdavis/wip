## Why

`wip` is ready to be shared with the team, but there is no version identity, no release process, and no distribution mechanism. Adding versioning and releases makes `wip` installable via `go install`, distributable as a pre-built binary, and clearly identified in bug reports.

## What Changes

- Add `wip version` subcommand that prints the version string and current OS/architecture
- Embed version at build time via `-ldflags`; dev builds use a `dev-<timestamp>` format; release builds use the git tag (e.g. `v0.0.1`)
- Add `go install` support via `debug.ReadBuildInfo()` fallback for users who install from source
- Extend Makefile with `VERSION` variable (defaulting to `dev-<timestamp>`), `LDFLAGS`, and a `release-darwin-arm64` cross-compile target
- Add `.github/workflows/release.yml` that triggers on `v*` tags, builds a darwin/arm64 binary, archives it as a `.tar.gz`, generates SHA256 checksums, and publishes a GitHub Release with auto-generated notes

## Capabilities

### New Capabilities

- `version-command`: The `wip version` subcommand — version string embedding, build info fallback, and OS/architecture display
- `release-workflow`: GitHub Actions release workflow triggered by semver tags — cross-compile, archive, checksum, and publish to GitHub Releases

### Modified Capabilities

- `makefile`: Makefile gains `VERSION`, `LDFLAGS`, and `release-darwin-arm64` target

## Impact

- `main.go`: new `case "version"` routing entry
- `cmd/version.go`: new file
- `Makefile`: extended with version injection and release targets
- `.github/workflows/release.yml`: new workflow file
- `go.mod` / `go.sum`: no new dependencies (`runtime/debug` is stdlib)
