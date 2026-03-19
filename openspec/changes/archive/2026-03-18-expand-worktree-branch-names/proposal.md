## Why

The current worktree name validator (`^[a-zA-Z0-9_-]+$`) rejects branch names that contain `/`, preventing common conventions like `feature/my-thing` or `team/user/ticket-123`. Worktree names should support any branch name that git itself considers valid.

## What Changes

- **Expand branch name validation**: Replace the hardcoded regex with a call to `git check-ref-format --branch <name>`, delegating validation to git.
- **Decouple branch name from path segment**: When a branch name contains `/`, replace each `/` with `-` to construct the worktree directory path (e.g., `feature/my-thing` → `worktrees/<sub>/feature-my-thing`).
- **Branch name is the canonical input everywhere**: `worktree add` and `worktree remove` both accept the branch name (with `/`); path derivation is an internal detail.
- **Expand `worktree list` output**: Add a third column showing the branch name checked out in each worktree, alongside the existing submodule and path columns.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `worktree-add`: Validation rule changes from regex to `git check-ref-format`; path construction now replaces `/` with `-`; git commands use the branch name directly.
- `worktree-list`: Output format changes from two columns (`submodule  path`) to three columns (`submodule  path  branch`); branch is read from each worktree via `git branch --show-current`.
- `worktree-remove`: Path derivation now replaces `/` with `-` to locate the worktree directory; `--delete-branch` uses the original branch name.

## Impact

- `cmd/worktree.go`: `validateWorktreeName` replaced by `validateBranchName` using `git check-ref-format`.
- `cmd/worktree_add.go`: Path construction uses `/`→`-` replacement.
- `cmd/worktree_list.go`: Output gains a branch column; reads branch via `git -C <path> branch --show-current`.
- `cmd/worktree_remove.go`: Path derivation uses `/`→`-` replacement; `--delete-branch` uses the raw branch name argument.
- Existing tests that assert on the old regex validation error message will need updating.
