## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any sync operation, `wip ref sync` SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip ref sync` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip ref sync` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Default behavior discovers all registered refs from .gitmodules
When no `--name` flag is provided, `wip ref sync` SHALL read `.gitmodules` at the repository root to discover all registered ref names and operate on all of them.

#### Scenario: Refs present in .gitmodules
- **WHEN** `.gitmodules` contains one or more submodule entries
- **THEN** the command collects all ref names and proceeds to update them

#### Scenario: No refs registered
- **WHEN** `.gitmodules` exists but contains no submodule entries, or `.gitmodules` does not exist
- **THEN** the CLI prints a message indicating there are no refs to sync and exits with code 0

### Requirement: Sync updates each ref to latest HEAD on its tracked branch
Each ref SHALL be updated by invoking `git submodule update --remote <name>` as a subprocess. This fetches the remote and checks out the latest commit on the ref's configured branch, ignoring any SHA previously committed in the parent repo.

#### Scenario: Ref updated to latest branch HEAD
- **WHEN** a ref's remote branch has new commits
- **THEN** after sync the ref's working tree reflects the latest commit on that branch

#### Scenario: Ref already at latest — no-op
- **WHEN** a ref is already at the latest commit on its tracked branch
- **THEN** the sync operation for that ref completes successfully with no changes

### Requirement: Default behavior updates all refs concurrently
When no `--name` flag is provided, `wip ref sync` SHALL launch update operations for all refs concurrently and wait for all to complete before producing output.

#### Scenario: Multiple refs run concurrently
- **WHEN** two or more refs are registered
- **THEN** their update operations are launched concurrently rather than sequentially

### Requirement: --name flag limits sync to a single named ref
`wip ref sync` SHALL accept an optional `--name <name>` flag. When provided, only that ref is updated; discovery and concurrency are skipped.

#### Scenario: Named ref updated in isolation
- **WHEN** the user runs `wip ref sync --name api`
- **THEN** the CLI updates only the ref named `api`

### Requirement: Named ref must exist in .gitmodules
When `--name` is provided, the CLI SHALL verify the given name corresponds to a registered ref. If not found, the CLI SHALL print an error and exit with a non-zero code without invoking any git subprocess.

#### Scenario: Named ref found
- **WHEN** the user runs `wip ref sync --name api` and `api` is registered in `.gitmodules`
- **THEN** the command proceeds to update that ref

#### Scenario: Named ref not found
- **WHEN** the user runs `wip ref sync --name foo` and `foo` is not registered
- **THEN** the CLI prints an error and exits with a non-zero code

### Requirement: Output is buffered and reported after all operations complete
`wip ref sync` SHALL NOT stream per-ref output to the terminal during execution. After all update operations finish, the CLI SHALL report one result line per ref. A successful ref SHALL be reported as `✓ <name>`. A failed ref SHALL be reported as `✗ <name>: <error>`.

#### Scenario: All succeed — success lines printed
- **WHEN** all ref updates succeed
- **THEN** the CLI prints `✓ <name>` for each ref after all complete, in any order

#### Scenario: Some fail — failure lines include error detail
- **WHEN** one or more ref updates fail
- **THEN** the CLI prints `✗ <name>: <error>` for each failed ref, with the captured error detail from the git subprocess

#### Scenario: No interleaved output during concurrent runs
- **WHEN** multiple submodules are being updated concurrently
- **THEN** no partial or interleaved output appears on the terminal until all operations have finished

### Requirement: Exit code is non-zero if any ref update failed
After reporting all results, `wip ref sync` SHALL exit with code 0 only if every ref update succeeded. If one or more updates failed, the CLI SHALL exit with a non-zero code.

#### Scenario: All succeed — exit code 0
- **WHEN** every ref update operation returns success
- **THEN** the CLI exits with code 0

#### Scenario: Any failure — non-zero exit code
- **WHEN** at least one ref update operation returns a non-zero exit code
- **THEN** the CLI exits with a non-zero code

#### Scenario: Partial failure — all results still reported
- **WHEN** some refs succeed and others fail
- **THEN** the CLI reports `✓ <name>` for those that succeeded and `✗ <name>: <error>` for those that failed, then exits with a non-zero code
