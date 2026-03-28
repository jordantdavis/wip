## 1. Core: WipProject type and findWipProject

- [x] 1.1 Add `WipProject` struct to `cmd/config.go` with `Root string` and `Config *WipConfig` fields
- [x] 1.2 Implement `findWipProject()` in `cmd/config.go`: get cwd, check if within home directory (fail immediately if not), walk upward checking for `.wip.yml` at each level, stop after checking home, return `*WipProject` or error
- [x] 1.3 Add unit tests for `findWipProject` in `cmd/config_test.go`: found in cwd, found in parent, not found within home, invoked from outside home

## 2. Main dispatch: chdir before subcommand

- [x] 2.1 In `main.go`, for the `submodule` and `worktree` cases, call `cmd.FindWipProject()` (exported) and `os.Chdir(project.Root)` before dispatching to `cmd.Submodule` or `cmd.Worktree`
- [x] 2.2 Ensure `init` and `version` cases are left unchanged (no project discovery)
- [x] 2.3 Print a clear error and exit non-zero if `FindWipProject()` fails

## 3. wip root command

- [x] 3.1 Create `cmd/root.go` implementing `Root(args []string)`: call `FindWipProject()`, print `project.Root` to stdout, exit 0; on error print to stderr and exit non-zero
- [x] 3.2 Add `root` case to the dispatch switch in `main.go` (calls `cmd.Root`, does NOT chdir — just finds and prints)
- [x] 3.3 Add `root` to the usage list in `main.go`'s `printUsage()`

## 4. README documentation

- [x] 4.1 Add a note under the Commands section header explaining that all `wip submodule` and `wip worktree` commands work from any subdirectory of a wip project
- [x] 4.2 Add a `### wip root` section to the README with usage example showing it printing the project root path
- [x] 4.3 Update the top-level usage listing in `main.go`'s `printUsage()` if it doesn't already mention `root`

## 5. Verification

- [x] 5.1 Run `go build ./...` — confirm no compilation errors
- [x] 5.2 Run `go vet ./...` — confirm no vet issues
- [x] 5.3 Run `go test ./...` — confirm all existing tests pass
- [x] 5.4 Manual smoke test: run `wip worktree list` from a subdirectory of a wip project and confirm it works correctly
- [x] 5.5 Manual smoke test: run `wip root` from a subdirectory and confirm it prints the project root path
- [x] 5.6 Manual smoke test: run `wip worktree list` from outside any wip project and confirm a clear error is printed
