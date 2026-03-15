## ADDED Requirements

### Requirement: --on-worktree-create flag accepts ordered hook commands
`wip submodule add` SHALL accept a repeatable `--on-worktree-create` flag placed before the URL. Each use appends a command string to the list in CLI argument order. The flag MAY be omitted entirely. When provided, the collected commands SHALL be written to `.wip.yml` under `submodules.<name>.on-worktree-create` as an ordered list.

#### Scenario: Single --on-worktree-create command
- **WHEN** the user runs `wip submodule add --on-worktree-create "npm install" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-create: ["npm install"]` under the submodule entry

#### Scenario: Multiple --on-worktree-create commands preserve order
- **WHEN** the user runs `wip submodule add --on-worktree-create "npm install" --on-worktree-create "cp .env.example .env" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-create: ["npm install", "cp .env.example .env"]` in that order

#### Scenario: --on-worktree-create omitted
- **WHEN** the user runs `wip submodule add <url>` without `--on-worktree-create`
- **THEN** no `on-worktree-create` entry is written to `.wip.yml` for this submodule

## MODIFIED Requirements

### Requirement: Working directory must be a git repository
Before executing any submodule operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip submodule add` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip submodule add` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code
