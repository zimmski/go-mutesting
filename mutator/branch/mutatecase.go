package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register(MutatorCase{}.String(), func() mutator.Mutator {
		return NewMutatorCase()
	})
}

// MutatorCase implements a mutator for case
type MutatorCase struct{}

// NewMutatorCase returns a new instance of a MutatorCase mutator
func NewMutatorCase() *MutatorCase {
	return &MutatorCase{}
}

// String implements the String method of the Stringer interface
func (m MutatorCase) String() string {
	return "branch/case"
}

// Mutations returns a list of possible mutations for the given node.
func (m *MutatorCase) Mutations(node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.CaseClause)
	if !ok {
		return nil
	}

	old := n.Body

	return []mutator.Mutation{
		mutator.Mutation{
			Change: func() {
				n.Body = []ast.Stmt{
					astutil.CreateNoopOfStatements(n.Body),
				}
			},
			Reset: func() {
				n.Body = old
			},
		},
	}
}
