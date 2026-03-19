## Context

`wip worktree add` currently validates worktree names with a hardcoded regex (`^[a-zA-Z0-9_-]+$`) and uses the name as both the git branch name and the filesystem path component. This prevents common branch naming conventions like `feature/my-thing`. The branch name and path segment need to be decoupled: the user always thinks in branch names, and the tool derives a safe path component from them.

Three commands are affected: `worktree add`, `worktree remove`, and `worktree list`.

## Goals / Non-Goals

**Goals:**
- Support any branch name that `git check-ref-format` considers valid.
- Derive the worktree path by replacing `/` with `-` in the branch name.
- Branch name is the canonical user-facing input for both `add` and `remove`.
- `worktree list` shows submodule, path segment, and branch name.

**Non-Goals:**
- Detecting or warning about path collisions (e.g., `feature/foo` and `feature-foo` mapping to the same path). Git and the filesystem will naturally reject duplicate paths.
- Changing the worktrees directory layout.
- Supporting branch name changes after a worktree is created.

## Decisions

### 1. Delegate branch name validation to `git check-ref-format`

**Decision:** Replace `validateWorktreeName` with `validateBranchName`, which runs `git check-ref-format --branch <name>` and reports an error if it exits non-zero.

**Alternatives considered:**
- Implement a regex approximating git's rules. Rejected: git's `check-ref-format` rules are complex and have edge cases (e.g., `..`, `@{`, `.lock` suffix, control characters). Maintaining a parallel implementation invites drift.
- No validation at all, let git fail. Rejected: the error messages from git are less user-friendly and don't mention wip's CLI conventions.

**Rationale:** Git already runs as a subprocess throughout this codebase. One extra subprocess call is negligible, and correctness is guaranteed.

### 2. Path derivation: replace `/` with `-`

**Decision:** `worktreePath := strings.ReplaceAll(branchName, "/", "-")` is the sole transformation applied. The result is used as the final directory name under `worktrees/<submodule>/`.

**Alternatives considered:**
- URL-encode or percent-encode slashes. Rejected: ugly and non-obvious on the filesystem.
- Use only the last segment after the final `/`. Rejected: loses context (`feature/auth` and `bugfix/auth` would both become `auth`), causing likely collisions.

**Rationale:** Simple, reversible in a user's head, consistent with common conventions (GitHub Actions branch-to-path transforms use the same approach).

### 3. Branch name is the canonical input for `worktree remove`

**Decision:** `worktree remove` accepts the branch name (e.g., `feature/my-thing`), derives the path by applying the same `/`→`-` replacement, and uses the original branch name for `--delete-branch`.

**Alternatives considered:**
- Accept the path form (e.g., `feature-my-thing`). Rejected: the path form is a derived internal detail; `--delete-branch` would need to reconstruct the branch name, which is ambiguous (`feature-foo` could come from `feature/foo` or `feature-foo`).

**Rationale:** Consistent UX — users pass the same string to `add` and `remove`. The mapping is the tool's responsibility.

### 4. `worktree list` reads branch via `git branch --show-current`

**Decision:** For each worktree directory discovered on disk, run `git -C <abs-path> branch --show-current` to get the checked-out branch name. Output three columns: `<submodule>  <path-segment>  <branch>`.

**Alternatives considered:**
- Read `<worktree>/.git` file then parse the gitdir's `HEAD` ref directly. Possible, but brittle — git's internal ref format could change. `branch --show-current` is the stable public interface.
- Show only the branch name instead of the path segment. Rejected: the path segment is what users use with `remove`, so showing it aids discoverability.

**Rationale:** `git branch --show-current` is the canonical way to read the current branch in a git worktree. The extra subprocess per worktree is acceptable for a listing command.

## Risks / Trade-offs

- **Path collisions**: `feature/foo` and `feature-foo` map to the same path. The filesystem or git will reject the second creation. No extra handling needed, but the error message will come from git and may be slightly confusing. → Acceptable for now; a future improvement could detect this at validation time.
- **Detached HEAD worktrees**: If a worktree is in detached HEAD state, `git branch --show-current` returns an empty string. `list` will show a blank branch column. → Acceptable; detached HEAD is an unusual state for wip-managed worktrees.
- **`git check-ref-format` subprocess cost**: One extra git call per `add`/`remove` invocation. Negligible.

## Migration Plan

No data migration needed. The directory layout does not change for worktrees created under the old naming rules (those names were already valid under the new rules too, since `[a-zA-Z0-9_-]+` is a subset of valid git branch names).

Existing worktrees continue to work. `worktree remove` with a name like `my-feature` will still find `worktrees/<sub>/my-feature` correctly (no slash replacement needed when there are no slashes).
