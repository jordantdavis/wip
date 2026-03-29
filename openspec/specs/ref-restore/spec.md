## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any restore operation, `wip ref restore` SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip ref restore` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip ref restore` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Restore initializes and populates all registered refs
`wip ref restore` SHALL initialize and populate all refs registered in `.gitmodules` by invoking `git submodule update --init --remote` as a subprocess. Using `--remote` ensures each ref is checked out at the latest commit on its tracked branch rather than the SHA committed in the parent repo.

#### Scenario: Refs initialized from scratch after clone
- **WHEN** a teammate clones the wip project and runs `wip ref restore`
- **THEN** all registered refs are cloned and checked out at the latest commit on their tracked branch

#### Scenario: Already-initialized refs are updated
- **WHEN** `wip ref restore` is run on a project where refs are already initialized
- **THEN** each ref is updated to the latest commit on its tracked branch (idempotent behavior)

#### Scenario: No refs registered
- **WHEN** `.gitmodules` contains no submodule entries or does not exist
- **THEN** the CLI prints a message indicating there are no refs to restore and exits with code 0

### Requirement: Restore runs all refs concurrently
`wip ref restore` SHALL launch initialization operations for all refs concurrently and wait for all to complete before producing output.

#### Scenario: Multiple refs initialized concurrently
- **WHEN** two or more refs are registered
- **THEN** their initialization operations are launched concurrently rather than sequentially

### Requirement: Output is buffered and reported after all operations complete
`wip ref restore` SHALL NOT stream per-ref output to the terminal during execution. After all operations finish, the CLI SHALL report one result line per ref: `✓ <name>` for success or `✗ <name>: <error>` for failure.

#### Scenario: All succeed
- **WHEN** all ref initializations succeed
- **THEN** the CLI prints `✓ <name>` for each ref after all complete

#### Scenario: Some fail
- **WHEN** one or more ref initializations fail
- **THEN** the CLI prints `✗ <name>: <error>` for each failed ref with captured error detail

### Requirement: Exit code is non-zero if any ref initialization failed
`wip ref restore` SHALL exit with code 0 only if every initialization succeeded. If one or more failed, it SHALL exit with a non-zero code.

#### Scenario: All succeed — exit code 0
- **WHEN** every ref initialization returns success
- **THEN** the CLI exits with code 0

#### Scenario: Any failure — non-zero exit code
- **WHEN** at least one ref initialization fails
- **THEN** the CLI exits with a non-zero code
