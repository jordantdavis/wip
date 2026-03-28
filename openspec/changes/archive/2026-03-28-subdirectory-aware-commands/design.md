## Context

All `wip` subcommands currently assume the process's current working directory is the project root — the directory containing `.wip.yml`. This assumption is baked into three places:

1. `loadWipConfig` / `saveWipConfig` — read/write `.wip.yml` relative to cwd
2. `checkGitRepo` — checks for `.git` in cwd
3. `repoRoot` — returns `os.Getwd()`

Git commands that run without an explicit `Dir` (e.g., `git submodule add`, `git config --file .gitmodules`) also implicitly run against cwd.

The fix is to locate the true project root early and normalize cwd before any subcommand logic runs.

## Goals / Non-Goals

**Goals:**
- `wip submodule` and `wip worktree` commands work from any subdirectory of a wip project
- Project root is discovered by walking up from cwd, stopping at the user's home directory
- Invocations from outside the home directory fail immediately with a clear error
- All existing subcommand behavior is preserved exactly — only the working directory is affected

**Non-Goals:**
- `wip init` participating in project discovery (it creates a project, not uses one)
- `wip version` participating in project discovery (no project context needed)
- Supporting projects outside the user's home directory
- Handling symlinked paths specially

## Decisions

### D1: `os.Chdir` to project root before dispatch (Option Z)

Rather than threading a `root` string through every function that spawns a git process, the main dispatch normalizes cwd to the project root once, before any subcommand runs. All existing code that relies on cwd (git exec calls without `Dir`, relative path construction) then works without modification.

**Alternatives considered:**
- Pass `root` explicitly to every git `exec.Command` as `cmd.Dir` — correct but requires touching ~10 call sites including helpers like `parseSubmodules` and `submoduleExists` that have no current need for a root parameter.
- Store root in a package-level variable — simpler than threading but global mutable state is harder to test and reason about.

### D2: `WipProject` struct bundles root + config (Option C)

`findWipProject()` returns a `*WipProject{Root string, Config *WipConfig}`. The config is loaded during the walk, so there's no second read. Callers in `main.go` get everything they need in one call.

```
type WipProject struct {
    Root   string      // absolute path to directory containing .wip.yml
    Config *WipConfig  // parsed contents of .wip.yml
}
```

After `os.Chdir(project.Root)`, `requireWipConfig()` continues to work unchanged (reads `.wip.yml` from the now-correct cwd), so subcommand internals need no changes.

### D3: Home directory as the walk ceiling

The upward walk is bounded by `os.UserHomeDir()`. If cwd is outside the home directory, the function returns an error immediately (no walk attempted). The check uses path prefix matching with a separator guard to avoid false matches (e.g., `/Users/jordanfoo` must not match home `/Users/jordan`):

```
dir == home || strings.HasPrefix(dir, home+string(os.PathSeparator))
```

Within the home boundary, the walk proceeds upward checking each directory for `.wip.yml`, stopping when it finds one or when it reaches `home`. A filesystem root guard (`filepath.Dir(dir) == dir`) is kept as a safety net.

**Alternatives considered:**
- Walk to filesystem root (`/`) — works but could pick up stray `.wip.yml` files in system directories; semantically wrong for a developer tool.
- Walk to git root — couples discovery to git, and a wip project root and git root are the same by construction anyway.

## Risks / Trade-offs

- **`os.Chdir` is a process-wide side effect** → Acceptable because `wip` is a short-lived CLI process. The chdir happens once, early, before any real work. No goroutines are running at that point.
- **Tests use `os.Chdir` to set up temp environments** → Since `findWipProject` checks the immediate directory first (before walking up), tests that chdir to a temp dir containing `.wip.yml` will find it immediately. No test changes expected.
- **User runs `wip` from exactly their home directory** → The walk checks `home` itself before stopping, so a `.wip.yml` at `~` is found correctly. Unusual but valid.

## Open Questions

None. All decisions are resolved.
