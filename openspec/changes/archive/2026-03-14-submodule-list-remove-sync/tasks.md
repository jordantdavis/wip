## 1. Shared: .gitmodules parsing

- [x] 1.1 Implement `parseSubmodules() ([]submodule, error)` that runs `git config --file .gitmodules --get-regexp` to extract all submodule names and URLs
- [x] 1.2 Implement `submoduleExists(name string) (bool, error)` that checks `.gitmodules` for a named entry (used by remove and sync --name)

## 2. submodule list

- [x] 2.1 Implement `submoduleList(args []string)` with `flag.NewFlagSet`, calling `checkGitRepo()` then `parseSubmodules()`
- [x] 2.2 Sort results alphabetically by name and print `<name>  <url>` per line
- [x] 2.3 Print `no submodules found` and exit 0 when list is empty
- [x] 2.4 Wire `list` case into the `submodule` routing switch in `cmd/submodule.go`

## 3. submodule remove

- [x] 3.1 Implement `submoduleRemove(args []string)` with `flag.NewFlagSet`, calling `checkGitRepo()`, requiring name positional arg, then calling `submoduleExists()`
- [x] 3.2 Run `git submodule deinit -f <name>` as subprocess streaming stdout/stderr; stop and exit on non-zero
- [x] 3.3 Run `git rm -f <name>` as subprocess streaming stdout/stderr; stop and exit on non-zero
- [x] 3.4 Remove `.git/modules/<name>` with `os.RemoveAll`; print error and exit on failure
- [x] 3.5 Wire `remove` case into the routing switch

## 4. submodule sync

- [x] 4.1 Implement `submoduleSync(args []string)` with `flag.NewFlagSet` and `--name` string flag, calling `checkGitRepo()`
- [x] 4.2 When `--name` provided: validate existence with `submoduleExists()`, run single `git submodule update --init --remote <name>`, print `✓`/`✗` result, exit accordingly
- [x] 4.3 When no `--name`: call `parseSubmodules()`, print message and exit 0 if empty
- [x] 4.4 Launch one goroutine per submodule running `git submodule update --init --remote <name>`, capturing combined stdout+stderr into a per-result buffer using `sync.WaitGroup`
- [x] 4.5 After all goroutines complete, print `✓ <name>` or `✗ <name>: <captured output>` for each result
- [x] 4.6 Exit non-zero if any submodule failed
- [x] 4.7 Wire `sync` case into the routing switch
