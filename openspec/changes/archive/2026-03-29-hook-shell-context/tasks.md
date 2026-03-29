## 1. Core Hook Execution

- [x] 1.1 Extract a shared `runHooks` helper in `cmd/` that accepts ref name, worktree name, worktree path, repo root, and hook list — builds the `WIP_*` env var map, executes each hook via `sh -c`, and prints ✓/✗
- [x] 1.2 Replace the inline hook loop in `worktree_add.go` with a call to `runHooks`
- [x] 1.3 Replace the inline hook loop in `worktree_launch.go` with a call to `runHooks`

## 2. Tests

- [x] 2.1 Add test: hook can read `WIP_REF_NAME` via `sh -c` and the value matches the ref argument
- [x] 2.2 Add test: hook can read `WIP_WORKTREE_NAME` via `sh -c` and the value matches the worktree argument
- [x] 2.3 Add test: hook can read `WIP_WORKTREE_PATH` and the value is the absolute worktree path
- [x] 2.4 Add test: hook can read `WIP_ROOT` and the value is the repo root
- [x] 2.5 Add test: compound command (`echo a && echo b`) executes both parts
- [x] 2.6 Verify existing tests in `worktree_launch_test.go` and `worktree_test.go` still pass with the `sh -c` change
