## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any submodule operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip submodule list` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip submodule list` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Submodule data is read from git's native storage
The CLI SHALL read submodule information from git's native `.gitmodules` file at the root of the repository. The data source SHALL NOT be a custom persistence layer. For each submodule, the CLI SHALL retrieve at minimum its name and URL. If `.gitmodules` is absent or contains no submodule entries, the CLI SHALL treat the submodule list as empty.

#### Scenario: .gitmodules exists with entries
- **WHEN** the repository contains a `.gitmodules` file with one or more submodule entries
- **THEN** the CLI reads each entry's name and URL from git's native storage and uses those values for output

#### Scenario: .gitmodules absent
- **WHEN** no `.gitmodules` file exists in the repository root
- **THEN** the CLI treats the list as empty and proceeds to the empty-state output

#### Scenario: .gitmodules exists but has no submodule entries
- **WHEN** the `.gitmodules` file exists but contains no `[submodule]` sections
- **THEN** the CLI treats the list as empty and proceeds to the empty-state output

### Requirement: Output format is one submodule per line with name and URL
For each submodule found, the CLI SHALL print exactly one line to stdout in the format:

```
<name>  <url>
```

where `<name>` and `<url>` are separated by two spaces. No header line, no trailing blank lines, and no additional decorations SHALL be added.

#### Scenario: Single submodule
- **WHEN** exactly one submodule is registered
- **THEN** the CLI prints one line containing the submodule's name and URL separated by two spaces

#### Scenario: Multiple submodules
- **WHEN** two or more submodules are registered
- **THEN** the CLI prints one line per submodule, each containing that submodule's name and URL separated by two spaces

### Requirement: Empty state produces a descriptive message
When no submodules are registered in the repository, the CLI SHALL print a single informational message to stdout indicating that no submodules are present. The CLI SHALL exit with code 0.

#### Scenario: No submodules present
- **WHEN** the repository has no registered submodules
- **THEN** the CLI prints a message such as `no submodules found` and exits with code 0

### Requirement: Submodules are listed in alphabetical order by name
The CLI SHALL sort the output lines alphabetically by submodule name using case-sensitive lexicographic ordering. This ensures deterministic output regardless of the order entries appear in `.gitmodules`.

#### Scenario: Submodules listed in alphabetical order
- **WHEN** the repository contains submodules with names `zebra`, `alpha`, and `monkey`
- **THEN** the CLI prints them in the order `alpha`, `monkey`, `zebra`

#### Scenario: Single submodule order is trivially correct
- **WHEN** the repository contains exactly one submodule
- **THEN** the CLI prints that single submodule without any ordering concern
