## MODIFIED Requirements

### Requirement: .wip.yml is required for wip submodule add and wip worktree add
Both `wip submodule add` and `wip worktree add` SHALL check for the presence of `.wip.yml` before proceeding. The check is performed as part of project discovery — `findWipProject()` locates and loads the nearest `.wip.yml` by walking up from the current directory. If no `.wip.yml` is found within the user's home directory tree, the CLI SHALL print an error directing the user to run `wip init` and exit with a non-zero code.

#### Scenario: .wip.yml present
- **WHEN** the user runs `wip submodule add` or `wip worktree add` and `.wip.yml` exists at or above the current directory within the home tree
- **THEN** the command proceeds past the config check using the discovered project root

#### Scenario: .wip.yml absent
- **WHEN** the user runs `wip submodule add` or `wip worktree add` and no `.wip.yml` exists at or above the current directory within the home tree
- **THEN** the CLI prints an error indicating `.wip.yml` is missing and suggests running `wip init`, then exits with a non-zero code
