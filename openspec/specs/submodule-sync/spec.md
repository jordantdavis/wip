## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any sync operation, `wip submodule sync` SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip submodule sync` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip submodule sync` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Default behavior discovers all registered submodules from .gitmodules
When no `--name` flag is provided, `wip submodule sync` SHALL read `.gitmodules` at the repository root to discover all registered submodule names. The set of discovered names SHALL be the complete list of submodules operated on.

#### Scenario: Submodules present in .gitmodules
- **WHEN** `.gitmodules` contains one or more submodule entries
- **THEN** the command collects all submodule names from `.gitmodules` and proceeds to update them

#### Scenario: No submodules registered
- **WHEN** `.gitmodules` exists but contains no submodule entries, or `.gitmodules` does not exist
- **THEN** the CLI prints a message indicating there are no submodules to sync and exits with code 0

### Requirement: Default behavior updates all submodules concurrently
When no `--name` flag is provided, `wip submodule sync` SHALL launch update operations for all discovered submodules concurrently. Each submodule SHALL be updated by invoking `git submodule update --init --remote <name>` as a subprocess. The command SHALL wait for all concurrent operations to complete before producing any output.

#### Scenario: All submodules update successfully
- **WHEN** all submodule update operations succeed
- **THEN** the CLI waits for all to complete, then reports success for each and exits with code 0

#### Scenario: Multiple submodules run concurrently
- **WHEN** two or more submodules are registered
- **THEN** their update operations are launched concurrently rather than sequentially

### Requirement: --name flag limits sync to a single named submodule
`wip submodule sync` SHALL accept an optional `--name <name>` flag. When provided, the CLI SHALL update only the submodule with that name, skipping discovery and concurrency. The single submodule SHALL be updated by invoking `git submodule update --init --remote <name>`.

#### Scenario: Named submodule updated in isolation
- **WHEN** the user runs `wip submodule sync --name foo`
- **THEN** the CLI updates only the submodule named `foo` and reports its result

#### Scenario: --name flag skips concurrency
- **WHEN** the user provides `--name`
- **THEN** only one submodule update is performed, with no concurrent operations

### Requirement: Named submodule must exist in .gitmodules
When `--name` is provided, the CLI SHALL verify that the given name corresponds to a registered submodule in `.gitmodules`. If the name is not found, the CLI SHALL print an error and exit with a non-zero code without invoking any git subprocess.

#### Scenario: Named submodule found
- **WHEN** the user runs `wip submodule sync --name foo` and `foo` is registered in `.gitmodules`
- **THEN** the command proceeds to update that submodule

#### Scenario: Named submodule not found
- **WHEN** the user runs `wip submodule sync --name bar` and `bar` is not registered in `.gitmodules`
- **THEN** the CLI prints an error indicating the submodule `bar` was not found and exits with a non-zero code

### Requirement: Output is buffered and reported after all operations complete
`wip submodule sync` SHALL NOT stream per-submodule output to the terminal during execution. After all update operations finish, the CLI SHALL report one result line per submodule. A successful submodule SHALL be reported as `✓ <name>`. A failed submodule SHALL be reported as `✗ <name>: <error>`.

#### Scenario: All succeed — success lines printed
- **WHEN** all submodule updates succeed
- **THEN** the CLI prints `✓ <name>` for each submodule after all complete, in any order

#### Scenario: Some fail — failure lines include error detail
- **WHEN** one or more submodule updates fail
- **THEN** the CLI prints `✗ <name>: <error>` for each failed submodule, with the captured error detail from the git subprocess

#### Scenario: No interleaved output during concurrent runs
- **WHEN** multiple submodules are being updated concurrently
- **THEN** no partial or interleaved output appears on the terminal until all operations have finished

### Requirement: Exit code is non-zero if any submodule update failed
After reporting all results, `wip submodule sync` SHALL exit with code 0 only if every submodule update succeeded. If one or more updates failed, the CLI SHALL exit with a non-zero code.

#### Scenario: All succeed — exit code 0
- **WHEN** every submodule update operation returns success
- **THEN** the CLI exits with code 0

#### Scenario: Any failure — non-zero exit code
- **WHEN** at least one submodule update operation returns a non-zero exit code
- **THEN** the CLI exits with a non-zero code

#### Scenario: Partial failure — all results still reported
- **WHEN** some submodules succeed and others fail
- **THEN** the CLI reports `✓ <name>` for those that succeeded and `✗ <name>: <error>` for those that failed, then exits with a non-zero code
