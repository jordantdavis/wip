## ADDED Requirements

### Requirement: .wip.yml is scaffolded at repo root by wip init
`wip init` SHALL create a `.wip.yml` file at the repository root if one does not already exist. If `.wip.yml` already exists, `wip init` SHALL leave it unchanged.

#### Scenario: .wip.yml does not exist
- **WHEN** the user runs `wip init` and no `.wip.yml` exists at the repo root
- **THEN** a `.wip.yml` file is created with an empty `submodules` map

#### Scenario: .wip.yml already exists
- **WHEN** the user runs `wip init` and `.wip.yml` already exists at the repo root
- **THEN** the file is left unchanged and the command exits with code 0

### Requirement: .wip.yml schema
`.wip.yml` SHALL use the following structure:

```yaml
submodules:
  <name>:
    on-worktree-create:
      - <command>
    on-worktree-launch:
      - <command>
```

All keys SHALL be kebab-case. The `submodules` map MAY be empty. The `on-worktree-create` and `on-worktree-launch` fields are both optional per submodule entry. Each command in either list is a string.

#### Scenario: Empty config is valid
- **WHEN** `.wip.yml` contains only an empty `submodules` map
- **THEN** the file is valid and all commands proceed without error

#### Scenario: Config with on-worktree-create list is valid
- **WHEN** `.wip.yml` contains a submodule entry with an `on-worktree-create` list of strings
- **THEN** the file is parsed correctly and the commands are accessible

#### Scenario: Config with on-worktree-launch list is valid
- **WHEN** `.wip.yml` contains a submodule entry with an `on-worktree-launch` list of strings
- **THEN** the file is parsed correctly and the commands are accessible

#### Scenario: Config with both hook lists is valid
- **WHEN** `.wip.yml` contains a submodule entry with both `on-worktree-create` and `on-worktree-launch` lists
- **THEN** both lists are parsed correctly and independently accessible

### Requirement: .wip.yml is required for wip submodule add and wip worktree add
Both `wip submodule add` and `wip worktree add` SHALL check for the presence of `.wip.yml` at the repo root before proceeding. If absent, the CLI SHALL print an error directing the user to run `wip init` and exit with a non-zero code.

#### Scenario: .wip.yml present
- **WHEN** the user runs `wip submodule add` or `wip worktree add` and `.wip.yml` exists
- **THEN** the command proceeds past the config check

#### Scenario: .wip.yml absent
- **WHEN** the user runs `wip submodule add` or `wip worktree add` and `.wip.yml` does not exist
- **THEN** the CLI prints an error indicating `.wip.yml` is missing and suggests running `wip init`, then exits with a non-zero code
