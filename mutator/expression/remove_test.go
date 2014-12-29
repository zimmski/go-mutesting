package expression

import (
	"testing"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutateElse(t *testing.T) {
	test.Mutator(
		t,
		NewMutatorRemoveTerm(),
		"../../testdata/expression/remove.go",
		6,
	)
}
