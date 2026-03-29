## 1. Config Layer

- [x] 1.1 Rename `WipConfig.Submodules` to `WipConfig.Refs` and `SubmoduleConfig` to `RefConfig` in `cmd/config.go`
- [x] 1.2 Add `Branch string` field to `RefConfig` with YAML key `branch`; default to `"main"` when absent
- [x] 1.3 Update `wip init` to scaffold `refs:` instead of `submodules:` in the empty `.wip.yml`

## 2. Ref Add

- [x] 2.1 Create `cmd/ref_add.go` with `refAdd(args []string)` function
- [x] 2.2 Add `--name`, `--branch` (default `main`), `--on-worktree-create`, and `--on-worktree-launch` flags
- [x] 2.3 Validate git repo, URL format, and name (no path separators)
- [x] 2.4 Execute `git submodule add -b <branch> <url> [<name>]` as subprocess
- [x] 2.5 Write `ignore = all` to the new submodule's entry in `.gitmodules` after successful add
- [x] 2.6 Write hook config to `.wip.yml` under `refs.<name>` when flags are provided

## 3. Ref List

- [x] 3.1 Create `cmd/ref_list.go` with `refList(args []string)` function
- [x] 3.2 Read all submodule entries from `.gitmodules` (name, URL, branch)
- [x] 3.3 Print one line per ref in `<name>  <branch>  <url>` format, sorted alphabetically by name
- [x] 3.4 Print empty-state message when no refs are registered

## 4. Ref Remove

- [x] 4.1 Create `cmd/ref_remove.go` with `refRemove(args []string)` function
- [x] 4.2 Validate git repo and that the named ref exists in `.gitmodules`
- [x] 4.3 Execute the three-step removal: `git submodule deinit -f`, `git rm -f`, `rm -rf .git/modules/<name>` in order, stopping on first failure
- [x] 4.4 Remove the ref's entry from `.wip.yml` after successful removal

## 5. Ref Sync

- [x] 5.1 Create `cmd/ref_sync.go` with `refSync(args []string)` function
- [x] 5.2 Add optional `--name` flag; when absent, discover all refs from `.gitmodules`
- [x] 5.3 Execute `git submodule update --remote <name>` for each ref; run all concurrently when no `--name` is given
- [x] 5.4 Buffer output and report `✓ <name>` / `✗ <name>: <error>` per ref after all complete
- [x] 5.5 Exit non-zero if any update failed

## 6. Ref Restore

- [x] 6.1 Create `cmd/ref_restore.go` with `refRestore(args []string)` function
- [x] 6.2 Discover all refs from `.gitmodules`; print empty-state message and exit 0 if none
- [x] 6.3 Execute `git submodule update --init --remote` for all refs concurrently
- [x] 6.4 Buffer output and report `✓ <name>` / `✗ <name>: <error>` per ref after all complete
- [x] 6.5 Exit non-zero if any initialization failed

## 7. Router and Entrypoint

- [x] 7.1 Create `cmd/ref.go` router that dispatches `add`, `list`, `remove`, `sync`, `restore` subcommands
- [x] 7.2 Update `main.go` to route `wip ref` to `cmd.Ref()` and remove the `wip submodule` route
- [x] 7.3 Update `printUsage()` in `main.go` to list `ref` and remove `submodule`

## 8. Cleanup

- [x] 8.1 Delete `cmd/submodule.go`, `cmd/submodule_add.go`, `cmd/submodule_list.go`, `cmd/submodule_remove.go`, `cmd/submodule_sync.go`
- [x] 8.2 Update any references to `SubmoduleConfig` or `WipConfig.Submodules` in worktree commands (`cmd/worktree_add.go`, `cmd/worktree_launch.go`) to use `RefConfig` and `WipConfig.Refs`

## 9. Tests

- [x] 9.1 Update config tests in `cmd/config_test.go` to use `refs` key and `RefConfig` type
- [x] 9.2 Add tests for `--branch` flag defaulting to `main` in ref add
- [x] 9.3 Verify `ignore = all` is written to `.gitmodules` after `wip ref add`
- [x] 9.4 Verify `wip ref restore` uses `--remote` flag (not just `--init`)
