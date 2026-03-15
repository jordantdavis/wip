## Context

The `wip` CLI manages git submodules from the repository root. Git worktrees are a built-in git mechanism for checking out multiple branches of a repository simultaneously into separate directories. Each submodule is its own git repository, so worktrees must be created and managed within each submodule's git context, not the parent repo's.

The existing submodule commands establish clear patterns: validate inputs early, run git as a subprocess, stream output to the terminal, and exit with git's exit code. The worktree commands follow the same patterns with one additional complexity — git commands must run with the submodule directory as their working directory, not the repo root.

## Goals / Non-Goals

**Goals:**
- Add `wip worktree add`, `wip worktree list`, and `wip worktree remove` following existing CLI conventions
- Worktrees are stored at a conventional path: `<repo root>/worktrees/<submodule>/<worktree>/`
- The worktree name doubles as the branch name
- `list` discovers worktrees from the filesystem (no git state required)
- `remove` is safe by default; branch deletion is opt-in

**Non-Goals:**
- Managing worktrees that exist outside `worktrees/` (untracked by convention)
- Syncing or reconciling worktrees across machines
- Support for detached HEAD worktrees
- Any modification to existing submodule commands

## Decisions

### Working directory for git worktree commands

**Decision:** Set `cmd.Dir` to `<repo root>/<submodule>/` when constructing `exec.Cmd` for `git worktree` operations.

Git worktree commands operate on a specific repository — in this case the submodule's repo. The submodule is checked out at `<repo root>/<submodule>/`. Setting `cmd.Dir` on the `exec.Cmd` struct is the cleanest way to run git in a different directory without changing the process working directory.

**Alternative considered:** `cd` into the directory via a shell command. Rejected — unnecessary shell invocation and harder to control.

### Path construction

**Decision:** Resolve all paths as absolute using `os.Getwd()` at command entry. The worktree target path passed to git is `<abs repo root>/worktrees/<submodule>/<worktree>`.

Using absolute paths avoids ambiguity when `cmd.Dir` is set to the submodule directory — a relative path like `../../worktrees/...` would be fragile and confusing.

### Worktree discovery for `list`

**Decision:** Scan the filesystem at `worktrees/` using `os.ReadDir` — enumerate subdirectories (submodule names), then enumerate their subdirectories (worktree names).

**Alternative considered:** Run `git worktree list` inside each known submodule. Rejected — requires knowing which submodules exist and running N git processes. Filesystem scan is simpler and doesn't depend on submodule state. The `worktrees/` directory is entirely owned by `wip` by convention, so its contents are authoritative.

### Branch creation vs. existing branch

**Decision:** Default behavior creates a new branch (`git worktree add -b <worktree> <path>`). The `--existing-branch` flag switches to checkout mode (`git worktree add <path> <worktree>`), which requires the branch to already exist in the submodule repo.

Let git surface the error naturally if `--existing-branch` is used with a non-existent branch — consistent with how other commands delegate error handling to git.

### Branch deletion safety

**Decision:** `--delete-branch` uses `git branch -d` (safe delete). If the branch has unmerged commits, git will refuse and print an informative error. The user can handle that case manually.

Force deletion (`-D`) is not offered by `wip` — if needed the user can run git directly. This matches the principle of safe defaults.

### `worktrees/` directory creation

**Decision:** `worktree add` creates `worktrees/<submodule>/` via `os.MkdirAll` before invoking git if the path does not exist. This is transparent to the user.

### Worktree name validation

**Decision:** Worktree names must match `^[a-zA-Z0-9_-]+$`. This is stricter than git branch name rules but safe for both filesystem paths and branch names with no special handling required.

## Risks / Trade-offs

- **Filesystem as source of truth for `list`**: If a worktree directory exists but the git worktree registration was corrupted or manually removed, `list` will show it even though git doesn't know about it. Mitigation: document that `list` is filesystem-based; users can use `git worktree list` directly inside a submodule for authoritative state.
- **Submodule directory must exist**: `worktree add` requires the submodule to be checked out (not just registered). If the submodule exists in `.gitmodules` but was never initialized, git will error. Mitigation: let git's error message surface naturally — it will be informative.
- **git 2.5+ required**: `git worktree` was introduced in git 2.5 (2015). This is a reasonable minimum; no version check is added.
