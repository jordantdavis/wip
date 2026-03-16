## ADDED Requirements

### Requirement: Version subcommand
`wip` SHALL provide a `version` subcommand that prints the binary's version string and the OS/architecture it was compiled for, then exits with status 0.

#### Scenario: Version output format
- **WHEN** user runs `wip version`
- **THEN** output is `wip <version> (<os>/<arch>)` followed by a newline (e.g. `wip v0.0.1 (darwin/arm64)`)

#### Scenario: Version exits cleanly
- **WHEN** user runs `wip version`
- **THEN** the process exits with status 0

### Requirement: Version string — release build
When the binary is built with `-ldflags "-X github.com/jordantdavis/wip/cmd.version=<tag>"`, the version string SHALL reflect the injected value.

#### Scenario: Release version is displayed
- **WHEN** the binary was built with `VERSION=v0.0.1` via `make build` or `make release-darwin-arm64`
- **THEN** `wip version` prints `wip v0.0.1 (darwin/arm64)`

### Requirement: Version string — dev build
When built via `make build` without an explicit `VERSION`, the version string SHALL be `dev-<timestamp>` where timestamp is `YYYYMMddHHMMss` at build time.

#### Scenario: Dev build shows timestamp
- **WHEN** user runs `make build` with no `VERSION` argument
- **THEN** `wip version` prints `wip dev-<timestamp> (darwin/arm64)`

### Requirement: Version string — go install fallback
When the binary is installed via `go install` and no ldflags version was injected, the version string SHALL be read from the binary's embedded build info.

#### Scenario: go install surfaces module version
- **WHEN** user runs `go install github.com/jordantdavis/wip@v0.0.1`
- **THEN** `wip version` prints `wip v0.0.1 (darwin/arm64)`

### Requirement: Version string — unknown fallback
When neither ldflags injection nor build info provides a version, the version string SHALL be `"unknown"`.

#### Scenario: Raw go build shows unknown
- **WHEN** user runs `go build .` directly (without make)
- **THEN** `wip version` prints `wip unknown (darwin/arm64)`
