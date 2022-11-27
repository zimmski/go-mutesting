export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export PKG := github.com/zimmski/osutil

export UNIT_TEST_TIMEOUT := 480

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

all: install-dependencies install-tools install lint test-verbose
.PHONY: all

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
.PHONY: clean

clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
.PHONY: clean-coverage

install:
	go install -v ./...
.PHONY: install

install-dependencies:
	go get -t -v ./...
	go test -i -v ./...
.PHONY: install-dependencies

install-tools:
	# Install linting tools
	go get -u -v golang.org/x/lint/golint/...
	go get -u -v github.com/kisielk/errcheck/...
	go get -u -v honnef.co/go/tools/cmd/megacheck

	# Install code coverage tools
	go get -u -v github.com/onsi/ginkgo/ginkgo/...
	go get -u -v github.com/modocache/gover/...
	go get -u -v github.com/mattn/goveralls/...
.PHONY: install-tools

lint:
	$(ROOT_DIR)/scripts/lint.sh
.PHONY: lint

test:
	go test -race -test.timeout $(UNIT_TEST_TIMEOUT)s $(PKG_TEST)
.PHONY: test

test-with-coverage:
	ginkgo -r -cover -race -skipPackage="testdata" $(PKG_TEST)
.PHONY: test-with-coverage

test-verbose:
	go test -race -test.timeout $(UNIT_TEST_TIMEOUT)s -v $(PKG_TEST)
.PHONY: test-verbose
