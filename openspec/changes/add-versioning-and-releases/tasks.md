## 1. Version Command

- [x] 1.1 Create `cmd/version.go` with `var version = "unknown"`, `getVersion()` (ldflags → build info → "unknown"), and `Version()` printing `wip <version> (<os>/<arch>)`
- [x] 1.2 Add `case "version": cmd.Version()` to the switch in `main.go`
- [x] 1.3 Verify `wip version` output format matches `wip <version> (darwin/arm64)`

## 2. Makefile

- [x] 2.1 Add `VERSION ?= dev-$(shell date +%Y%m%d%H%M%S)` and `LDFLAGS = -X github.com/jordantdavis/wip/cmd.version=$(VERSION)` variables
- [x] 2.2 Update `build` target to pass `-ldflags "$(LDFLAGS)"` to `go build`
- [x] 2.3 Add `release-darwin-arm64` target: `GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o wip-darwin-arm64 .`
- [x] 2.4 Verify `make build` produces a binary where `wip version` shows `dev-<timestamp>`
- [x] 2.5 Verify `make build VERSION=v0.0.1` produces a binary where `wip version` shows `v0.0.1`

## 3. Release Workflow

- [x] 3.1 Create `.github/workflows/release.yml` triggered on `push: tags: ['v*']`
- [x] 3.2 Add build step: `make release-darwin-arm64 VERSION=${{ github.ref_name }}`
- [x] 3.3 Add archive step: `tar -czf wip-${{ github.ref_name }}-darwin-arm64.tar.gz wip-darwin-arm64`
- [x] 3.4 Add checksum step: `sha256sum wip-*.tar.gz > checksums.txt`
- [x] 3.5 Add publish step: `gh release create ${{ github.ref_name }} --generate-notes wip-*.tar.gz checksums.txt`
- [ ] 3.6 Tag `v0.0.1` and push to verify the full release workflow end-to-end
