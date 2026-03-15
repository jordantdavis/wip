## MODIFIED Requirements

### Requirement: Working directory must be a git repository
Before executing any worktree operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip worktree add` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip worktree add` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code
