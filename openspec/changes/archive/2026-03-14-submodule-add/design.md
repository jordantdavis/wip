## Context

Fresh Go project (`main.go` is a stub). The `wip` CLI is being built from scratch as a personal AI coding workflow tool. This change establishes the CLI skeleton and delivers the first real subcommand. The `flags` package is mandated. No external CLI frameworks.

## Goals / Non-Goals

**Goals:**
- Establish a two-level subcommand routing pattern that scales to future subcommands
- Implement `wip submodule add <url> [<path>]` with full validation

**Non-Goals:**
- Any submodule subcommands beyond `add`
- Help text beyond what `flag.Usage` provides
- Shell completion

## Decisions

### 1. Built-in `flags` package over a third-party CLI framework
Using `flag.NewFlagSet` per subcommand. Each level of routing is a manual switch on `os.Args`. This keeps zero external deps for CLI parsing and fits the project's preference.

Alternatives considered: `cobra`, `urfave/cli` — both add deps and abstractions that aren't needed for a small personal tool.

### 2. URL as mandatory positional argument, name as optional named flag
URL is read from `flagSet.Args()[0]` after flag parsing. `--name` is optional and defaults to empty (git derives the name and checkout directory from the URL).

When `--name` is provided it serves dual purpose: passed as `--name <name>` to set git's internal submodule name, and as the path positional to set the checkout directory at repo root. This lets the same remote URL be added more than once under distinct names, avoiding the git conflict on duplicate submodule names.

### 3. Two-level routing: main → submodule → add
`main.go` switches on `os.Args[1]` and delegates `os.Args[2:]` to `cmd/submodule.go`. `submodule.go` switches on the first remaining arg to reach `add`. Each subcommand owns its own `flag.NewFlagSet`.

### 4. Name validation: no path separators
A name like `libs/foo` would cause git to create a nested directory, conflicting with the "initialized at root" constraint. Reject any name containing `/` or `\` before invoking git.

## Risks / Trade-offs

- **Subcommand routing boilerplate** — manual switch statements will repeat as commands grow. Acceptable for now; extract a router if/when it becomes unwieldy. → No mitigation needed yet.
- **Duplicate default names** — without `--name`, git derives the submodule name from the URL. Adding the same URL twice without `--name` will fail; the user must supply `--name` to disambiguate. This is expected behavior, not a bug.
- **git must be installed** — the command shells out to `git`. No check for git binary presence beyond the working directory being a git repo. → Document as a prerequisite.
