## Requirements

### Requirement: Build target
The Makefile SHALL provide a `build` target that compiles the binary to `./wip`.

#### Scenario: Build produces binary
- **WHEN** user runs `make build`
- **THEN** a `wip` binary is produced at the repo root

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
