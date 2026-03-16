## Context

`wip` is a Go CLI with no version identity today. It is being prepared to share with a team. The module is publicly hosted (`github.com/jordantdavis/wip`), so both `go install` and pre-built binary downloads are viable distribution paths. The project uses a hand-rolled Makefile and GitHub Actions CI.

## Goals / Non-Goals

**Goals:**
- Embed version string at build time via ldflags
- Dev builds auto-stamp a `dev-<timestamp>` version; release builds use the git tag
- `go install` path surfaces the correct version via `debug.ReadBuildInfo()`
- `wip version` prints version + OS/arch (e.g. `wip v0.0.1 (darwin/arm64)`)
- GitHub Actions releases darwin/arm64 binary on semver tag push
- Artifacts: `.tar.gz` archive + `checksums.txt` per release

**Non-Goals:**
- Homebrew tap (future)
- Linux or Windows builds (future — architecture is designed to expand)
- Changelog authoring (GitHub auto-generated notes are sufficient for now)
- Code signing or notarization

## Decisions

### Version variable lives in `cmd` package

**Decision**: Declare `var version = "unknown"` in `cmd/version.go` and use ldflags path `github.com/jordantdavis/wip/cmd.version`.

**Alternatives considered**:
- `main.go` — simpler path (`-X main.version`), but `main` is routing-only; putting state there is inconsistent with the existing command pattern where each command family owns its own file.

**Rationale**: Keeps `cmd/` as the single home for all command logic. `main.go` stays a pure router.

---

### Two-layer version resolution

**Decision**: Check ldflags-injected `version` first; fall back to `debug.ReadBuildInfo()` for `go install` users.

```
version var == "unknown"?
  └─ yes → ReadBuildInfo().Main.Version != "(devel)"?
              └─ yes → use it
              └─ no  → return "unknown"
  └─ no  → use version var
```

**Rationale**: `go install` does not support ldflags injection, but Go embeds the resolved module version in binary build metadata. This gives `go install @v0.0.1` users a correct version string without any extra work.

---

### Dev build format: `dev-YYYYMMddHHMMss`

**Decision**: Makefile defaults `VERSION` to `dev-$(shell date +%Y%m%d%H%M%S)`.

**Alternatives considered**:
- `dev-<git-sha>` — requires git to be available and the working tree to be clean; fragile in some environments.
- Static `"dev"` — provides no build identity, making it impossible to distinguish two dev builds.

**Rationale**: Timestamp is always available, cheap to generate, and makes dev binaries uniquely identifiable without git dependency.

---

### Archive format: `wip-<version>-<os>-<arch>.tar.gz`

**Decision**: Single binary inside a `.tar.gz`, plus a separate `checksums.txt` with SHA256 of all archives.

**Rationale**: Standard convention for Go CLI tools. Easy to verify integrity. Checksums file enables future Homebrew formula generation.

---

### Release notes: GitHub auto-generated

**Decision**: Use `gh release create --generate-notes`.

**Rationale**: Commit messages are high quality. Auto-generation is sufficient for an early-stage tool. Hand-written changelogs can be added later without changing the workflow structure.

## Risks / Trade-offs

- **`go build .` without make** → version shows `"unknown"`. Acceptable — documented behavior. Developers should use `make build`.
- **darwin/arm64 only** → limits early adopters to Apple Silicon Macs. Acceptable given current team. Expanding platforms is additive: one new Makefile target + one new workflow step per platform.
- **Tag-driven releases only** → no way to cut a release without a git tag. This is intentional — the tag is the version source of truth.

## Migration Plan

No migration required. This is purely additive:
- New file: `cmd/version.go`
- New file: `.github/workflows/release.yml`
- Modified: `main.go` (add `case "version"`)
- Modified: `Makefile` (add `VERSION`, `LDFLAGS`, `release-darwin-arm64`)

No existing behavior changes. Existing binaries built without ldflags continue to work; they will show `"unknown"` for version.
