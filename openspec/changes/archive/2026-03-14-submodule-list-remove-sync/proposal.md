## Why

The `wip submodule add` command lets users register submodules by name, but there is no way to inspect what's registered, remove one, or pull the latest changes. These three operations complete the core submodule lifecycle.

## What Changes

- Add `wip submodule list` — reads `.gitmodules` and prints each submodule's name and URL, one per line, sorted alphabetically
- Add `wip submodule remove <name>` — fully removes a submodule by name via the standard three-step git sequence (deinit, rm, purge modules cache)
- Add `wip submodule sync [--name <name>]` — updates all submodules concurrently (or a single named one), reports buffered pass/fail results after all complete

## Capabilities

### New Capabilities

- `submodule-list`: List registered submodules with name and URL from git's native `.gitmodules` storage
- `submodule-remove`: Fully remove a submodule by name using the git deinit + rm + modules cache purge sequence
- `submodule-sync`: Concurrently update all submodules (or a single named one) and report per-submodule results

### Modified Capabilities

## Impact

- `cmd/submodule.go`: Three new subcommand functions added; routing switch extended with `list`, `remove`, `sync` cases
- No new dependencies — uses only standard library (`flag`, `os/exec`, `sync`)
- `.gitmodules` parsing added (read via `git config --file .gitmodules` subprocess or direct file read)
