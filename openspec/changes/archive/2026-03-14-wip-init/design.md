## Context

`wip` is a CLI tool for managing git submodules and worktrees. All existing commands (`submodule`, `worktree`) guard against running outside a git repository via `checkGitRepo()`. There is currently no way to bootstrap a new project — `init` fills that gap.

The command is intentionally minimal now but designed to be a staging point for future project-level initialization steps.

## Goals / Non-Goals

**Goals:**
- Add `wip init` as a top-level command that runs `git init` in the current directory if needed
- Make `init` idempotent — safe to run repeatedly as more steps are added over time
- Follow existing code structure and conventions

**Non-Goals:**
- Accepting a directory path argument (always operates on cwd)
- Creating any wip-specific config or directory structure (future concern)
- Producing verbose output beyond minimal confirmation

## Decisions

**Operate on cwd only**
No path argument. `wip init` always initializes the current directory. Rationale: keeps the interface simple and consistent with how `submodule` and `worktree` commands resolve the repo root via `repoRoot()` (which returns `os.Getwd()`).

**Repo root detection via `git rev-parse --git-dir`**
Rather than checking for `.git` directory presence, `init` runs `git rev-parse --git-dir` and branches on the result:
- Exits non-zero → not in any git repo → run `git init`
- Returns `.git` → current directory is the repo root → no-op, continue
- Returns anything else → inside a repo but not at root → error and exit non-zero

This approach is more robust than a raw `.git` check: it correctly handles the case where the user is inside a subdirectory or submodule of an existing repo, preventing accidental nesting. The error message in the third case is kept simple: "not at the root of a git repository".

**Single file: `cmd/init.go`**
No subcommands needed. The command lives in one file with an `Init(args []string)` function, matching the pattern of `cmd/submodule.go` and `cmd/worktree.go`.

## Risks / Trade-offs

**Bare repos not supported** → `git rev-parse --git-dir` returns an absolute path for bare repos, so they fall into the "not at root" error case. Bare repo support is out of scope.

**Silent on already-initialized** → No output when already at a git root. This is intentional for idempotency but means users get no feedback confirming the repo state. Acceptable for a tool used by developers who can run `git status` themselves.
