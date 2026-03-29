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
`wip worktree remove` SHALL accept the ref name as the first positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Ref name provided
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the command proceeds with `my-lib` as the ref name

#### Scenario: Ref name omitted
- **WHEN** the user runs `wip worktree remove` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Ref must exist
The submodule name SHALL be validated against `.gitmodules`. If the submodule is not registered, the CLI SHALL print an error and exit with a non-zero code.

#### Scenario: Ref exists
- **WHEN** the ref name is registered in `.gitmodules`
- **THEN** the command proceeds past the ref existence check

#### Scenario: Ref does not exist
- **WHEN** the ref name is not registered in `.gitmodules`
- **THEN** the CLI prints an error indicating the ref was not found and exits with a non-zero code

### Requirement: Worktree name is a required positional argument
`wip worktree remove` SHALL accept the worktree name as the second positional argument. If absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Worktree name provided
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the command proceeds with `my-feature` as the worktree name

#### Scenario: Worktree name omitted
- **WHEN** the user runs `wip worktree remove my-lib` with no second argument
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Worktree name must be a valid git branch name
The worktree argument is treated as a branch name. The CLI SHALL validate it by running `git check-ref-format --branch <name>`. If the command exits with a non-zero code, the CLI SHALL print an error and exit with a non-zero code.

#### Scenario: Valid simple branch name
- **WHEN** the worktree argument is `my-feature`
- **THEN** the name passes validation

#### Scenario: Valid slash-delimited branch name
- **WHEN** the worktree argument is `feature/my-thing`
- **THEN** the name passes validation

#### Scenario: Invalid branch name
- **WHEN** the worktree argument contains characters disallowed by git
- **THEN** the CLI prints a validation error and exits with a non-zero code

### Requirement: Worktree path is derived from the branch name by replacing slashes
The worktree directory path SHALL be computed by replacing every `/` character in the branch name with `-`. The CLI SHALL look up and operate on `<repo root>/worktrees/<ref>/<derived-path>/`.

#### Scenario: Branch name with no slashes
- **WHEN** the branch name is `my-feature`
- **THEN** the CLI looks for the worktree at `worktrees/<ref>/my-feature`

#### Scenario: Branch name with slash
- **WHEN** the branch name is `feature/my-thing`
- **THEN** the CLI looks for the worktree at `worktrees/<ref>/feature-my-thing`

### Requirement: Worktree path must exist
The CLI SHALL verify that the derived worktree path exists before invoking git. If the path does not exist, the CLI SHALL print an error and exit with a non-zero code.

#### Scenario: Worktree path exists
- **WHEN** `worktrees/<ref>/<worktree>/` exists on the filesystem
- **THEN** the command proceeds to git execution

#### Scenario: Worktree path does not exist
- **WHEN** `worktrees/<ref>/<worktree>/` does not exist
- **THEN** the CLI prints an error indicating the worktree was not found and exits with a non-zero code

### Requirement: CLI removes the worktree at the derived path
The CLI SHALL execute `git worktree remove <abs-derived-path>` with its working directory set to `<repo root>/<ref>/`. The CLI SHALL stream stdout and stderr to the terminal and exit with git's exit code.

#### Scenario: Worktree with slash-delimited branch name removed successfully
- **WHEN** the user runs `wip worktree remove my-lib feature/my-thing`
- **THEN** the CLI runs `git worktree remove <abs-path>/worktrees/my-lib/feature-my-thing` inside `my-lib/` and exits with code 0

#### Scenario: Worktree with simple branch name removed successfully
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the CLI runs `git worktree remove <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

#### Scenario: Git command fails
- **WHEN** git returns a non-zero exit code (e.g., worktree has uncommitted changes)
- **THEN** the CLI exits with the same non-zero code

### Requirement: --delete-branch uses the full original branch name
When `--delete-branch` is provided, the CLI SHALL execute `git branch -d <branch-name>` using the original branch name argument (including any `/` characters), not the derived path segment.

#### Scenario: Branch with slash deleted successfully
- **WHEN** `--delete-branch` is provided and the branch `feature/my-thing` is fully merged
- **THEN** the CLI removes the worktree at `feature-my-thing` and runs `git branch -d feature/my-thing`, exiting with code 0

#### Scenario: Branch has unmerged commits
- **WHEN** `--delete-branch` is provided but the branch has unmerged commits
- **THEN** git refuses to delete the branch and the CLI exits with a non-zero code

#### Scenario: Default behavior preserves the branch
- **WHEN** `--delete-branch` is not provided
- **THEN** only the worktree is removed; the branch remains in the ref repo
