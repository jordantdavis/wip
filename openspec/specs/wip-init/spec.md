## ADDED Requirements

### Requirement: Init initializes a git repository in the current directory
`wip init` SHALL use `git rev-parse --git-dir` to determine the state of the current directory and act accordingly:
- If the command exits non-zero, no git repo exists — SHALL run `git init`
- If the command returns `.git`, the current directory is a git repo root — SHALL continue without error or output
- If the command returns anything other than `.git`, the current directory is inside a git repo but not at its root — SHALL print an error and exit with a non-zero code

#### Scenario: No git repo present
- **WHEN** the user runs `wip init` in a directory that is not inside any git repo
- **THEN** `git init` is executed in the current directory and exits successfully

#### Scenario: Already at git repo root
- **WHEN** the user runs `wip init` in a directory that is the root of a git repo
- **THEN** the command completes successfully without running `git init` and without printing any output

#### Scenario: Inside a git repo but not at root
- **WHEN** the user runs `wip init` from a subdirectory inside an existing git repo
- **THEN** an error is printed to stderr indicating the current directory is not at the root, and the command exits with a non-zero code

#### Scenario: Git init failure
- **WHEN** the user runs `wip init` and `git init` fails (e.g. permissions issue)
- **THEN** the error is printed to stderr and the command exits with a non-zero code

### Requirement: Init is idempotent
Running `wip init` multiple times from the same git repo root SHALL produce the same result as running it once. Subsequent runs SHALL NOT fail, overwrite state, or produce spurious output. This includes `.wip.yml`: if it already exists, it SHALL be left unchanged.

#### Scenario: Repeated invocation
- **WHEN** the user runs `wip init` twice from the same git repo root
- **THEN** both invocations exit with code 0, the second produces no output, and `.wip.yml` is not overwritten
