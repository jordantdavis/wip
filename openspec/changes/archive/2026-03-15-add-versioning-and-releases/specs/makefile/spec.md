## ADDED Requirements

### Requirement: VERSION variable
The Makefile SHALL define a `VERSION` variable that defaults to `dev-<timestamp>` (format: `YYYYMMddHHMMss`) evaluated at `make` invocation time. It SHALL be overridable by passing `VERSION=<value>` on the command line.

#### Scenario: Default VERSION is a dev timestamp
- **WHEN** user runs `make build` without specifying `VERSION`
- **THEN** `VERSION` resolves to a string of the form `dev-20260315143022`

#### Scenario: VERSION can be overridden
- **WHEN** user runs `make build VERSION=v0.0.1`
- **THEN** `VERSION` is `v0.0.1`

### Requirement: LDFLAGS variable
The Makefile SHALL define an `LDFLAGS` variable set to `-X github.com/jordantdavis/wip/cmd.version=$(VERSION)`.

#### Scenario: LDFLAGS injects version
- **WHEN** `LDFLAGS` is passed to `go build`
- **THEN** the `cmd.version` variable in the binary is set to the current `VERSION` value

### Requirement: release-darwin-arm64 target
The Makefile SHALL provide a `release-darwin-arm64` target that cross-compiles the binary for `GOOS=darwin GOARCH=arm64` with `LDFLAGS` applied, outputting `wip-darwin-arm64`.

#### Scenario: Cross-compile produces arm64 binary
- **WHEN** user runs `make release-darwin-arm64 VERSION=v0.0.1`
- **THEN** a file named `wip-darwin-arm64` is produced, compiled for darwin/arm64 with version `v0.0.1` embedded

## MODIFIED Requirements

### Requirement: Build target
The Makefile SHALL provide a `build` target that compiles the binary to `./wip` with `LDFLAGS` applied.

#### Scenario: Build produces binary
- **WHEN** user runs `make build`
- **THEN** a `wip` binary is produced at the repo root

#### Scenario: Build injects version
- **WHEN** user runs `make build`
- **THEN** the binary is compiled with `-ldflags "$(LDFLAGS)"` so `wip version` reflects the current `VERSION`
