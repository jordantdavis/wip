## Context

`wip` already has an `on-worktree-create` hook that runs once when a worktree is first created. The execution model (sequential blocking commands, cwd = worktree dir, ✓/✗ output per command) is established and working. This change adds a parallel `on-worktree-launch` hook type and a new `wip worktree launch` subcommand to trigger it.

The primary use case is re-entering a worktree in a ready state: syncing to latest, reinstalling deps, and optionally starting an interactive environment (editor, tmux session, dev server).

## Goals / Non-Goals

**Goals:**
- Add `wip worktree launch <submodule> <worktree>` subcommand
- Add `--on-worktree-launch` flag to `wip submodule add`
- Extend `SubmoduleConfig` and `.wip.yml` schema with `on-worktree-launch`
- Verify worktree directory exists before running any hooks
- Print informational message when no hooks are configured

**Non-Goals:**
- Background or detached process management — users control this via the commands they configure
- Modifying `wip worktree add` — launch is a separate, explicit user invocation
- Shell expansion or templating in hook commands

## Decisions

### Execution model: same as on-worktree-create

`on-worktree-launch` hooks run sequentially and blocking, with cwd set to the worktree directory. This is identical to `on-worktree-create`. The final command naturally takes over the terminal if it is interactive (e.g. `claude`, `tmux new-session`).

Alternatives considered:
- **Background all commands**: Adds process tracking complexity with no clear benefit; users who want backgrounded processes can write `myserver &` or use tmux themselves.
- **Special "last command" exec**: Replacing the process via `syscall.Exec` for the last hook would give cleaner terminal handoff but adds OS-specific complexity. Not worth it — `exec.Command` with streaming I/O works correctly for interactive commands.

### Failure behavior: warn and continue (same as on-worktree-create)

If a hook exits non-zero, print a warning and continue running remaining hooks. The overall exit code is 0. This mirrors `on-worktree-create` and keeps the behavior consistent. Launch hooks are expected to be idempotent; a failure in one (e.g. `git pull` with no network) shouldn't prevent subsequent hooks from running.

### Worktree path resolution: same formula as worktree add/remove

The worktree path is `<repo root>/worktrees/<submodule>/<worktree-path-segment>` where `worktreePathSegment` replaces `/` with `-`. This is already used in `worktree_add.go` and `worktree_remove.go`.

### No-hooks message exits 0

When no `on-worktree-launch` hooks are configured, print an informational message to stdout and exit 0. This is not an error — the user may be exploring, or hooks may not yet be configured for that submodule.

### New file: cmd/worktree_launch.go

Consistent with the existing pattern (`worktree_add.go`, `worktree_remove.go`). The router in `worktree.go` gains a `"launch"` case.

## Risks / Trade-offs

- **Idempotency is user responsibility** → Document clearly in help text that launch hooks run every time; commands like `cp .env.example .env` will clobber edits. Mitigation: help text notes this explicitly.
- **Worktree existence check uses `os.Stat`** → If the worktree directory was created outside `wip`, it will still be found. This is intentional and acceptable.
