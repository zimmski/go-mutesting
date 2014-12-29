package branch

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutateIf(t *testing.T) {
	m := NewMutatorIf()
	NotNil(t, m)

	test.Mutator(
		t,
		"../../testdata/branch/mutateif.go",
		m,
		2,
	)
}
