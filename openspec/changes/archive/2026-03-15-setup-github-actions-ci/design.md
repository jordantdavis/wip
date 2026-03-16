## Context

The repo uses GitHub as its remote and has a `Makefile` with `build`, `test`, `fmt-check`, and `vet` targets. No CI exists today. GitHub Actions is the natural choice — it's native to GitHub, requires no external service, and is free for public repos.

## Goals / Non-Goals

**Goals:**
- Automated CI on every push to `main` and every PR targeting `main`
- Three parallel jobs: `build-test`, `fmt-check`, `vet`
- Each job surfaces as an independent status check on PRs
- Branch protection on `main` blocks merge when any check fails

**Non-Goals:**
- Caching (can be added later)
- Deployment or release automation
- Matrix builds across multiple Go versions or OSes

## Decisions

### Single workflow file with two triggers
Use one `.github/workflows/ci.yml` with both `push` (branches: main) and `pull_request` (branches: main) triggers rather than two separate files. The jobs are identical for both events; a single file avoids drift.

### Three jobs, not four
`build` and `test` are combined into one job (`build-test`) because `go test ./...` implicitly compiles all packages. Running `make build` first catches linker/binary issues (e.g. `main` package), then `make test` runs the full suite. `fmt-check` and `vet` are independent jobs with no dependency on each other or on `build-test`.

### Pin Go version from go.mod
The workflow pins Go `1.26.1` to match `go.mod`. This ensures CI matches the version used in development and avoids surprise breakage from runner image updates.

### Branch protection configured manually
Branch protection rules require a one-time manual setup in GitHub repo settings. The required check names (`build-test`, `fmt-check`, `vet`) must be added after the workflow has run at least once so they appear in the GitHub UI dropdown. This is documented as a post-deploy step.

## Risks / Trade-offs

- **Check names must match job names exactly** → Job `name:` fields in the workflow file are the source of truth for required check names in branch protection. If job names change, branch protection rules must be updated manually.
- **go 1.26.1 availability on runners** → `actions/setup-go@v6` will download the version if it's not pre-installed on the runner image. This adds a small cold-start cost but is reliable.
- **No required checks on direct pushes to main** → Branch protection blocks PR merges but not direct pushes by repo admins (unless "Do not allow bypassing" is enabled). The design notes this should be enabled.
