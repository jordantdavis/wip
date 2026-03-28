### Requirement: wip root prints the project root path
`wip root` SHALL print the absolute path of the nearest `.wip.yml` directory discovered by walking upward from the current working directory, followed by a newline, then exit with code 0.

#### Scenario: Invoked from the project root
- **WHEN** the user runs `wip root` from the directory containing `.wip.yml`
- **THEN** that directory's absolute path is printed to stdout

#### Scenario: Invoked from a subdirectory
- **WHEN** the user runs `wip root` from a subdirectory of the project
- **THEN** the absolute path of the ancestor directory containing `.wip.yml` is printed to stdout

#### Scenario: No project found
- **WHEN** the user runs `wip root` and no `.wip.yml` is found within the home directory tree
- **THEN** an error message is printed to stderr and the command exits with a non-zero code

### Requirement: wip root uses the same discovery logic as other subcommands
`wip root` SHALL use the same `findWipProject` function as `wip submodule` and `wip worktree`. It is subject to the same home directory ceiling and outside-home failure behavior.

#### Scenario: Invoked from outside home directory
- **WHEN** the user runs `wip root` from a directory outside their home directory
- **THEN** the command exits immediately with a non-zero code and an error message

### Requirement: wip root is documented in the README
The README SHALL include a section describing `wip root` and a note explaining that all subcommands work from any subdirectory of a wip project.

#### Scenario: README contains wip root section
- **WHEN** a user reads the README
- **THEN** they can find documentation for `wip root` with a usage example showing it printing the project root path
