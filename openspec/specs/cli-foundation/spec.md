## ADDED Requirements

### Requirement: CLI entry point routes to subcommands
The CLI SHALL read `os.Args[1]` to determine the top-level subcommand and delegate remaining arguments to the appropriate handler. An unknown or missing subcommand SHALL print usage and exit with a non-zero code.

#### Scenario: Known subcommand is dispatched — ref
- **WHEN** the user runs `wip ref <args>`
- **THEN** the ref handler is invoked with the remaining args

#### Scenario: Known subcommand is dispatched — worktree
- **WHEN** the user runs `wip worktree <args>`
- **THEN** the worktree handler is invoked with the remaining args

#### Scenario: Known subcommand is dispatched — init
- **WHEN** the user runs `wip init`
- **THEN** the init handler is invoked

#### Scenario: Unknown subcommand
- **WHEN** the user runs `wip unknowncmd`
- **THEN** the CLI prints usage information and exits with a non-zero code

#### Scenario: No subcommand provided
- **WHEN** the user runs `wip` with no arguments
- **THEN** the CLI prints usage information and exits with a non-zero code

### Requirement: Subcommands own their flag sets
Each subcommand SHALL define its own `flag.FlagSet` and parse only its portion of `os.Args`. Flags from one subcommand SHALL NOT bleed into another.

#### Scenario: Subcommand flag isolation
- **WHEN** a flag is defined on `ref add`
- **THEN** that flag is not recognized at the top-level or by other subcommands
