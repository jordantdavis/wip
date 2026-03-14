## 1. Project Setup

- [x] 1.1 Create `cmd/` directory structure

## 2. CLI Foundation

- [x] 2.1 Rewrite `main.go` with subcommand router (switch on `os.Args[1]`, delegate to `cmd/`)
- [x] 2.2 Add usage output and non-zero exit for missing or unknown subcommand

## 3. Submodule Command

- [x] 3.1 Create `cmd/submodule.go` with subcommand router for `submodule` (switch on first remaining arg)
- [x] 3.2 Add usage output and non-zero exit for missing or unknown submodule subcommand

## 4. Submodule Add — Validation

- [x] 4.1 Implement git repo check (`.git` exists in working directory)
- [x] 4.2 Implement URL validation (non-empty, matches `https://`, `http://`, `git://`, or `git@` SSH form)
- [x] 4.3 Implement path validation (resolve to absolute, verify within cwd, verify does not exist)

## 5. Submodule Add — Execution

- [x] 5.1 Implement `add` subcommand: parse positional args (`url` required, `path` optional)
- [x] 5.2 Run all validations in order (git repo → url → path)
- [x] 5.3 Execute `git submodule add <url> [<path>]` as subprocess, stream stdout/stderr, propagate exit code
