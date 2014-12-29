package branch

import (
	"testing"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutatorIf(t *testing.T) {
	test.Mutator(
		t,
		NewMutatorIf(),
		"../../testdata/branch/mutateif.go",
		2,
	)
}
