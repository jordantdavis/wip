## 1. Create Makefile

- [x] 1.1 Add `build` target: `go build -o wip .`
- [x] 1.2 Add `test` target: `go test ./...`
- [x] 1.3 Add `fmt` target: `go fmt ./...`
- [x] 1.4 Add `fmt-check` target: fail with file list if `gofmt -l .` output is non-empty
- [x] 1.5 Add `vet` target: `go vet ./...`
- [x] 1.6 Add `install` target: `go install .`
- [x] 1.7 Add `clean` target: `rm -f wip`
- [x] 1.8 Add `check` target: depends on `fmt-check vet test`
- [x] 1.9 Declare all targets as `.PHONY`
- [x] 1.10 Verify bare `make` with no arguments produces an error (no default target)
