## Why

`wip` currently requires an existing git repository — there's no way to bootstrap a new wip project from scratch. An `init` command establishes the entry point for new projects and provides a foundation for future project-level setup steps.

## What Changes

- Add `wip init` as a new top-level command
- `wip init` runs `git init` in the current directory if no git repo is present
- If the current directory is already the root of a git repo, `wip init` is a no-op (idempotent)
- If the current directory is inside an existing git repo but not at its root, `wip init` exits with an error
- Uses `git rev-parse --git-dir` to detect repo state: exits non-zero means no repo, `.git` means at root, anything else means not at root
- Registers `init` in `main.go` dispatch and usage

## Capabilities

### New Capabilities

- `wip-init`: Initialize a wip project in the current directory, starting with git repo initialization

### Modified Capabilities

- `cli-foundation`: Register the new `init` command in the top-level dispatcher and usage output

## Impact

- New file: `cmd/init.go`
- Modified: `main.go` (add `case "init"` and update usage)
- No breaking changes
