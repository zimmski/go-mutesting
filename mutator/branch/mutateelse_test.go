package branch

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutateElse(t *testing.T) {
	m := NewMutatorElse()
	NotNil(t, m)

	test.Mutator(
		t,
		"../../testdata/branch/mutateelse.go",
		m,
		1,
	)
}
