## Why

The `wip` CLI needs a foundation and its first useful subcommand: `wip submodule add`, which initializes a git submodule in the working directory. This gives the AI coding workflow a consistent, validated interface for managing submodules rather than running raw git commands.

## What Changes

- Introduce the `wip` CLI entry point with two-level subcommand routing (`main.go` → `cmd/`)
- Add the `submodule` top-level command
- Add the `submodule add <url> [<path>]` subcommand that executes `git submodule add`
- Implement input validation: git repo detection, URL format, path constraints

## Capabilities

### New Capabilities

- `cli-foundation`: Entry point, subcommand routing, and project structure for the `wip` CLI
- `submodule-add`: The `wip submodule add <url> [<path>]` command — validation and git execution

### Modified Capabilities

## Impact

- `main.go`: Replaced with real CLI entry point and subcommand router
- `cmd/submodule.go`: New file implementing the submodule command and `add` subcommand
