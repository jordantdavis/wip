## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any ref operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip ref add` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip ref add` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: URL is a required positional argument
`wip ref add` SHALL accept the repository URL as the first positional argument. If the URL is absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: URL provided
- **WHEN** the user runs `wip ref add https://github.com/org/repo`
- **THEN** the command proceeds with that URL

#### Scenario: URL omitted
- **WHEN** the user runs `wip ref add` with no arguments
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

### Requirement: Name is an optional flag that sets the ref identity and checkout directory
If provided via `--name`, the value SHALL be used as both the git-internal submodule name and the checkout directory at the repository root. The name SHALL NOT contain path separators (`/` or `\`).

#### Scenario: No name provided — git default behavior
- **WHEN** the user runs `wip ref add <url>` without `--name`
- **THEN** the command derives the name from the URL basename (stripping `.git`) and uses it as the checkout directory

#### Scenario: Name provided
- **WHEN** the user provides `--name api`
- **THEN** the command checks out the repo at `./api` with internal name `api`

#### Scenario: Name with path separator rejected
- **WHEN** the user provides `--name libs/foo`
- **THEN** the CLI prints an error indicating the name must not contain path separators and exits with a non-zero code

### Requirement: Branch is an optional flag defaulting to main
`wip ref add` SHALL accept an optional `--branch` flag. If omitted, the branch defaults to `main`. The branch value SHALL be stored in `.gitmodules` and used by `wip ref sync` and `wip ref restore` to determine which remote branch to track.

#### Scenario: No branch provided — defaults to main
- **WHEN** the user runs `wip ref add <url>` without `--branch`
- **THEN** the ref is configured with `branch = main` in `.gitmodules`

#### Scenario: Custom branch provided
- **WHEN** the user runs `wip ref add --branch develop <url>`
- **THEN** the ref is configured with `branch = develop` in `.gitmodules`

### Requirement: Ref is added as a git submodule with branch tracking and ignore = all
After all validations pass, the CLI SHALL execute `git submodule add -b <branch> <url> [<name>]` as a subprocess. After the submodule is added, the CLI SHALL write `ignore = all` to the submodule's entry in `.gitmodules`. The CLI SHALL exit with the same exit code as the git process if it fails.

#### Scenario: Submodule added with branch tracking
- **WHEN** all validations pass and git executes successfully
- **THEN** the ref appears in `.gitmodules` with `branch = <branch>` and `ignore = all` set

#### Scenario: Git command fails
- **WHEN** `git submodule add` returns a non-zero exit code
- **THEN** the CLI exits with the same non-zero code

### Requirement: wip ref add always writes a ref entry to .wip.yml
After a successful `git submodule add`, `wip ref add` SHALL always write an entry to `.wip.yml` under `refs.<name>` containing the `url` and `branch` fields. This ensures `.wip.yml` contains all information needed to restore the ref on a fresh clone without consulting `.gitmodules`.

#### Scenario: Ref entry written
- **WHEN** the user runs `wip ref add <url>`
- **THEN** `.wip.yml` is updated with a `refs.<name>` entry containing `url` and `branch`

### Requirement: Hooks are not configurable via wip ref add
`wip ref add` SHALL NOT accept hook flags. Hook commands (`on-worktree-create`, `on-worktree-launch`) are configured by editing `.wip.yml` directly after the ref is added.

#### Scenario: Hooks added after ref registration
- **WHEN** the user wants to configure hooks for a ref
- **THEN** they edit `.wip.yml` directly, adding `on-worktree-create` and/or `on-worktree-launch` lists under the ref entry
