## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip worktree list` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip worktree list` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Worktrees are discovered by scanning the filesystem
`wip worktree list` SHALL discover worktrees by reading the `worktrees/` directory at the repo root. For each subdirectory (treated as a ref name), it SHALL enumerate its subdirectories (treated as worktree names).

#### Scenario: worktrees/ directory does not exist
- **WHEN** the `worktrees/` directory does not exist at the repo root
- **THEN** the CLI prints "no worktrees found" and exits with code 0

#### Scenario: worktrees/ directory is empty
- **WHEN** `worktrees/` exists but contains no subdirectories
- **THEN** the CLI prints "no worktrees found" and exits with code 0

#### Scenario: Submodule directory has no worktrees
- **WHEN** `worktrees/<ref>/` exists but contains no subdirectories
- **THEN** that ref produces no output rows

### Requirement: Output is a flat three-column table sorted alphabetically
The CLI SHALL print one row per worktree in the format `<ref>  <path-segment>  <branch>` (two spaces as separator between each column), sorted first by ref name then by path segment. The branch column SHALL contain the branch name currently checked out in that worktree, obtained by running `git -C <abs-worktree-path> branch --show-current`. If the worktree is in detached HEAD state, the branch column SHALL be empty.

#### Scenario: Single worktree with simple branch name
- **WHEN** `worktrees/my-lib/feature-x/` exists and has branch `feature-x` checked out
- **THEN** the CLI prints `my-lib  feature-x  feature-x`

#### Scenario: Single worktree with slash-delimited branch name
- **WHEN** `worktrees/my-lib/feature-my-thing/` exists and has branch `feature/my-thing` checked out
- **THEN** the CLI prints `my-lib  feature-my-thing  feature/my-thing`

#### Scenario: Multiple worktrees across submodules
- **WHEN** worktrees exist for multiple submodules
- **THEN** the CLI prints all rows sorted by ref name, then path segment within each ref, with the branch column populated for each

#### Scenario: Worktree in detached HEAD state
- **WHEN** a worktree directory exists but its HEAD is detached (not on a branch)
- **THEN** the row is printed with an empty branch column
