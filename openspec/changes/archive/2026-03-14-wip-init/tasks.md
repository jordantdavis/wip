## 1. Implement init command

- [x] 1.1 Create `cmd/init.go` with `Init(args []string)` function that runs `git rev-parse --git-dir` to determine repo state: exits non-zero → run `git init`; returns `.git` → no-op; returns anything else → error "not at the root of a git repository"
- [x] 1.2 Register `case "init"` in `main.go` dispatch switch
- [x] 1.3 Add `init` to the usage output in `main.go`
