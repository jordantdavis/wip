## 1. CLI Foundation

- [x] 1.1 Add `worktree` case to the top-level command switch in `main.go` that calls `cmd.Worktree(os.Args[2:])`
- [x] 1.2 Update `printUsage()` in `main.go` to include the `worktree` command

## 2. Worktree Dispatcher

- [x] 2.1 Create `cmd/worktree.go` with a `Worktree(args []string)` dispatcher routing to `add`, `list`, and `remove` subcommands
- [x] 2.2 Add `worktreeUsage()` in `cmd/worktree.go` listing the three subcommands

## 3. Shared Utilities

- [x] 3.1 Add `validateWorktreeName(name string) error` to validate names match `[a-zA-Z0-9_-]+`
- [x] 3.2 Add `repoRoot() (string, error)` (or use `os.Getwd()` inline) for absolute path construction in worktree commands

## 4. worktree add

- [x] 4.1 Create `cmd/worktree_add.go` with `worktreeAdd(args []string)`
- [x] 4.2 Define `flag.FlagSet` with `--existing-branch` boolean flag
- [x] 4.3 Validate: check git repo, parse and validate submodule name (exists in `.gitmodules`), parse and validate worktree name character set
- [x] 4.4 Construct absolute worktree path: `<repo root>/worktrees/<submodule>/<worktree>`
- [x] 4.5 Create `worktrees/<submodule>/` via `os.MkdirAll` if it does not exist
- [x] 4.6 Build `exec.Cmd` for `git worktree add -b <worktree> <abs-path>` (default) or `git worktree add <abs-path> <worktree>` (`--existing-branch`), with `cmd.Dir` set to `<repo root>/<submodule>/`
- [x] 4.7 Stream stdout/stderr to terminal and exit with git's exit code

## 5. worktree list

- [x] 5.1 Create `cmd/worktree_list.go` with `worktreeList(args []string)`
- [x] 5.2 Check git repo, then check if `worktrees/` directory exists; print "no worktrees found" and exit 0 if absent or empty
- [x] 5.3 Enumerate `worktrees/<submodule>/*/` entries using `os.ReadDir` at two levels
- [x] 5.4 Collect all `(submodule, worktree)` pairs, sort by submodule then worktree name
- [x] 5.5 Print rows as `<submodule>  <worktree>` (two-space separator)

## 6. worktree remove

- [x] 6.1 Create `cmd/worktree_remove.go` with `worktreeRemove(args []string)`
- [x] 6.2 Define `flag.FlagSet` with `--delete-branch` boolean flag
- [x] 6.3 Validate: check git repo, parse and validate submodule name (exists in `.gitmodules`), parse and validate worktree name character set
- [x] 6.4 Verify `worktrees/<submodule>/<worktree>/` exists on the filesystem; error if not
- [x] 6.5 Build `exec.Cmd` for `git worktree remove <abs-path>` with `cmd.Dir` set to `<repo root>/<submodule>/`; stream output and exit on failure
- [x] 6.6 If `--delete-branch` is set and worktree removal succeeded, run `git branch -d <worktree>` inside the submodule directory; stream output and exit with git's exit code
