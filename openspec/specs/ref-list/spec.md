## ADDED Requirements

### Requirement: Working directory must be a git repository
Before executing any ref operation, the CLI SHALL verify the current working directory contains a `.git` entry. If not, it SHALL print an error and exit with a non-zero code.

#### Scenario: Valid git repository
- **WHEN** the user runs `wip ref list` inside a git repository
- **THEN** the command proceeds past the git repo check

#### Scenario: Not a git repository
- **WHEN** the user runs `wip ref list` outside a git repository
- **THEN** the CLI prints an error indicating the directory is not a git repository and exits with a non-zero code

### Requirement: Ref data is read from .gitmodules
The CLI SHALL read ref information from `.gitmodules` at the repository root. For each ref, the CLI SHALL retrieve its name, URL, and configured branch. If `.gitmodules` is absent or contains no submodule entries, the CLI SHALL treat the ref list as empty.

#### Scenario: .gitmodules exists with entries
- **WHEN** the repository contains a `.gitmodules` file with one or more submodule entries
- **THEN** the CLI reads each entry's name, URL, and branch from `.gitmodules`

#### Scenario: .gitmodules absent
- **WHEN** no `.gitmodules` file exists in the repository root
- **THEN** the CLI treats the list as empty and proceeds to the empty-state output

### Requirement: Output format is one ref per line with name, branch, and URL
For each ref found, the CLI SHALL print exactly one line to stdout in the format:

```
<name>  <branch>  <url>
```

where fields are separated by two spaces. No header line, no trailing blank lines, and no additional decorations SHALL be added.

#### Scenario: Single ref
- **WHEN** exactly one ref is registered
- **THEN** the CLI prints one line containing the ref's name, branch, and URL separated by two spaces

#### Scenario: Multiple refs
- **WHEN** two or more refs are registered
- **THEN** the CLI prints one line per ref, each containing that ref's name, branch, and URL

### Requirement: Empty state produces a descriptive message
When no refs are registered, the CLI SHALL print a single informational message to stdout indicating that no refs are present and exit with code 0.

#### Scenario: No refs present
- **WHEN** the repository has no registered refs
- **THEN** the CLI prints a message such as `no refs found` and exits with code 0

### Requirement: Refs are listed in alphabetical order by name
The CLI SHALL sort output lines alphabetically by ref name using case-sensitive lexicographic ordering.

#### Scenario: Refs listed in alphabetical order
- **WHEN** the repository contains refs with names `zebra`, `alpha`, and `monkey`
- **THEN** the CLI prints them in the order `alpha`, `monkey`, `zebra`
