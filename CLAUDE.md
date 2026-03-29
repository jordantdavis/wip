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

`wip` is a Go CLI that provides a high-level interface over git submodule and worktree operations for monorepo-style projects, exposing them as `wip ref` and `wip worktree` commands.

**Entry point:** `main.go` routes on `os.Args[1]` to `cmd.Init()`, `cmd.Ref()`, or `cmd.Worktree()`.

**Command pattern:** Each subcommand family has a router file (`cmd/ref.go`, `cmd/worktree.go`) that owns shared helpers and routes to individual operation files (`ref_add.go`, `worktree_list.go`, etc.). Each operation owns its own `flag.FlagSet`.

**Git integration:** All git operations are delegated via `os/exec.Command("git", ...)` — stdout/stderr are streamed directly to the user and exit codes are propagated.

**Worktree directory convention:**
```
repo/
├── <ref>/          ← ref checkouts (branch-tracking, ignore=all)
└── worktrees/
    └── <ref>/
        └── <worktree-name>/
```

**Concurrency:** `ref sync` and `ref restore` parallelize updates across refs using goroutines and `sync.WaitGroup`.

**Hook execution:** Hooks are run via `runHooks` in `cmd/hooks.go`. Each hook string is executed as `sh -c "<hook>"` (not split with `strings.Fields`), enabling shell features like `&&`, pipes, and redirects. The following env vars are injected into every hook subprocess: `WIP_REF_NAME`, `WIP_WORKTREE_NAME`, `WIP_WORKTREE_PATH`, `WIP_ROOT`. Hooks are configured exclusively in `.wip.yml` — there are no CLI flags for setting hooks.

## OpenSpec

This project uses OpenSpec for spec-driven development. Specifications live in `openspec/specs/`, and change artifacts (proposals, designs, tasks) live in `openspec/changes/`. Use the `/openspec-propose`, `/openspec-apply-change`, and `/openspec-archive-change` skills to work within this workflow.
