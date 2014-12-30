.PHONY: all clean coverage debug-install dependencies fmt install lint markdown test testverbose tools

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: tools dependencies install test

clean:
	go clean -i ./...
	go clean -i -race ./...
coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
crosscompile:
	gox -os="linux" ./...
debug-install: clean
	go install -race -v ./...
dependencies:
	go get -t -v ./...
	go build -v ./...
fmt:
	gofmt -l -w $(ROOT_DIR)/
install: clean
	go install -v ./...
	go generate ./...
	go install -v ./...
lint: install fmt
	errcheck github.com/zimmski/go-mutesting/... || true
	golint ./... | grep --invert-match -P "(_string.go:)" || true
	go tool vet -all=true -v=true $(ROOT_DIR)/ 2>&1 | grep --invert-match -P "(Checking file|\%p of wrong type|can't check non-constant format)" || true
markdown:
	orange
test:
	go test -race ./...
testverbose:
	go test -race -v ./...
tools:
	# generation
	go get -u golang.org/x/tools/cmd/godoc
	go get -u golang.org/x/tools/cmd/stringer

	# linting
	go get -u golang.org/x/tools/cmd/vet
	go get -u github.com/golang/lint
	go install github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

	# code coverage
	go get -u golang.org/x/tools/cmd/cover
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
