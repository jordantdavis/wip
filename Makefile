.PHONY: build release-darwin-arm64 test fmt fmt-check vet install clean check

VERSION ?= dev-$(shell date +%Y%m%d%H%M%S)
LDFLAGS = -X github.com/jordantdavis/wip/cmd.version=$(VERSION)

_no_default:
	@echo "No default target. Specify one of: build test fmt fmt-check vet install clean check"; exit 1

build:
	go build -ldflags "$(LDFLAGS)" -o wip .

release-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o wip-darwin-arm64 .

test:
	go test ./...

fmt:
	go fmt ./...

fmt-check:
	@files=$$(gofmt -l .); if [ -n "$$files" ]; then echo "Unformatted files:"; echo "$$files"; exit 1; fi

vet:
	go vet ./...

install:
	go install .

clean:
	rm -f wip

check: fmt-check vet test
