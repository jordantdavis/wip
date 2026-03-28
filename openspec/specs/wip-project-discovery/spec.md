### Requirement: Project root is discovered by walking up from cwd
`wip` SHALL locate the project root by walking upward from the current working directory, checking each directory for a `.wip.yml` file, and returning the first directory in which one is found.

#### Scenario: .wip.yml found in current directory
- **WHEN** the current directory contains `.wip.yml`
- **THEN** the current directory is used as the project root

#### Scenario: .wip.yml found in a parent directory
- **WHEN** the current directory does not contain `.wip.yml` but an ancestor directory within the user's home does
- **THEN** that ancestor directory is used as the project root

#### Scenario: No .wip.yml found within home directory
- **WHEN** no `.wip.yml` exists in the current directory or any ancestor up to and including the user's home directory
- **THEN** the CLI SHALL print an error directing the user to run `wip init` and exit with a non-zero code

### Requirement: Walk is bounded by the user's home directory
The upward walk SHALL stop after checking the user's home directory. Directories above the home directory SHALL NOT be searched.

#### Scenario: .wip.yml exists at home directory
- **WHEN** the user's home directory contains `.wip.yml`
- **THEN** it is found and the home directory is used as the project root

#### Scenario: .wip.yml would only exist above home
- **WHEN** no `.wip.yml` is found at or below the user's home directory
- **THEN** the command fails with an error, even if a `.wip.yml` exists further up the filesystem

### Requirement: Invocation from outside home directory fails immediately
If the current working directory is not within the user's home directory (i.e., it does not equal home and is not a descendant of home), the CLI SHALL fail immediately without walking the directory tree.

#### Scenario: cwd is outside home
- **WHEN** the user runs a `wip` subcommand from a directory that is not within their home directory (e.g., `/tmp`, `/var`, `/`)
- **THEN** the command exits immediately with a non-zero code and an error message, without performing any directory walk

### Requirement: Subcommands run relative to the project root
Once the project root is discovered, all subcommands (`wip submodule`, `wip worktree`) SHALL execute as if they were invoked from the project root directory. Path construction, git commands, and config reads all resolve relative to the project root.

#### Scenario: Command invoked from a subdirectory
- **WHEN** the user runs `wip worktree list` from a subdirectory of the project
- **THEN** the output is identical to running the same command from the project root

#### Scenario: wip init is exempt
- **WHEN** the user runs `wip init`
- **THEN** project discovery does NOT run; `wip init` operates on the current directory as before
