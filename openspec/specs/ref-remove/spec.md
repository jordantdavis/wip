## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any ref removal operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip ref remove` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip ref remove` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Name is a required positional argument
`wip ref remove` SHALL accept the ref name as the first positional argument. If the name is absent or empty, the CLI SHALL print usage and exit with a non-zero code.

#### Scenario: Name provided
- **WHEN** the user runs `wip ref remove foo`
- **THEN** the command proceeds with `foo` as the ref name

#### Scenario: Name omitted
- **WHEN** the user runs `wip ref remove` with no arguments
- **THEN** the CLI prints usage and exits with a non-zero code

### Requirement: Named ref must exist before removal is attempted
Before executing any removal steps, the CLI SHALL verify that a ref with the given name is registered in `.gitmodules`. If no matching entry is found, the CLI SHALL print an error and exit with a non-zero code without executing any removal steps.

#### Scenario: Ref exists
- **WHEN** the named ref is present in `.gitmodules`
- **THEN** the command proceeds to the removal sequence

#### Scenario: Ref does not exist
- **WHEN** no ref with the given name is registered
- **THEN** the CLI prints an error indicating the ref was not found and exits with a non-zero code

### Requirement: Removal executes three steps in order
After validation passes, the CLI SHALL perform full removal by executing the following three operations in order:
1. `git submodule deinit -f <name>` — deinitializes the submodule and clears the working tree checkout
2. `git rm -f <name>` — removes the submodule from the working tree and from `.gitmodules`
3. `rm -rf .git/modules/<name>` — removes the cached module data from the git object store

Each step SHALL complete before the next begins. Stdout and stderr from git commands SHALL be streamed to the terminal in real time.

#### Scenario: All three steps execute in sequence
- **WHEN** all validations pass
- **THEN** the CLI runs `git submodule deinit -f <name>`, then `git rm -f <name>`, then `rm -rf .git/modules/<name>`, in that order

#### Scenario: Output is streamed
- **WHEN** a git command produces output during removal
- **THEN** that output is written to the terminal as it is produced, not buffered until completion

### Requirement: Failure in any step stops execution
If any step exits with a non-zero code, the CLI SHALL immediately stop, print an error indicating which step failed, and exit with the failing step's exit code. Subsequent steps SHALL NOT be attempted.

#### Scenario: deinit step fails
- **WHEN** `git submodule deinit -f <name>` exits with a non-zero code
- **THEN** the CLI reports the failure, does not proceed to subsequent steps, and exits with the same non-zero code

#### Scenario: Successful removal exits with code 0
- **WHEN** all three steps complete successfully
- **THEN** the CLI exits with code 0

### Requirement: Ref entry is removed from .wip.yml after successful removal
After all three removal steps succeed, the CLI SHALL remove the ref's entry from `.wip.yml` if one exists. If no entry exists in `.wip.yml` for the ref, this step is a no-op.

#### Scenario: .wip.yml entry removed
- **WHEN** the ref has an entry in `.wip.yml` and removal succeeds
- **THEN** the entry is removed from `.wip.yml`

#### Scenario: No .wip.yml entry — no-op
- **WHEN** the ref has no entry in `.wip.yml` and removal succeeds
- **THEN** `.wip.yml` is unchanged and the command exits with code 0
