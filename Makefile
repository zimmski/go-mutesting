.PHONY: all clean clean-coverage generate install install-dependencies install-tools lint test test-verbose test-verbose-with-coverage

export ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
export PKG := github.com/zimmski/go-mutesting
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

export TEST_TIMEOUT_IN_SECONDS := 240

$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

ifdef ARGS
	PKG_TEST := $(ARGS)
else
	PKG_TEST := $(PKG)/...
endif

all: install-dependencies install-tools install lint test
.PHONY: all

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
.PHONY: clean

clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
.PHONY: clean-coverage

generate: clean
	go generate $(PKG)/...
.PHONY: generate

install:
	go install -v $(PKG)/...
.PHONY: install

install-dependencies:
	go get
.PHONY: install-dependencies

install-tools:
	# generation
	go install golang.org/x/tools/cmd/stringer

	# linting
	go install golang.org/x/lint/golint/...
	go install github.com/kisielk/errcheck/...
	go install honnef.co/go/tools/...

	# code coverage
	go install golang.org/x/tools/cmd/cover
	go install github.com/onsi/ginkgo/ginkgo/...
	go install github.com/modocache/gover/...
	go install github.com/mattn/goveralls/...
.PHONY: install-tools

lint:
	$(ROOT_DIR)/scripts/lint.sh
.PHONY: lint

test:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" $(PKG_TEST)
.PHONY: test

test-verbose:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v $(PKG_TEST)
.PHONY: test-verbose

test-verbose-with-coverage:
	ginkgo -r -v -cover -race -skipPackage="testdata"
.PHONY: test-verbose-with-coverage
