## ADDED Requirements

### Requirement: wip worktree launch requires submodule and worktree arguments
`wip worktree launch` SHALL accept exactly two positional arguments: `<submodule>` and `<worktree>`. If either is missing, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Both arguments provided
- **WHEN** the user runs `wip worktree launch my-service feat/my-feature`
- **THEN** the command proceeds with submodule `my-service` and worktree `feat/my-feature`

#### Scenario: Worktree argument missing
- **WHEN** the user runs `wip worktree launch my-service` with no worktree argument
- **THEN** the CLI prints usage and exits with a non-zero code

#### Scenario: No arguments provided
- **WHEN** the user runs `wip worktree launch` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: wip worktree launch verifies the worktree directory exists
Before running any hooks, `wip worktree launch` SHALL verify that the worktree directory (`<repo root>/worktrees/<submodule>/<worktree-path-segment>`) exists on disk. If the directory does not exist, the CLI SHALL print an error message and exit with a non-zero code.

#### Scenario: Worktree directory exists
- **WHEN** the worktree directory is present on disk
- **THEN** the command proceeds to hook execution

#### Scenario: Worktree directory does not exist
- **WHEN** the worktree directory is absent
- **THEN** the CLI prints an error indicating the worktree was not found and exits with a non-zero code

### Requirement: wip worktree launch prints an informational message when no hooks are configured
If the submodule has no `on-worktree-launch` entry in `.wip.yml`, `wip worktree launch` SHALL print a message indicating that no hooks are configured for the submodule and exit with code 0.

#### Scenario: No on-worktree-launch hooks configured
- **WHEN** the submodule has no `on-worktree-launch` list in `.wip.yml`
- **THEN** the CLI prints "no on-worktree-launch hooks configured for <submodule>" and exits with code 0

#### Scenario: Submodule has no entry in .wip.yml
- **WHEN** the submodule name does not appear in the `.wip.yml` submodules map
- **THEN** the CLI prints "no on-worktree-launch hooks configured for <submodule>" and exits with code 0

### Requirement: wip worktree launch runs on-worktree-launch hooks sequentially in the worktree directory
`wip worktree launch` SHALL execute each command in the `on-worktree-launch` list in order, with the working directory set to the worktree path. Commands run blocking and sequentially. The final command MAY be interactive and take over the terminal.

#### Scenario: Single hook runs successfully
- **WHEN** the submodule has one `on-worktree-launch` command and it exits with code 0
- **THEN** the command runs in the worktree directory, ✓ is printed, and the CLI exits with code 0

#### Scenario: Multiple hooks run in order
- **WHEN** the submodule has multiple `on-worktree-launch` commands
- **THEN** each command runs sequentially in the worktree directory in list order

#### Scenario: Hook failure prints warning and continues
- **WHEN** a command in the `on-worktree-launch` list exits non-zero
- **THEN** ✗ is printed for that command, subsequent commands still run, and the CLI exits with code 0

#### Scenario: Last command is interactive
- **WHEN** the last `on-worktree-launch` command is an interactive process (e.g. `claude`)
- **THEN** it takes over the terminal and the CLI exits when that process exits
