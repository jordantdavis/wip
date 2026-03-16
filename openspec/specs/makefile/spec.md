## Requirements

### Requirement: Build target
The Makefile SHALL provide a `build` target that compiles the binary to `./wip` with `LDFLAGS` applied.

#### Scenario: Build produces binary
- **WHEN** user runs `make build`
- **THEN** a `wip` binary is produced at the repo root

#### Scenario: Build injects version
- **WHEN** user runs `make build`
- **THEN** the binary is compiled with `-ldflags "$(LDFLAGS)"` so `wip version` reflects the current `VERSION`

### Requirement: Test target
The Makefile SHALL provide a `test` target that runs all Go tests.

#### Scenario: Tests run
- **WHEN** user runs `make test`
- **THEN** `go test ./...` is executed across all packages

### Requirement: Format target
The Makefile SHALL provide a `fmt` target that reformats all Go source files in place.

#### Scenario: Files are reformatted
- **WHEN** user runs `make fmt`
- **THEN** `go fmt ./...` rewrites any unformatted files

### Requirement: Format check target
The Makefile SHALL provide a `fmt-check` target that fails if any Go source files are not formatted, without modifying them.

#### Scenario: Unformatted files cause failure
- **WHEN** user runs `make fmt-check` and one or more files are not formatted
- **THEN** the target exits with a non-zero status and lists the unformatted files

#### Scenario: Formatted files pass
- **WHEN** user runs `make fmt-check` and all files are formatted
- **THEN** the target exits with status 0

### Requirement: Vet target
The Makefile SHALL provide a `vet` target that runs `go vet` across all packages.

#### Scenario: Vet runs
- **WHEN** user runs `make vet`
- **THEN** `go vet ./...` is executed

### Requirement: Install target
The Makefile SHALL provide an `install` target that builds and installs the binary into `$GOPATH/bin`.

#### Scenario: Binary is installed
- **WHEN** user runs `make install`
- **THEN** `go install .` places the `wip` binary in `$GOPATH/bin`

### Requirement: Clean target
The Makefile SHALL provide a `clean` target that removes the built binary from the repo root.

#### Scenario: Binary is removed
- **WHEN** user runs `make clean`
- **THEN** `./wip` is deleted if it exists

### Requirement: Check target
The Makefile SHALL provide a `check` target that runs `fmt-check`, `vet`, and `test` in sequence.

#### Scenario: Check fails on unformatted code
- **WHEN** user runs `make check` and files are unformatted
- **THEN** the target fails at `fmt-check` before running vet or test

#### Scenario: Check passes when all gates pass
- **WHEN** user runs `make check` and code is formatted, vet-clean, and tests pass
- **THEN** the target exits with status 0

### Requirement: No default target
The Makefile SHALL NOT define a default target. Running bare `make` SHALL produce an error.

#### Scenario: Bare make fails
- **WHEN** user runs `make` with no arguments
- **THEN** make exits with an error indicating no default target

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
