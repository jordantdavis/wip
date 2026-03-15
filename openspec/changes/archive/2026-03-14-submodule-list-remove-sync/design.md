## Context

`wip submodule add` is implemented and working. It uses `git submodule add` to register submodules, with an optional `--name` flag that sets both the git-internal submodule name and the checkout directory. All submodule metadata lives in git's native `.gitmodules` file — there is no custom persistence layer.

The project uses the standard library `flag` package for flag parsing, manual switch routing at two levels (`main.go` → `cmd/submodule.go`), and `os/exec` for subprocess delegation. No external dependencies.

## Goals / Non-Goals

**Goals:**
- Add `list`, `remove`, and `sync` to the existing `submodule` command
- Reuse existing patterns: `flag.NewFlagSet`, `checkGitRepo()`, `os/exec` subprocess delegation
- Keep zero external dependencies

**Non-Goals:**
- Submodule pinning or branch tracking configuration
- Interactive selection of submodules
- Progress bars or TUI output for sync
- Shell completions

## Decisions

### 1. Read `.gitmodules` via `git config --file .gitmodules` subprocess

`git config --file .gitmodules --get-regexp 'submodule\..*\.name'` yields all submodule names. A second call with `\.url` yields URLs. This avoids writing a custom INI parser and stays consistent with the pattern of delegating to git.

Alternative considered: parse `.gitmodules` directly with `bufio.Scanner`. Rejected — more brittle, requires handling git's INI quoting edge cases.

### 2. `remove` validates existence before touching anything

The command checks `.gitmodules` for the named submodule before running any destructive git commands. This gives a clean, friendly error instead of letting `git submodule deinit` fail with a cryptic message.

### 3. `remove` streams git output; `sync` buffers it

`remove` is a sequential three-step operation — streaming matches the existing `add` pattern and is appropriate for a single linear process. `sync` runs operations concurrently, so streaming would produce interleaved output. Sync collects each subprocess's combined stdout+stderr into a buffer, then prints structured results after all goroutines complete.

### 4. `sync` uses goroutines + `sync.WaitGroup`

One goroutine per submodule, results collected into a slice with a mutex or pre-allocated by index (index is safe since slice length is fixed before launch). `sync.WaitGroup` waits for all. No worker pool — submodule counts are small and the bottleneck is network I/O, not CPU.

### 5. `sync` exit code is non-zero if any submodule failed

Matches the pattern of `git submodule update` itself. Lets CI scripts detect partial failures without parsing output.

### 6. `list` sorts alphabetically by name

Deterministic output regardless of `.gitmodules` order. Simple `sort.Strings` on the name slice.

## Risks / Trade-offs

- **`.gitmodules` parsing via subprocess** — two `git config` calls per `list`/`sync` invocation adds slight latency. Acceptable for a personal tool with small submodule counts.
- **`rm -rf .git/modules/<name>` in `remove`** — destructive and not reversible. Mitigated by the existence check before any removal step begins. No further mitigation; this is the correct and documented removal procedure.
- **Concurrent sync with no concurrency limit** — if someone has dozens of submodules, all network calls fire simultaneously. Acceptable for a personal tool; a semaphore can be added later if needed.
