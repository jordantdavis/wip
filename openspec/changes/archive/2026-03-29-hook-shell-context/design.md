## Context

Hook commands in `.wip.yml` are currently executed via `exec.Command(parts[0], parts[1:]...)` after splitting the hook string with `strings.Fields`. This gives hooks no access to contextual values (ref name, worktree name, paths) and prevents shell features like pipes, redirects, and `&&` chaining from working.

Both `worktree_add.go` and `worktree_launch.go` share the same hook execution pattern — a loop that splits, runs, and prints ✓/✗. The fix applies identically to both.

## Goals / Non-Goals

**Goals:**
- Switch hook execution to `sh -c "<hook>"` so hooks are full shell commands
- Inject a standard set of `WIP_*` env vars before each hook runs
- Apply consistently to both `on-worktree-create` and `on-worktree-launch`
- Establish the `WIP_*` env vars as a forward contract — future hooks inherit the same set

**Non-Goals:**
- Changing hook configuration syntax in `.wip.yml`
- Adding new hook types
- Supporting non-POSIX shells (cmd.exe, PowerShell)

## Decisions

### Use `sh -c` instead of direct exec

**Decision:** Execute hooks as `exec.Command("sh", "-c", hook)` rather than splitting with `strings.Fields`.

**Rationale:** Hook strings in `.wip.yml` are written as shell commands. Users expect shell features — `&&`, pipes, redirects, variable expansion. Direct exec breaks all of these silently. `sh -c` is how npm, git hooks tooling, and Make all execute inline command strings.

**Alternative considered:** Go template interpolation (`{{.RefName}}`) into the command string before direct exec. Rejected because it introduces unfamiliar syntax, doesn't support shell features, and breaks for values containing spaces. `sh -c` + env vars solves all the same problems more idiomatically.

### Inject env vars, not template variables

**Decision:** Provide context as `WIP_*` environment variables set on the hook subprocess, not as template substitutions in the hook string.

**Rationale:** Env vars are the universal Unix convention for ambient context in child processes — used by git (`GIT_DIR`, `GIT_WORK_TREE`), npm (`npm_package_name`), and every CI system. They're inherited by all subprocesses the hook spawns, work in any language, and don't require any special syntax in the hook string itself. With `sh -c` execution, users can reference them as `$WIP_REF_NAME` inline when needed.

### Standard env var set

All hook executions SHALL inject these variables:

| Variable | Value | Example |
|---|---|---|
| `WIP_REF_NAME` | The ref name as given in `.wip.yml` | `backend` |
| `WIP_WORKTREE_NAME` | The worktree name argument | `my-feature` |
| `WIP_WORKTREE_PATH` | Absolute path to the worktree directory | `/Users/jordan/proj/worktrees/backend/my-feature` |
| `WIP_ROOT` | Absolute path to the repo root | `/Users/jordan/proj` |

These four variables are available in both `on-worktree-create` and `on-worktree-launch`. Future hook types SHALL also inject this same set (extended as needed).

### Hook label shows the raw hook string

**Decision:** The ✓/✗ output label shows the hook string as written in config, not a shell-expanded form.

**Rationale:** The config string is stable and matches what the user wrote. It's the right identifier for "which hook ran." The expanded form (with env var values substituted) is available in shell output if the hook itself echoes anything.

## Risks / Trade-offs

**Shell injection via hook strings** → Hook strings come from `.wip.yml` which is a project config file committed to the repo. Users who can edit `.wip.yml` already have arbitrary code execution. No additional attack surface.

**`sh` not available** → `sh` is present on all POSIX systems. This is the same assumption git makes for hooks. Not a practical concern.

**Hooks that relied on direct exec behavior** → Existing hooks like `"npm install"` or `"echo hello"` work identically under `sh -c`. No migration needed.

**Spaces in env var values** → With `sh -c`, users referencing `$WIP_WORKTREE_PATH` must quote it (`"$WIP_WORKTREE_PATH"`) if the path contains spaces, just as with any shell variable. This is standard shell practice, not a new constraint.
