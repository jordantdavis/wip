## Why

The repo has no automated CI. Without it, broken builds, unformatted code, and failing tests can land on `main` undetected and PRs have no automated quality gate before merge.

## What Changes

- Add a GitHub Actions workflow that runs on every push to `main` and every PR targeting `main`
- Run `make build` + `make test` as a combined job
- Run `make fmt-check` as a standalone job
- Run `make vet` as a standalone job
- All three jobs run in parallel; each surfaces as an independent status check on PRs
- Branch protection on `main` will require all three checks to pass before merge is allowed

## Capabilities

### New Capabilities

- `github-actions-ci`: GitHub Actions workflow providing build, format, vet, and test CI checks on push to main and on PRs targeting main

### Modified Capabilities

<!-- none -->

## Impact

- New file: `.github/workflows/ci.yml`
- No changes to application code, tests, or the Makefile
- Requires branch protection rule configuration in GitHub repo settings (manual step, documented in design)
