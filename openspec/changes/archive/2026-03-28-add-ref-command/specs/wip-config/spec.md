## MODIFIED Requirements

### Requirement: .wip.yml is scaffolded at repo root by wip init
`wip init` SHALL create a `.wip.yml` file at the repository root if one does not already exist. If `.wip.yml` already exists, `wip init` SHALL leave it unchanged.

#### Scenario: .wip.yml does not exist
- **WHEN** the user runs `wip init` and no `.wip.yml` exists at the repo root
- **THEN** a `.wip.yml` file is created with an empty `refs` map

#### Scenario: .wip.yml already exists
- **WHEN** the user runs `wip init` and `.wip.yml` already exists at the repo root
- **THEN** the file is left unchanged and the command exits with code 0

### Requirement: .wip.yml schema
`.wip.yml` SHALL use the following structure:

```yaml
refs:
  <name>:
    url: <url>
    branch: <branch>
    on-worktree-create:
      - <command>
    on-worktree-launch:
      - <command>
```

All keys SHALL be kebab-case. The `refs` map MAY be empty. The `url` field stores the git remote URL for the ref. The `branch` field is optional per ref entry and defaults to `main` when absent. The `on-worktree-create` and `on-worktree-launch` fields are both optional per ref entry.

#### Scenario: Empty config is valid
- **WHEN** `.wip.yml` contains only an empty `refs` map
- **THEN** the file is valid and all commands proceed without error

#### Scenario: Config with url and branch fields is valid
- **WHEN** `.wip.yml` contains a ref entry with `url` and `branch` fields
- **THEN** the file is parsed correctly and both values are accessible

#### Scenario: Config without branch field defaults to main
- **WHEN** `.wip.yml` contains a ref entry without a `branch` field
- **THEN** the branch is treated as `main`

#### Scenario: Config with hook lists is valid
- **WHEN** `.wip.yml` contains a ref entry with `on-worktree-create` and/or `on-worktree-launch` lists
- **THEN** both lists are parsed correctly and independently accessible

### Requirement: .wip.yml is required for wip ref add and wip worktree add
Both `wip ref add` and `wip worktree add` SHALL check for the presence of `.wip.yml` before proceeding. If no `.wip.yml` is found within the user's home directory tree, the CLI SHALL print an error directing the user to run `wip init` and exit with a non-zero code.

#### Scenario: .wip.yml present
- **WHEN** the user runs `wip ref add` or `wip worktree add` and `.wip.yml` exists at or above the current directory within the home tree
- **THEN** the command proceeds past the config check

#### Scenario: .wip.yml absent
- **WHEN** the user runs `wip ref add` or `wip worktree add` and no `.wip.yml` is found
- **THEN** the CLI prints an error suggesting `wip init` and exits with a non-zero code

## REMOVED Requirements

### Requirement: Legacy submodules schema key
**Reason**: The `submodules` top-level key is replaced by `refs` to align with the `wip ref` command family.
**Migration**: Rename `submodules:` to `refs:` in `.wip.yml`. All nested structure is otherwise unchanged.
