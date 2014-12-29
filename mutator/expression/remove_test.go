package expression

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/go-mutesting/test"
)

func TestMutateElse(t *testing.T) {
	m := NewMutatorRemoveTerm()
	NotNil(t, m)

	test.Mutator(
		t,
		"../../testdata/expression/remove.go",
		m,
		6,
	)
}
