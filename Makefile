.PHONY: all clean clean-coverage generate install install-dependencies install-markdown install-tools lint markdown test test-verbose test-with-coverage

export ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
export PKG := github.com/zimmski/go-mutesting
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

all: install-tools install-dependencies install lint test

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
generate: clean
	go generate $(PKG)/...
install: generate
	go install -v $(PKG)/...
install-dependencies:
	go get -t -v $(PKG)/...
	go build -v $(PKG)/...
install-markdown:
	go get -u -v github.com/noraesae/orange-cat/...
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
markdown:
	orange
test:
	go test -race -test.timeout 60s $(PKG)/...
test-verbose:
	go test -race -test.timeout 60s -v $(PKG)/...
test-with-coverage:
	ginkgo -r -cover -race -skipPackage="testdata"
