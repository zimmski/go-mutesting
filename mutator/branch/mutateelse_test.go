package branch

import (
	"testing"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutatorElse(t *testing.T) {
	test.Mutator(
		t,
		NewMutatorElse(),
		"../../testdata/branch/mutateelse.go",
		1,
	)
}
