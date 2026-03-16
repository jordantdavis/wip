## ADDED Requirements

### Requirement: Release trigger
The release workflow SHALL trigger automatically when a tag matching `v*` is pushed to the repository.

#### Scenario: Semver tag triggers release
- **WHEN** a tag like `v0.0.1` is pushed to the repository
- **THEN** the release workflow starts

#### Scenario: Non-tag push does not trigger release
- **WHEN** a commit is pushed to `main` without a tag
- **THEN** the release workflow does NOT run

### Requirement: darwin/arm64 binary build
The release workflow SHALL build a `wip` binary for `darwin/arm64` with the tag name injected as the version string via ldflags.

#### Scenario: Binary is built with correct version
- **WHEN** the release workflow runs for tag `v0.0.1`
- **THEN** the produced binary reports `wip v0.0.1 (darwin/arm64)` when `wip version` is run

### Requirement: Archive artifact
The release workflow SHALL package the binary into a `.tar.gz` archive named `wip-<tag>-darwin-arm64.tar.gz`.

#### Scenario: Archive is created
- **WHEN** the release workflow runs for tag `v0.0.1`
- **THEN** a file named `wip-v0.0.1-darwin-arm64.tar.gz` is produced containing the `wip` binary

### Requirement: Checksums file
The release workflow SHALL produce a `checksums.txt` file containing the SHA256 hash of each release archive.

#### Scenario: Checksum is generated
- **WHEN** the release workflow runs
- **THEN** `checksums.txt` contains one line per archive in the format `<sha256>  <filename>`

### Requirement: GitHub Release publication
The release workflow SHALL create a GitHub Release for the tag, attach the archive and `checksums.txt`, and use auto-generated release notes.

#### Scenario: Release is published with artifacts
- **WHEN** the release workflow completes successfully
- **THEN** a GitHub Release exists for the tag with `wip-<tag>-darwin-arm64.tar.gz` and `checksums.txt` attached and release notes auto-generated from commits
