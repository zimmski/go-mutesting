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

all: install-tools install-dependencies install lint test

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
generate: clean
	go generate $(PKG)/...
install:
	go install -v $(PKG)/...
install-dependencies:
	go get -t -v $(PKG)/...
	go build -v $(PKG)/...
install-tools:
	# generation
	go get -u -v golang.org/x/tools/cmd/stringer

	# linting
	go get -u -v github.com/golang/lint/...
	go get -u -v github.com/kisielk/errcheck/...

	# code coverage
	go get -u -v golang.org/x/tools/cmd/cover
	go get -u -v github.com/onsi/ginkgo/ginkgo/...
	go get -u -v github.com/modocache/gover/...
	go get -u -v github.com/mattn/goveralls/...
lint:
	$(ROOT_DIR)/scripts/lint.sh
test:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" $(PKG_TEST)
test-verbose:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v $(PKG_TEST)
test-verbose-with-coverage:
	ginkgo -r -v -cover -race -skipPackage="testdata"
