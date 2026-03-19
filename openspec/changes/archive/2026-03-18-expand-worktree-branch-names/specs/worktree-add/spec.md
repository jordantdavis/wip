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
- **WHEN** the worktree argument contains characters disallowed by git (e.g., `..`, `@{`, trailing `.`, control characters)
- **THEN** the CLI prints a validation error and exits with a non-zero code

### Requirement: Worktree path is derived from the branch name by replacing slashes
The worktree directory path SHALL be computed by replacing every `/` character in the branch name with `-`. The resulting string is used as the leaf directory name under `<repo root>/worktrees/<submodule>/`.

#### Scenario: Branch name with no slashes
- **WHEN** the branch name is `my-feature`
- **THEN** the worktree path is `worktrees/<submodule>/my-feature`

#### Scenario: Branch name with one slash
- **WHEN** the branch name is `feature/my-thing`
- **THEN** the worktree path is `worktrees/<submodule>/feature-my-thing`

#### Scenario: Branch name with multiple slashes
- **WHEN** the branch name is `team/user/ticket-123`
- **THEN** the worktree path is `worktrees/<submodule>/team-user-ticket-123`

### Requirement: Default behavior creates a new branch using the full branch name
By default, `wip worktree add` SHALL create a new branch using the full branch name argument (including any `/` characters) and check it out at the derived worktree path. The CLI SHALL execute `git worktree add -b <branch-name> <abs-worktree-path>` with its working directory set to `<repo root>/<submodule>/`.

#### Scenario: New branch with slash in name created successfully
- **WHEN** the user runs `wip worktree add my-lib feature/my-thing`
- **THEN** the CLI runs `git worktree add -b feature/my-thing <abs-path>/worktrees/my-lib/feature-my-thing` inside `my-lib/` and exits with code 0

#### Scenario: New branch with simple name created successfully
- **WHEN** the user runs `wip worktree add my-lib my-feature`
- **THEN** the CLI runs `git worktree add -b my-feature <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

### Requirement: --existing-branch checks out the named branch at the derived path
When `--existing-branch` is provided, the CLI SHALL execute `git worktree add <abs-worktree-path> <branch-name>` (without `-b`), where `<abs-worktree-path>` is derived from the branch name by replacing `/` with `-`.

#### Scenario: Existing branch with slash in name checked out successfully
- **WHEN** the user runs `wip worktree add my-lib feature/my-thing --existing-branch` and `feature/my-thing` branch exists
- **THEN** the CLI runs `git worktree add <abs-path>/worktrees/my-lib/feature-my-thing feature/my-thing` inside `my-lib/` and exits with code 0
