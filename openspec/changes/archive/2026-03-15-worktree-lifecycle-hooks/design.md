## Context

`wip` is currently stateless — no configuration file exists. All state is derived from `.gitmodules` and the filesystem. Adding lifecycle hooks requires introducing a config layer for the first time. The design must be minimal and non-breaking: a single `.wip.yml` at the repo root, processed with `go.yaml.in/yaml/v4`.

## Goals / Non-Goals

**Goals:**
- Introduce `.wip.yml` as the `wip` project config file, scaffolded by `wip init`
- Enforce `.wip.yml` presence in `wip submodule add` and `wip worktree add`
- Support per-submodule `on-worktree-create` hook lists, configured via `wip submodule add --on-worktree-create`
- Execute `on-worktree-create` hooks in order inside the newly created worktree directory
- Fail gracefully on hook errors (warn, leave worktree intact)

**Non-Goals:**
- Per-worktree hook overrides
- Hooks for other lifecycle events (on-remove, on-sync, etc.)
- `.wip.yml` validation/linting command
- `wip submodule list/remove/sync` or `wip worktree list/remove` requiring `.wip.yml`

## Decisions

### 1. `.wip.yml` schema

```yaml
submodules:
  <name>:
    on-worktree-create:
      - <command>
      - <command>
```

**Rationale:** Template-based (per submodule, not per worktree instance). Every worktree created from a given submodule gets the same initialization. Kebab-case keys throughout, consistent with CLI flag naming. The file is minimal — only submodules with hooks need entries.

**Alternative considered:** Per-worktree instance records in `.wip.yml`. Rejected — adds complexity and creates a second source of truth that can drift from the filesystem.

### 2. `.wip.yml` as required project marker

`wip submodule add` and `wip worktree add` SHALL fail with a nudge to `wip init` if `.wip.yml` is absent. This makes `.wip.yml` the canonical marker of a `wip`-managed repo (analogous to `go.mod` or `package.json`).

**Rationale:** Clear contract — if you're using `wip add` commands, you've opted into the `wip` project model. Hard failure with a clear message is better than silently skipping hooks or auto-creating the file mid-command.

**Alternative considered:** Auto-create `.wip.yml` on first use. Rejected — side effects in non-init commands are surprising and harder to reason about.

### 3. Shared validation helper

A package-level function `requireWipConfig() (*WipConfig, error)` (or similar) reads and parses `.wip.yml`. Both `submodule_add.go` and `worktree_add.go` call it at the top of their execution. This avoids duplicating the check.

### 4. `--on-worktree-create` as a repeatable flag

Implemented via a custom `flag.Value` type (`stringList`) that appends on each use. Order is preserved by CLI argument order. Maps directly to the YAML list. Flag is positioned before the URL, consistent with `--name`.

```bash
wip submodule add --on-worktree-create "npm install" --on-worktree-create "cp .env.example .env" <url>
```

### 5. Hook execution model

Each command in the `on-worktree-create` list is run via `exec.Command` with `cwd` set to the newly created worktree directory. Commands are run sequentially in list order. Each command is split into program + args using `strings.Fields`. If a command fails (non-zero exit), a warning is printed to stderr and execution continues with the remaining commands.

**Alternative considered:** Run all commands as a single shell string via `sh -c`. Rejected — avoids shell dependency and is more explicit about what's being run. If users need shell features (`&&`, pipes), they can wrap in a script.

**Alternative considered:** Abort on first hook failure and remove the worktree. Rejected — leaves the user with a usable worktree and a clear warning; they can fix and re-run the hook manually.

### 6. YAML dependency

`go.yaml.in/yaml/v4` is added as the first external dependency. Chosen per the user's explicit requirement. This is the maintained successor to `gopkg.in/yaml.v3`.

## Risks / Trade-offs

- **`.wip.yml` out of sync with filesystem**: If a worktree is deleted manually without updating `.wip.yml`, the config entry for that submodule remains. This is intentional — `.wip.yml` is a config file, not a registry. No reconciliation needed.
- **Command splitting with `strings.Fields`**: Commands containing quoted arguments with spaces (e.g., `cp "my file" dest`) will not split correctly. Mitigation: document that complex commands should be wrapped in a shell script.
- **First external dependency**: Introduces `go.sum` and module download requirements. Low risk for a well-maintained YAML library.

## Migration Plan

No migration needed. `.wip.yml` is a new file. Existing repos without it will get a clear error message on `add` commands directing them to run `wip init`.

## Open Questions

None — all decisions resolved in exploration.
