## ADDED Requirements

### Requirement: on-worktree-launch hooks run when wip worktree launch is invoked
`wip worktree launch` SHALL read the `on-worktree-launch` list for the target submodule from `.wip.yml` and execute each command in order inside the worktree directory. Unlike `on-worktree-create`, these hooks are expected to run every time the user launches a worktree and SHOULD be idempotent.

#### Scenario: on-worktree-launch hook executes successfully
- **WHEN** the submodule has an `on-worktree-launch` list in `.wip.yml` and all commands exit with code 0
- **THEN** all commands run in order inside the worktree directory and the CLI exits with code 0

#### Scenario: No on-worktree-launch hook configured
- **WHEN** the submodule has no `on-worktree-launch` entry in `.wip.yml`
- **THEN** `wip worktree launch` prints an informational message and exits with code 0 without running any commands

### Requirement: on-worktree-launch hook failure produces a warning but continues
If any command in the `on-worktree-launch` list exits with a non-zero code, the CLI SHALL print a warning identifying which command failed. The CLI SHALL continue running any remaining commands in the list. The CLI SHALL exit with code 0.

#### Scenario: One hook command fails
- **WHEN** a command in the `on-worktree-launch` list exits non-zero
- **THEN** a warning is printed, subsequent commands still run, and the CLI exits with code 0

### Requirement: on-worktree-launch commands run with cwd set to the worktree directory
Each command in the `on-worktree-launch` list SHALL be executed with its working directory set to the absolute path of the worktree (`<repo root>/worktrees/<submodule>/<worktree>/`).

#### Scenario: Working directory is the worktree
- **WHEN** an `on-worktree-launch` command runs
- **THEN** its working directory is the worktree directory, not the repo root
