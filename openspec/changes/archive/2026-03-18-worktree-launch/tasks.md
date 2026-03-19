## 1. Config

- [x] 1.1 Add `OnWorktreeLaunch []string` field to `SubmoduleConfig` in `cmd/config.go` with yaml tag `on-worktree-launch`

## 2. submodule add

- [x] 2.1 Add `--on-worktree-launch` repeatable flag to `submoduleAdd` in `cmd/submodule_add.go`
- [x] 2.2 Write `on-worktree-launch` commands to `.wip.yml` when flag is provided (mirror the existing `on-worktree-create` logic)
- [x] 2.3 Update usage string to document `--on-worktree-launch`

## 3. worktree launch command

- [x] 3.1 Create `cmd/worktree_launch.go` with `worktreeLaunch(args []string)` function
- [x] 3.2 Parse positional args `<submodule>` and `<worktree>`; print usage and exit 1 if either is missing
- [x] 3.3 Load and require `.wip.yml`; load and verify git repo
- [x] 3.4 Resolve worktree path using existing `repoRoot()` and `worktreePathSegment()` helpers
- [x] 3.5 Verify worktree directory exists with `os.Stat`; print error and exit 1 if not found
- [x] 3.6 Check `on-worktree-launch` hooks for submodule; print no-hooks message and exit 0 if empty
- [x] 3.7 Execute hooks sequentially with cwd = worktree dir; print ✓/✗ per hook; continue on failure; exit 0

## 4. Router

- [x] 4.1 Add `"launch"` case to `Worktree()` switch in `cmd/worktree.go` routing to `worktreeLaunch`
- [x] 4.2 Add `launch` to `worktreeUsage()` help text

## 5. Tests

- [x] 5.1 Add tests for `worktreeLaunch`: missing args, worktree not found, no hooks configured, hooks execute with correct cwd
- [x] 5.2 Add tests for `submoduleAdd` `--on-worktree-launch` flag: single command, multiple commands preserve order, omitted flag writes nothing
- [x] 5.3 Run `go test ./...` and `go vet ./...` to verify all pass
