## ADDED Requirements

### Requirement: version command appears in help output
The CLI help output SHALL list `version` as an available command when usage is printed.

#### Scenario: version listed in help
- **WHEN** a user runs `wip` with no arguments or an unknown command
- **THEN** the help output includes a `version` entry with a short description

#### Scenario: version entry matches style of other commands
- **WHEN** the help output is printed
- **THEN** the `version` line uses the same indentation and lowercase verb-first description format as `init`, `submodule`, and `worktree`
