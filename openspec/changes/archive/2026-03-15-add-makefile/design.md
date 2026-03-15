## Context

The repo has no Makefile. Build commands exist only as documentation in CLAUDE.md. GitHub Actions workflows will soon need a standard interface for CI steps.

## Goals / Non-Goals

**Goals:**
- Provide a single `Makefile` with targets covering the full dev and CI workflow
- Distinguish between `fmt` (local, rewrites) and `fmt-check` (CI, read-only)
- Keep targets independently callable so GH Actions can invoke them with granularity

**Non-Goals:**
- Cross-platform builds or release packaging
- A `ci` aggregate target (CI calls individual targets directly)
- Docker or container-based build targets

## Decisions

**`fmt` vs `fmt-check` as separate targets**
`go fmt ./...` rewrites files in place — useful locally but wrong for CI (would silently "fix" what should be a failure). `gofmt -l .` lists unformatted files without modifying them; CI can fail on non-empty output. Keeping them as separate targets avoids the need for a flag or env var to switch modes.

**`check` uses `fmt-check` not `fmt`**
The local `check` aggregate (fmt-check + vet + test) mirrors what CI verifies. If `check` used `fmt` instead, passing locally would not guarantee passing in CI.

**No default target**
All targets must be called explicitly. This avoids ambiguity about what bare `make` does and forces intentional use.

**`go install .` for install target**
Installs from local source into `$GOPATH/bin`. Idiomatic for CLI tools under active development.

## Risks / Trade-offs

- `gofmt -l .` recurses into all subdirectories including vendored or generated code — not a concern now but worth noting if those are added later → Mitigation: scope with explicit paths if needed
- No default target means slightly more typing — acceptable given the explicitness preference
