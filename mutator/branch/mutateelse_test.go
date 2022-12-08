package branch

import (
	"testing"

	"github.com/osmosis-labs/go-mutesting/test"
)

func TestMutatorElse(t *testing.T) {
	test.Mutator(
		t,
		MutatorElse,
		"../../testdata/branch/mutateelse.go",
		1,
	)
}
