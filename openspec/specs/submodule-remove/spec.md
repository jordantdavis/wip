## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any submodule removal operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip submodule remove` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip submodule remove` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Name is a required positional argument
`wip submodule remove` SHALL accept the submodule name as the first positional argument. If the name is absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Name provided
- **WHEN** the user runs `wip submodule remove foo`
- **THEN** the command proceeds with `foo` as the submodule name

#### Scenario: Name omitted
- **WHEN** the user runs `wip submodule remove` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Named submodule must exist before removal is attempted
Before executing any removal steps, the CLI SHALL verify that a submodule with the given name is registered. It SHALL check `.gitmodules` or the git config for a matching submodule entry. If no matching entry is found, the CLI SHALL print an error and exit with a non-zero code without executing any removal steps.

#### Scenario: Submodule exists
- **WHEN** the named submodule is present in `.gitmodules` or git config
- **THEN** the command proceeds to the removal sequence

#### Scenario: Submodule does not exist
- **WHEN** no submodule with the given name is registered
- **THEN** the CLI prints an error indicating the submodule was not found and exits with a non-zero code

### Requirement: Removal executes three steps in order
After validation passes, the CLI SHALL perform full submodule removal by executing the following three operations in order:
1. `git submodule deinit -f <name>` — deinitializes the submodule, removing its entry from `.git/config` and clearing the working tree checkout
2. `git rm -f <name>` — removes the submodule from the working tree and from `.gitmodules`
3. `rm -rf .git/modules/<name>` — removes the cached module data from the git object store

Each step SHALL be completed before the next begins. Stdout and stderr from git commands SHALL be streamed to the terminal in real time.

#### Scenario: All three steps execute in sequence
- **WHEN** all validations pass
- **THEN** the CLI runs `git submodule deinit -f <name>`, then `git rm -f <name>`, then `rm -rf .git/modules/<name>`, in that order

#### Scenario: Output is streamed
- **WHEN** a git command produces output during removal
- **THEN** that output is written to the terminal as it is produced, not buffered until completion

### Requirement: Failure in any step stops execution and reports the error
If any step in the removal sequence exits with a non-zero code or otherwise fails, the CLI SHALL immediately stop, print an error indicating which step failed, and exit with the failing step's exit code. Subsequent steps SHALL NOT be attempted.

#### Scenario: deinit step fails
- **WHEN** `git submodule deinit -f <name>` exits with a non-zero code
- **THEN** the CLI prints an error identifying the deinit step as the failure, does not proceed to `git rm` or `rm -rf`, and exits with the same non-zero code

#### Scenario: git rm step fails
- **WHEN** `git submodule deinit -f <name>` succeeds but `git rm -f <name>` exits with a non-zero code
- **THEN** the CLI prints an error identifying the git rm step as the failure, does not proceed to `rm -rf`, and exits with the same non-zero code

#### Scenario: Module directory removal fails
- **WHEN** the first two steps succeed but `rm -rf .git/modules/<name>` fails
- **THEN** the CLI prints an error identifying the module directory removal as the failure and exits with a non-zero code

### Requirement: Successful removal exits with code 0
If all three removal steps complete without error, the CLI SHALL exit with code 0.

#### Scenario: Full successful removal
- **WHEN** all three steps complete successfully
- **THEN** the CLI exits with code 0
