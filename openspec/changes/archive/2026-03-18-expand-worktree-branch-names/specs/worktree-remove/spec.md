## MODIFIED Requirements

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
The worktree directory path SHALL be computed by replacing every `/` character in the branch name with `-`. The CLI SHALL look up and operate on `<repo root>/worktrees/<submodule>/<derived-path>/`.

#### Scenario: Branch name with no slashes
- **WHEN** the branch name is `my-feature`
- **THEN** the CLI looks for the worktree at `worktrees/<submodule>/my-feature`

#### Scenario: Branch name with slash
- **WHEN** the branch name is `feature/my-thing`
- **THEN** the CLI looks for the worktree at `worktrees/<submodule>/feature-my-thing`

### Requirement: CLI removes the worktree at the derived path
The CLI SHALL execute `git worktree remove <abs-derived-path>` with its working directory set to `<repo root>/<submodule>/`. The CLI SHALL stream stdout and stderr to the terminal and exit with git's exit code.

#### Scenario: Worktree with slash-delimited branch name removed successfully
- **WHEN** the user runs `wip worktree remove my-lib feature/my-thing`
- **THEN** the CLI runs `git worktree remove <abs-path>/worktrees/my-lib/feature-my-thing` inside `my-lib/` and exits with code 0

#### Scenario: Worktree with simple branch name removed successfully
- **WHEN** the user runs `wip worktree remove my-lib my-feature`
- **THEN** the CLI runs `git worktree remove <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

### Requirement: --delete-branch uses the full original branch name
When `--delete-branch` is provided, the CLI SHALL execute `git branch -d <branch-name>` using the original branch name argument (including any `/` characters), not the derived path segment.

#### Scenario: Branch with slash deleted successfully
- **WHEN** `--delete-branch` is provided and the branch `feature/my-thing` is fully merged
- **THEN** the CLI removes the worktree at `feature-my-thing` and runs `git branch -d feature/my-thing`, exiting with code 0
