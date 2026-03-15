## 1. Dependencies

- [x] 1.1 Add `go.yaml.in/yaml/v4` to go.mod and go.sum via `go get`

## 2. Config Package

- [x] 2.1 Create `cmd/config.go` with the `WipConfig` struct matching the `.wip.yml` schema (`Submodules map[string]SubmoduleConfig`, `OnWorktreeCreate []string`)
- [x] 2.2 Implement `loadWipConfig() (*WipConfig, error)` that reads and parses `.wip.yml` from the current directory
- [x] 2.3 Implement `requireWipConfig() (*WipConfig, error)` that calls `loadWipConfig` and returns a user-facing error with a `wip init` nudge if the file is absent
- [x] 2.4 Implement `saveWipConfig(cfg *WipConfig) error` that marshals and writes `.wip.yml`

## 3. wip init

- [x] 3.1 After the git init logic in `cmd/init.go`, scaffold `.wip.yml` with an empty `submodules` map if the file does not already exist

## 4. wip submodule add

- [x] 4.1 Add `stringList` custom `flag.Value` type to `cmd/config.go` (or a shared file) for repeatable flags
- [x] 4.2 Add `--on-worktree-create` repeatable flag (before URL, after `--name`) to the `flag.FlagSet` in `cmd/submodule_add.go`
- [x] 4.3 Call `requireWipConfig()` at the start of `SubmoduleAdd` execution (after flag parse, before git logic)
- [x] 4.4 After successful `git submodule add`, if `--on-worktree-create` was provided, update the config struct with the commands and call `saveWipConfig`

## 5. wip worktree add

- [x] 5.1 Call `requireWipConfig()` at the start of `WorktreeAdd` execution (after flag parse, before git logic)
- [x] 5.2 After successful `git worktree add`, look up `cfg.Submodules[submoduleName].OnWorktreeCreate`
- [x] 5.3 For each command in the list, run it via `exec.Command` with `cwd` set to the new worktree path, streaming stdout/stderr to the terminal
- [x] 5.4 On non-zero exit from any hook command, print a warning to stderr and continue; do not remove the worktree

## 6. Tests

- [x] 6.1 Test `loadWipConfig` with a valid `.wip.yml`, a missing file, and a malformed file
- [x] 6.2 Test `requireWipConfig` returns the correct error message when `.wip.yml` is absent
- [x] 6.3 Test `saveWipConfig` round-trips correctly through `loadWipConfig`
- [x] 6.4 Test `stringList` flag accumulates values in order
