package expression

import (
	"testing"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutatorComparison(t *testing.T) {
	test.Mutator(
		t,
		MutatorComparison,
		"../../testdata/expression/comparison.go",
		4,
	)
}
