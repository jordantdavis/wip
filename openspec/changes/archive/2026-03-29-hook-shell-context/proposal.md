## Why

Hook commands in `.wip.yml` currently run via direct `exec.Command`, giving hooks no access to contextual values like the ref name or worktree path. Users have no way to write hooks that adapt to their execution context without hardcoding values.

## What Changes

- Hook execution switches from direct `exec.Command` to `sh -c "<hook>"`, making hooks full shell commands
- Before each hook runs, a set of `WIP_*` environment variables is injected as ambient context
- The same env vars are available in both `on-worktree-create` and `on-worktree-launch` hooks
- The hook output label (✓/✗) shows the hook string as written in config

## Capabilities

### New Capabilities

- `hook-shell-execution`: Hooks run via `sh -c`, enabling shell features (pipes, redirects, `&&`, variable expansion)
- `hook-context-env-vars`: A standard set of `WIP_*` env vars is injected into every hook execution, providing ref name, worktree name, worktree path, and repo root as ambient context

### Modified Capabilities

- `worktree-add`: Hook execution behavior changes (direct exec → `sh -c` + env vars)
- `worktree-launch`: Hook execution behavior changes (direct exec → `sh -c` + env vars)

## Impact

- `cmd/worktree_add.go`: Hook execution loop updated
- `cmd/worktree_launch.go`: Hook execution loop updated
- `cmd/worktree_test.go`, `cmd/worktree_launch_test.go`: Tests updated to cover env var injection and shell features
- `.wip.yml` configs using hooks: no breaking change — existing hooks work unchanged under `sh -c`
