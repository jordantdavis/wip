## ADDED Requirements

### Requirement: CI runs on push to main
The system SHALL run the full CI suite automatically whenever a commit is pushed to the `main` branch.

#### Scenario: Push to main triggers CI
- **WHEN** a commit is pushed to `main`
- **THEN** the CI workflow runs all three jobs (`build-test`, `fmt-check`, `vet`)

#### Scenario: Push to non-main branch does not trigger CI
- **WHEN** a commit is pushed to any branch other than `main`
- **THEN** the CI workflow does not run

### Requirement: CI runs as PR status checks
The system SHALL run the full CI suite as status checks on every pull request targeting `main`.

#### Scenario: PR targeting main triggers CI
- **WHEN** a pull request is opened or updated with `main` as the base branch
- **THEN** the CI workflow runs all three jobs and each appears as an independent status check on the PR

#### Scenario: PR targeting non-main branch does not trigger CI
- **WHEN** a pull request is opened or updated with a base branch other than `main`
- **THEN** the CI workflow does not run

### Requirement: Build and test job
The CI SHALL run `make build` followed by `make test` as a single job named `build-test`.

#### Scenario: Build and test succeed
- **WHEN** `make build` and `make test` both exit with code 0
- **THEN** the `build-test` job passes

#### Scenario: Build fails
- **WHEN** `make build` exits with a non-zero code
- **THEN** the `build-test` job fails and `make test` does not run

### Requirement: Format check job
The CI SHALL run `make fmt-check` as a standalone job named `fmt-check`.

#### Scenario: All files are formatted
- **WHEN** all Go source files are `gofmt`-formatted
- **THEN** the `fmt-check` job passes

#### Scenario: Unformatted file detected
- **WHEN** one or more Go source files are not `gofmt`-formatted
- **THEN** the `fmt-check` job fails and lists the offending files

### Requirement: Vet job
The CI SHALL run `make vet` as a standalone job named `vet`.

#### Scenario: Vet passes
- **WHEN** `go vet ./...` finds no issues
- **THEN** the `vet` job passes

#### Scenario: Vet finds issues
- **WHEN** `go vet ./...` reports one or more issues
- **THEN** the `vet` job fails

### Requirement: Branch protection blocks merge on failed checks
The `main` branch SHALL require all CI status checks to pass before a pull request can be merged.

#### Scenario: All checks pass
- **WHEN** `build-test`, `fmt-check`, and `vet` all pass on a PR
- **THEN** the PR is eligible to merge

#### Scenario: One or more checks fail
- **WHEN** any of `build-test`, `fmt-check`, or `vet` fails on a PR
- **THEN** the merge button is blocked and the PR cannot be merged
