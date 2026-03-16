## 1. Workflow File

- [x] 1.1 Create `.github/workflows/ci.yml` with `push` (branches: main) and `pull_request` (branches: main) triggers
- [x] 1.2 Add `build-test` job: checkout@v5, setup-go@v6 (go 1.26.1), `make build`, `make test`
- [x] 1.3 Add `fmt-check` job: checkout@v5, setup-go@v6 (go 1.26.1), `make fmt-check`
- [x] 1.4 Add `vet` job: checkout@v5, setup-go@v6 (go 1.26.1), `make vet`

## 2. GitHub Repo Configuration

- [x] 2.1 Push workflow file to `main` (or merge a PR) so the jobs run once and status check names register with GitHub
- [x] 2.2 In Settings → Branches, add branch protection rule for `main` with required status checks: `build-test`, `fmt-check`, `vet`
- [x] 2.3 Enable "Require branches to be up to date before merging" and "Do not allow bypassing the above settings"
