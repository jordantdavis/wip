## Why

The `wip submodule` command uses raw git submodule machinery — SHA pinning, detached HEAD state, manual sync — for a narrow purpose: co-locating repos so AI tools like Claude Code and OpenCode can `@` reference their files in a session. The complexity of submodules far exceeds this use case, and the command name obscures the intent.

## What Changes

- **BREAKING**: `wip submodule` command is removed; replaced by `wip ref`
- **BREAKING**: `.wip.yml` config key `submodules` renamed to `refs`
- New `wip ref` command family manages git submodules under the hood with branch tracking (`branch = main` by default) and `ignore = all`, so git status stays clean and repos are always at latest HEAD
- `wip ref restore` replaces `git submodule update --init` for teammates setting up a cloned wip project; uses `--remote` to pull latest rather than the committed SHA
- `wip ref sync` pulls all refs to latest branch HEAD in one command
- The `on-worktree-create` and `on-worktree-launch` hook config moves from `submodules` to `refs` in `.wip.yml`

## Capabilities

### New Capabilities

- `ref-add`: Add a reference repo as a git submodule configured for branch tracking and clean status
- `ref-list`: List configured reference repos and their sync status
- `ref-remove`: Remove a reference repo and clean up its submodule registration
- `ref-sync`: Pull all reference repos to latest HEAD on their tracked branch
- `ref-restore`: Initialize and populate reference repos after cloning a wip project

### Modified Capabilities

- `wip-config`: Config schema changes — `submodules` key renamed to `refs`, submodule config entries gain optional `branch` field (defaults to `main`)

## Impact

- `cmd/submodule.go`, `cmd/submodule_add.go`, `cmd/submodule_list.go`, `cmd/submodule_remove.go`, `cmd/submodule_sync.go` — removed
- New `cmd/ref.go`, `cmd/ref_add.go`, `cmd/ref_list.go`, `cmd/ref_remove.go`, `cmd/ref_sync.go`, `cmd/ref_restore.go`
- `cmd/config.go` — `WipConfig.Submodules` field renamed to `Refs`; `SubmoduleConfig` renamed to `RefConfig` with added `Branch` field
- `main.go` — routes `ref` command, removes `submodule` route
- `.wip.yml` files in existing projects require migration (rename `submodules:` to `refs:`)
