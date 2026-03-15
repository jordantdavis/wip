.PHONY: build test fmt fmt-check vet install clean check

_no_default:
	@echo "No default target. Specify one of: build test fmt fmt-check vet install clean check"; exit 1

build:
	go build -o wip .

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
