## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any worktree operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip worktree remove` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip worktree remove` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Submodule name is a required positional argument
`wip worktree remove` SHALL accept the submodule name as the first positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Submodule name provided
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the command proceeds with `my-lib` as the submodule name

#### Scenario: Submodule name omitted
- **WHEN** the user runs `wip worktree remove` with no arguments
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
`wip worktree remove` SHALL accept the worktree name as the second positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Worktree name provided
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the command proceeds with `my-feature` as the worktree name

#### Scenario: Worktree name omitted
- **WHEN** the user runs `wip worktree remove my-lib` with no second argument
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Worktree name must match the allowed character set
The worktree name SHALL only contain uppercase letters, lowercase letters, digits, hyphens, and underscores (`[a-zA-Z0-9_-]`). Any other character SHALL cause the CLI to print an error and exit with a non-zero code.

#### Scenario: Valid worktree name
- **WHEN** the worktree name contains only letters, digits, hyphens, and underscores
- **THEN** the name passes validation

#### Scenario: Worktree name with invalid characters
- **WHEN** the worktree name contains a disallowed character
- **THEN** the CLI prints a validation error and exits with a non-zero code

### Requirement: Worktree path must exist
The CLI SHALL verify that `<repo root>/worktrees/<submodule>/<worktree>/` exists before invoking git. If the path does not exist, the CLI SHALL print an error and exit with a non-zero code.

#### Scenario: Worktree path exists
- **WHEN** `worktrees/<submodule>/<worktree>/` exists on the filesystem
- **THEN** the command proceeds to git execution

#### Scenario: Worktree path does not exist
- **WHEN** `worktrees/<submodule>/<worktree>/` does not exist
- **THEN** the CLI prints an error indicating the worktree was not found and exits with a non-zero code

### Requirement: CLI removes the worktree
The CLI SHALL execute `git worktree remove <abs-path>` with its working directory set to `<repo root>/<submodule>/`, where `<abs-path>` is the absolute path to `worktrees/<submodule>/<worktree>/`. The CLI SHALL stream stdout and stderr to the terminal and exit with git's exit code.

#### Scenario: Worktree removed successfully
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the CLI runs `git worktree remove <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

#### Scenario: Git command fails
- **WHEN** git returns a non-zero exit code (e.g., worktree has uncommitted changes)
- **THEN** the CLI exits with the same non-zero code

### Requirement: --delete-branch flag also deletes the branch
When `--delete-branch` is provided, the CLI SHALL additionally execute `git branch -d <worktree>` inside the submodule directory after successfully removing the worktree. If the branch has unmerged commits, git will refuse with an error and the CLI SHALL exit with that exit code.

#### Scenario: Branch deleted successfully
- **WHEN** `--delete-branch` is provided and the branch is fully merged
- **THEN** the CLI removes the worktree and then deletes the branch, exiting with code 0

#### Scenario: Branch has unmerged commits
- **WHEN** `--delete-branch` is provided but the branch has unmerged commits
- **THEN** git refuses to delete the branch and the CLI exits with a non-zero code

#### Scenario: Default behavior preserves the branch
- **WHEN** `--delete-branch` is not provided
- **THEN** only the worktree is removed; the branch remains in the submodule repo
