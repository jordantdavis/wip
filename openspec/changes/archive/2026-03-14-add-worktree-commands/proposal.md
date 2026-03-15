## Why

The `wip` CLI manages git submodules but provides no way to create or manage git worktrees within those submodules. Worktrees are a standard git workflow for working on multiple branches simultaneously, and first-class support in `wip` would make it seamless to create, list, and remove worktrees scoped to individual submodules.

## What Changes

- Adds a new top-level `wip worktree` command with three subcommands: `add`, `list`, and `remove`
- `wip worktree add <submodule> <worktree>` creates a git worktree at `worktrees/<submodule>/<worktree>/` and a new branch named `<worktree>` inside the submodule repo; `--existing-branch` checks out an existing branch instead
- `wip worktree list` scans the filesystem under `worktrees/` and prints a flat two-column table of submodule and worktree names
- `wip worktree remove <submodule> <worktree>` removes the worktree at `worktrees/<submodule>/<worktree>/`; `--delete-branch` also deletes the branch (safe `-d`) from the submodule repo
- `worktrees/` directory is created automatically on first `add` if it does not exist

## Capabilities

### New Capabilities

- `worktree-add`: Create a git worktree for a submodule at a conventional path, optionally targeting an existing branch
- `worktree-list`: Discover and display all worktrees across all submodules by scanning the filesystem
- `worktree-remove`: Remove a git worktree for a submodule, with an option to also delete the associated branch

### Modified Capabilities

- `cli-foundation`: A new top-level subcommand `worktree` must be registered in the CLI entry point routing

## Impact

- `main.go`: Add `worktree` case to the top-level command switch
- `cmd/worktree.go`: New dispatcher for the `worktree` subcommands
- `cmd/worktree_add.go`: New file implementing `add`
- `cmd/worktree_list.go`: New file implementing `list`
- `cmd/worktree_remove.go`: New file implementing `remove`
- No changes to existing submodule commands or shared utilities
- No new external dependencies; relies on `git worktree` which is available in git 2.5+
