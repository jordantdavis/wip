## Context

`wip` is a CLI for managing git submodule + worktree workflows in monorepo-style projects. The current `wip submodule` command wraps raw git submodule operations, exposing SHA pinning, detached HEAD state, and manual sync complexity to the user.

The primary use case for `wip submodule` is co-locating external repos in a single directory so AI coding tools (Claude Code, OpenCode) can `@` reference their files during planning sessions. Git's file picker in these tools requires files to be tracked (not gitignored and not merely untracked) for `@` completion to work.

## Goals / Non-Goals

**Goals:**
- Replace `wip submodule` with `wip ref` — same git mechanism, better defaults, cleaner surface
- Always track branch HEAD; never require committing a SHA update to advance a ref
- Keep git status clean on the parent repo; never surface refs as modified/dirty
- Make post-clone setup a single command (`wip ref restore`) for teammates
- Preserve worktree hook config (`on-worktree-create`, `on-worktree-launch`)

**Non-Goals:**
- Supporting SHA pinning or reproducible builds (refs are for read-only AI context, not dependency locking)
- Network fetches or remote-only context (all refs are local clones)
- Sparse or shallow checkouts (full clones are required for complete AI context)
- Supporting git repos that are not on the local machine

## Decisions

### Keep git submodules as the underlying mechanism

**Decision:** Implement `wip ref` on top of git submodules.

**Why:** Git submodules create a gitlink entry (mode 160000) in the parent repo's tree — a single tracked object that represents the submodule directory without tracking its file contents. This is the only git primitive that satisfies all three requirements simultaneously: (1) the directory is tracked (so AI tools' git-based file pickers include it), (2) the file contents are not committed to the parent repo (so `git diff` stays clean), and (3) a URL is stored in `.gitmodules` (so the setup is reconstructable from the repo).

**Alternatives considered:**
- Plain clones at root, not gitignored: Dirs appear in `git status` as untracked — acceptable line count, but risky for `git add .` and noisy for teammates.
- Gitignored clones: Clean git status, but both Claude Code and OpenCode exclude gitignored files from `@` completions. The `respectGitignore: false` override in Claude Code has been broken since v2.0.75.
- Symlinks to external dirs: Claude Code's file picker does not recurse into symlinks pointing outside the worktree in git repos (closed "not planned"). OpenCode's symlink support was reverted.

### Configure submodules with `branch` tracking and `ignore = all`

**Decision:** All refs are added with `branch = main` (configurable via `--branch`) and `ignore = all` written to `.gitmodules`.

**Why:** `branch` tracking enables `git submodule update --remote` to pull the latest commit on the named branch without any SHA committed in the parent repo. `ignore = all` tells git to never report the submodule as modified in `git status` or `git diff`, regardless of divergence between the committed gitlink SHA and the current submodule state. Together these mean: sync is a single command, and the parent repo's git status stays clean at all times.

**Alternatives considered:**
- Default SHA pinning: Requires committing a SHA update to advance the ref. Generates noise commits ("update api pointer") with no semantic value for a read-only planning context.
- `ignore = dirty` only: Suppresses uncommitted working tree changes in the submodule but still reports "new commits" drift. Would still produce modified status after every `wip ref sync`.

### `wip ref restore` always uses `--remote`

**Decision:** `wip ref restore` runs `git submodule update --init --remote` rather than `git submodule update --init`.

**Why:** The committed SHA in the parent repo may be stale (refs are never explicitly pinned). Using `--remote` fetches the current branch HEAD for each ref, giving teammates the same view of the codebase rather than an arbitrary historical SHA.

### `wip ref` command family — no `git submodule` passthrough

**Decision:** The `wip ref` command family covers add, list, remove, sync, and restore. No general `git submodule` passthrough is provided.

**Why:** The goal is to hide submodule complexity entirely. Users should never need to run `git submodule` directly for ref management. The five operations cover the complete lifecycle.

### Config schema: `submodules` → `refs`, add `branch` field

**Decision:** Rename the top-level `.wip.yml` key from `submodules` to `refs`. Rename `SubmoduleConfig` to `RefConfig`. Add an optional `branch` field (default: `main`).

**Why:** The config key should match the command name. The `branch` field is needed to store the configured tracking branch per-ref, since different repos may use `main`, `master`, or a custom branch.

## Risks / Trade-offs

**`ignore = all` hides real divergence** → If a ref's remote moves significantly (breaking changes, force-push), `wip ref sync` will pull silently. Mitigation: `wip ref list` shows current SHA and branch for manual inspection; sync output reports each ref's result.

**Breaking change for existing `.wip.yml` files** → Projects using `wip submodule` must rename `submodules:` to `refs:` in `.wip.yml`. Mitigation: document the migration clearly; the rename is mechanical and one-time.

**`ignore = all` means `git submodule status` won't show refs as initialized** → Teammates may not know refs exist until they run `wip ref list`. Mitigation: `wip ref restore` prints each ref as it initializes; `wip ref list` reads `.gitmodules` directly and is always authoritative.

## Migration Plan

1. Remove `cmd/submodule*.go` files
2. Add `cmd/ref*.go` files implementing the new command surface
3. Update `cmd/config.go`: rename `WipConfig.Submodules` → `WipConfig.Refs`, `SubmoduleConfig` → `RefConfig`, add `Branch string` field
4. Update `main.go`: route `ref`, remove `submodule` route
5. Update `wip init` to scaffold `refs:` instead of `submodules:` in `.wip.yml`
6. Update `wip-config` spec (delta: rename schema key, add branch field)
7. Existing `.wip.yml` files require manual rename of `submodules:` → `refs:`

No rollback strategy is needed — this is a local CLI with no network-facing components. Teams can pin to a prior binary if needed.

## Open Questions

- Should `wip ref restore` skip already-initialized refs, or always run `--remote` to update them too? Current plan: always run `--remote` so restore is idempotent and also functions as a sync.
- Should `--branch` default to `main` hardcoded, or should `wip ref add` detect the remote's default branch? Starting with `main` as default; can add detection later.
