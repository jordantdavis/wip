## MODIFIED Requirements

### Requirement: Output is a flat three-column table sorted alphabetically
The CLI SHALL print one row per worktree in the format `<submodule>  <path-segment>  <branch>` (two spaces as separator between each column), sorted first by submodule name then by path segment. The branch column SHALL contain the branch name currently checked out in that worktree, obtained by running `git -C <abs-worktree-path> branch --show-current`. If the worktree is in detached HEAD state, the branch column SHALL be empty.

#### Scenario: Single worktree with simple branch name
- **WHEN** `worktrees/my-lib/feature-x/` exists and has branch `feature-x` checked out
- **THEN** the CLI prints `my-lib  feature-x  feature-x`

#### Scenario: Single worktree with slash-delimited branch name
- **WHEN** `worktrees/my-lib/feature-my-thing/` exists and has branch `feature/my-thing` checked out
- **THEN** the CLI prints `my-lib  feature-my-thing  feature/my-thing`

#### Scenario: Multiple worktrees across submodules
- **WHEN** worktrees exist for multiple submodules
- **THEN** the CLI prints all rows sorted by submodule name, then path segment within each submodule, with the branch column populated for each

#### Scenario: Worktree in detached HEAD state
- **WHEN** a worktree directory exists but its HEAD is detached (not on a branch)
- **THEN** the row is printed with an empty branch column
