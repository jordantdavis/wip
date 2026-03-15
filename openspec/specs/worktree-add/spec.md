## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any worktree operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip worktree add` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip worktree add` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Submodule name is a required positional argument
`wip worktree add` SHALL accept the submodule name as the first positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Submodule name provided
- **WHEN** the user runs `wip worktree add my-lib my-feature`
- **THEN** the command proceeds with `my-lib` as the submodule name

#### Scenario: Submodule name omitted
- **WHEN** the user runs `wip worktree add` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Submodule must exist
The submodule name SHALL be validated against `.gitmodules`. If the submodule is not registered, the CLI SHALL print an error and exit with a non-zero code.

#### Scenario: Submodule exists
- **WHEN** the submodule name is registered in `.gitmodules`
- **THEN** the command proceeds past the submodule existence check

#### Scenario: Submodule does not exist
- **WHEN** the submodule name is not registered in `.gitmodules`
- **THEN** the CLI prints an error indicating the submodule was not found and exits with a non-zero code

### Requirement: Worktree name is a required positional argument
`wip worktree add` SHALL accept the worktree name as the second positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Worktree name provided
- **WHEN** the user runs `wip worktree add my-lib my-feature`
- **THEN** the command proceeds with `my-feature` as the worktree name

#### Scenario: Worktree name omitted
- **WHEN** the user runs `wip worktree add my-lib` with no second argument
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Worktree name must match the allowed character set
The worktree name SHALL only contain uppercase letters, lowercase letters, digits, hyphens, and underscores (`[a-zA-Z0-9_-]`). Any other character SHALL cause the CLI to print an error and exit with a non-zero code.

#### Scenario: Valid worktree name
- **WHEN** the worktree name contains only letters, digits, hyphens, and underscores
- **THEN** the name passes validation

#### Scenario: Worktree name with invalid characters
- **WHEN** the worktree name contains a space, slash, dot, or other disallowed character
- **THEN** the CLI prints a validation error and exits with a non-zero code

### Requirement: CLI creates the worktrees directory if it does not exist
Before invoking git, the CLI SHALL create `<repo root>/worktrees/<submodule>/` using `os.MkdirAll` if it does not already exist.

#### Scenario: Directory does not exist
- **WHEN** `worktrees/<submodule>/` does not exist
- **THEN** the CLI creates it before running git

#### Scenario: Directory already exists
- **WHEN** `worktrees/<submodule>/` already exists
- **THEN** the CLI proceeds without error

### Requirement: Default behavior creates a new branch
By default, `wip worktree add` SHALL create a new branch named `<worktree>` and check it out at `<repo root>/worktrees/<submodule>/<worktree>/`. The CLI SHALL execute `git worktree add -b <worktree> <abs-path>` with its working directory set to `<repo root>/<submodule>/`.

#### Scenario: New branch worktree created successfully
- **WHEN** the user runs `wip worktree add my-lib my-feature` without `--existing-branch`
- **THEN** the CLI runs `git worktree add -b my-feature <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

#### Scenario: Git command fails
- **WHEN** git returns a non-zero exit code (e.g., branch already exists)
- **THEN** the CLI exits with the same non-zero code

### Requirement: --existing-branch flag checks out an existing branch
When `--existing-branch` is provided, the CLI SHALL execute `git worktree add <abs-path> <worktree>` (without `-b`), checking out the named branch. If the branch does not exist in the submodule repo, git will error and the CLI SHALL exit with that exit code.

#### Scenario: Existing branch checked out successfully
- **WHEN** the user runs `wip worktree add my-lib my-feature --existing-branch` and `my-feature` branch exists
- **THEN** the CLI runs `git worktree add <abs-path>/worktrees/my-lib/my-feature my-feature` inside `my-lib/` and exits with code 0

#### Scenario: Branch does not exist
- **WHEN** `--existing-branch` is used but the named branch does not exist in the submodule
- **THEN** git errors and the CLI exits with git's non-zero exit code

### Requirement: CLI streams git output to the terminal
The CLI SHALL stream stdout and stderr from the git subprocess directly to the terminal. The CLI SHALL exit with the same exit code as the git process.

#### Scenario: Output is streamed
- **WHEN** git produces output during worktree creation
- **THEN** the output appears in the terminal in real time
