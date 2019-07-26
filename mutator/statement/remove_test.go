package statement

import (
	"testing"

	"github.com/mccurdyc/go-mutesting/test"
)

func TestMutatorRemoveStatement(t *testing.T) {
	test.Mutator(
		t,
		MutatorRemoveStatement,
		"../../testdata/statement/remove.go",
		17,
	)
}
