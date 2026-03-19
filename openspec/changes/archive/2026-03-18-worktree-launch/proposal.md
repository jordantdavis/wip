## Why

When returning to an existing worktree, developers must manually run setup steps (pull latest, reinstall deps) and launch their environment (dev server, editor, tmux session). There is no automated way to re-enter a worktree in a ready state.

## What Changes

- `wip submodule add` gains a repeatable `--on-worktree-launch` flag that stores ordered commands in `.wip.yml` under `submodules.<name>.on-worktree-launch`
- New `wip worktree launch <submodule> <worktree>` subcommand that verifies the worktree exists, then runs its configured `on-worktree-launch` hooks sequentially in the worktree directory
- If no hooks are configured for the submodule, the command prints an informational message and exits cleanly
- Hook commands run sequentially and blocking; the last command may be interactive and take over the terminal (e.g. `claude`, `tmux new-session`)

## Capabilities

### New Capabilities

- `worktree-launch`: The `wip worktree launch` subcommand — signature, worktree path resolution, existence check, no-hooks message, and hook execution contract

### Modified Capabilities

- `submodule-add`: Gains the repeatable `--on-worktree-launch` flag, writing hooks to `.wip.yml`
- `worktree-lifecycle`: New `on-worktree-launch` hook type — configuration shape, execution contract (cwd, ordering, idempotency expectation, failure behavior), and how it differs from `on-worktree-create`
- `wip-config`: `SubmoduleConfig` gains the `on-worktree-launch` field

## Impact

- New file: `cmd/worktree_launch.go`
- Modified files: `cmd/submodule_add.go`, `cmd/worktree.go`, `cmd/config.go`
- No new dependencies
