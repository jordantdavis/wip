## Context

`wip` has a fully working `version` command routed in `main.go` and implemented in `cmd/version.go`. The `printUsage()` function in `main.go` lists available commands, but `version` is absent from that list. The fix is a single-line addition.

## Goals / Non-Goals

**Goals:**
- `version` appears in `wip` help output at the same indentation and style as existing commands

**Non-Goals:**
- Adding `--version` / `-v` flag support
- Changing the `version` command's output format
- Moving `printUsage` or restructuring help output

## Decisions

**Add `version` last in the command list** — `version` is a meta-command (about the tool itself), not an operational command. Convention places it after operational commands (`init`, `submodule`, `worktree`).

**Description text: `"print version information"`** — matches the verb-first, lowercase style of existing descriptions and accurately reflects the command's output.

## Risks / Trade-offs

None. The change is purely additive and touches a single line in `printUsage()`.
