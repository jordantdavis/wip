## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any submodule operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip submodule add` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip submodule add` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: URL is a required positional argument
`wip submodule add` SHALL accept the repository URL as the first positional argument. If the URL is absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: URL provided
- **WHEN** the user runs `wip submodule add https://github.com/org/repo`
- **THEN** the command proceeds with that URL

#### Scenario: URL omitted
- **WHEN** the user runs `wip submodule add` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: URL must be a valid git remote format
The URL SHALL match at least one of: `https://`, `http://`, `git://`, or SSH form (`git@<host>:<path>`). An invalid URL SHALL cause the CLI to print an error and exit with a non-zero code.

#### Scenario: HTTPS URL accepted
- **WHEN** the URL begins with `https://`
- **THEN** the URL passes validation

#### Scenario: SSH URL accepted
- **WHEN** the URL matches the pattern `git@<host>:<path>`
- **THEN** the URL passes validation

#### Scenario: Invalid URL rejected
- **WHEN** the URL does not match any accepted git remote format
- **THEN** the CLI prints a validation error and exits with a non-zero code

### Requirement: Name is an optional flag that sets the submodule identity and checkout location
If provided via `--name`, the value SHALL be used as both the git-internal submodule name and the checkout directory at the repository root. This allows the same remote URL to be added more than once under distinct names. The name SHALL NOT contain path separators (`/` or `\`).

#### Scenario: No name provided — git default behavior
- **WHEN** the user runs `wip submodule add <url>` without `--name`
- **THEN** the command invokes `git submodule add <url>`, letting git derive the name and checkout directory from the URL

#### Scenario: Name provided
- **WHEN** the user provides `--name foo`
- **THEN** the command invokes `git submodule add --name foo <url> foo`, checking out at `./foo` with internal name `foo`

#### Scenario: Same URL added twice with different names
- **WHEN** the user runs `wip submodule add <url>` and then `wip submodule add --name alt <url>`
- **THEN** both succeed, resulting in two submodule entries with distinct names and checkout directories

#### Scenario: Name with path separator rejected
- **WHEN** the user provides `--name libs/foo`
- **THEN** the CLI prints an error indicating the name must not contain path separators and exits with a non-zero code

### Requirement: CLI executes git submodule add
After all validations pass, the CLI SHALL execute `git submodule add <url> [<path>]` as a subprocess, streaming stdout and stderr to the terminal. The CLI SHALL exit with the same exit code as the git process.

#### Scenario: Successful submodule add
- **WHEN** all validations pass and git executes successfully
- **THEN** the submodule is initialized and the CLI exits with code 0

#### Scenario: Git command fails
- **WHEN** git submodule add returns a non-zero exit code
- **THEN** the CLI exits with the same non-zero code

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

### Requirement: --on-worktree-launch flag accepts ordered hook commands
`wip submodule add` SHALL accept a repeatable `--on-worktree-launch` flag placed before the URL. Each use appends a command string to the list in CLI argument order. The flag MAY be omitted entirely. When provided, the collected commands SHALL be written to `.wip.yml` under `submodules.<name>.on-worktree-launch` as an ordered list.

#### Scenario: Single --on-worktree-launch command
- **WHEN** the user runs `wip submodule add --on-worktree-launch "npm run dev" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-launch: ["npm run dev"]` under the submodule entry

#### Scenario: Multiple --on-worktree-launch commands preserve order
- **WHEN** the user runs `wip submodule add --on-worktree-launch "git pull" --on-worktree-launch "npm install" --on-worktree-launch "claude" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-launch: ["git pull", "npm install", "claude"]` in that order

#### Scenario: --on-worktree-launch omitted
- **WHEN** the user runs `wip submodule add <url>` without `--on-worktree-launch`
- **THEN** no `on-worktree-launch` entry is written to `.wip.yml` for this submodule

#### Scenario: Both --on-worktree-create and --on-worktree-launch provided
- **WHEN** the user provides both `--on-worktree-create` and `--on-worktree-launch` flags
- **THEN** both lists are written to `.wip.yml` under the same submodule entry independently
