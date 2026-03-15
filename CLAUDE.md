# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o wip .

# Run tests
go test ./...

# Run a single test
go test ./cmd/... -run TestFunctionName

# Vet
go vet ./...
```

## Architecture

`wip` is a Go CLI that provides a high-level interface over `git submodule` and `git worktree` operations for monorepo-style projects.

**Entry point:** `main.go` routes on `os.Args[1]` to `cmd.Init()`, `cmd.Submodule()`, or `cmd.Worktree()`.

**Command pattern:** Each subcommand family has a router file (`cmd/submodule.go`, `cmd/worktree.go`) that owns shared helpers and routes to individual operation files (`submodule_add.go`, `worktree_list.go`, etc.). Each operation owns its own `flag.FlagSet`.

**Git integration:** All git operations are delegated via `os/exec.Command("git", ...)` — stdout/stderr are streamed directly to the user and exit codes are propagated.

**Worktree directory convention:**
```
repo/
├── <submodule>/        ← actual submodule checkouts
└── worktrees/
    └── <submodule>/
        └── <worktree-name>/
```

**Concurrency:** `submodule sync` parallelizes updates across submodules using goroutines and `sync.WaitGroup`.

## OpenSpec

This project uses OpenSpec for spec-driven development. Specifications live in `openspec/specs/`, and change artifacts (proposals, designs, tasks) live in `openspec/changes/`. Use the `/openspec-propose`, `/openspec-apply-change`, and `/openspec-archive-change` skills to work within this workflow.
