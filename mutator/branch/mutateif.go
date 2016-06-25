package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register(MutatorIf{}.String(), func() mutator.Mutator {
		return NewMutatorIf()
	})
}

// NewMutatorIf returns a new instance of a MutatorIf mutator
func NewMutatorIf() *MutatorIf {
	return &MutatorIf{}
}

// MutatorIf implements a mutator for if and else if branches
type MutatorIf struct{}

// String implements the String method of the Stringer interface
func (m MutatorIf) String() string {
	return "branch/if"
}

// Mutations returns a list of possible mutations for the given node.
func (m *MutatorIf) Mutations(node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil
	}

	old := n.Body.List

	return []mutator.Mutation{
		mutator.Mutation{
			Change: func() {
				n.Body.List = []ast.Stmt{
					astutil.CreateNoopOfStatement(n.Body),
				}
			},
			Reset: func() {
				n.Body.List = old
			},
		},
	}
}
