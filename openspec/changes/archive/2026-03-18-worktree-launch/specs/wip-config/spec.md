## MODIFIED Requirements

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
