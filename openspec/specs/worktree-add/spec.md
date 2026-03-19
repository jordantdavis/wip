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

### Requirement: CLI creates the worktrees directory if it does not exist
Before invoking git, the CLI SHALL create `<repo root>/worktrees/<submodule>/` using `os.MkdirAll` if it does not already exist.

#### Scenario: Directory does not exist
- **WHEN** `worktrees/<submodule>/` does not exist
- **THEN** the CLI creates it before running git

#### Scenario: Directory already exists
- **WHEN** `worktrees/<submodule>/` already exists
- **THEN** the CLI proceeds without error

### Requirement: Default behavior creates a new branch using the full branch name
By default, `wip worktree add` SHALL create a new branch using the full branch name argument (including any `/` characters) and check it out at the derived worktree path. The CLI SHALL execute `git worktree add -b <branch-name> <abs-worktree-path>` with its working directory set to `<repo root>/<submodule>/`.

#### Scenario: New branch with slash in name created successfully
- **WHEN** the user runs `wip worktree add my-lib feature/my-thing`
- **THEN** the CLI runs `git worktree add -b feature/my-thing <abs-path>/worktrees/my-lib/feature-my-thing` inside `my-lib/` and exits with code 0

#### Scenario: New branch with simple name created successfully
- **WHEN** the user runs `wip worktree add my-lib my-feature`
- **THEN** the CLI runs `git worktree add -b my-feature <abs-path>/worktrees/my-lib/my-feature` inside `my-lib/` and exits with code 0

#### Scenario: Git command fails
- **WHEN** git returns a non-zero exit code (e.g., branch already exists)
- **THEN** the CLI exits with the same non-zero code

### Requirement: --existing-branch checks out the named branch at the derived path
When `--existing-branch` is provided, the CLI SHALL execute `git worktree add <abs-worktree-path> <branch-name>` (without `-b`), where `<abs-worktree-path>` is derived from the branch name by replacing `/` with `-`.

#### Scenario: Existing branch with slash in name checked out successfully
- **WHEN** the user runs `wip worktree add my-lib feature/my-thing --existing-branch` and `feature/my-thing` branch exists
- **THEN** the CLI runs `git worktree add <abs-path>/worktrees/my-lib/feature-my-thing feature/my-thing` inside `my-lib/` and exits with code 0

#### Scenario: Branch does not exist
- **WHEN** `--existing-branch` is used but the named branch does not exist in the submodule
- **THEN** git errors and the CLI exits with git's non-zero exit code

### Requirement: CLI streams git output to the terminal
The CLI SHALL stream stdout and stderr from the git subprocess directly to the terminal. The CLI SHALL exit with the same exit code as the git process.

#### Scenario: Output is streamed
- **WHEN** git produces output during worktree creation
- **THEN** the output appears in the terminal in real time

### Requirement: on-worktree-create hooks run after successful worktree creation
After `git worktree add` completes successfully, `wip worktree add` SHALL check `.wip.yml` for an `on-worktree-create` list under the target submodule. If present, each command SHALL be executed in list order with the working directory set to the newly created worktree path.

#### Scenario: on-worktree-create hook executes successfully
- **WHEN** the submodule has an `on-worktree-create` list in `.wip.yml` and all commands exit with code 0
- **THEN** all commands run in order inside the worktree directory and the CLI exits with code 0

#### Scenario: No on-worktree-create hook configured
- **WHEN** the submodule has no `on-worktree-create` entry in `.wip.yml`
- **THEN** `wip worktree add` completes after git without running any hooks

#### Scenario: Submodule has no entry in .wip.yml
- **WHEN** the submodule name does not appear in the `.wip.yml` submodules map
- **THEN** `wip worktree add` completes after git without running any hooks

### Requirement: on-worktree-create hook failure produces a warning but leaves the worktree intact
If any command in the `on-worktree-create` list exits with a non-zero code, the CLI SHALL print a warning to stderr identifying which command failed. The CLI SHALL continue running any remaining commands in the list. The worktree directory SHALL NOT be removed. The CLI SHALL exit with code 0.

#### Scenario: One hook command fails
- **WHEN** a command in the `on-worktree-create` list exits non-zero
- **THEN** a warning is printed to stderr, subsequent commands still run, the worktree remains on disk, and the CLI exits with code 0

#### Scenario: All hook commands fail
- **WHEN** all commands in the `on-worktree-create` list exit non-zero
- **THEN** a warning is printed for each failure, the worktree remains on disk, and the CLI exits with code 0

### Requirement: on-worktree-create commands run with cwd set to the worktree directory
Each command in the `on-worktree-create` list SHALL be executed with its working directory set to the absolute path of the newly created worktree (`<repo root>/worktrees/<submodule>/<worktree>/`).

#### Scenario: Working directory is the worktree
- **WHEN** an `on-worktree-create` command runs
- **THEN** its working directory is the newly created worktree directory, not the repo root
