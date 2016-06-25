package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register(MutatorElse{}.String(), func() mutator.Mutator {
		return NewMutatorElse()
	})
}

// MutatorElse implements a mutator for else branches
type MutatorElse struct{}

// NewMutatorElse returns a new instance of a MutatorElse mutator
func NewMutatorElse() *MutatorElse {
	return &MutatorElse{}
}

// String implements the String method of the Stringer interface
func (m MutatorElse) String() string {
	return "branch/else"
}

// Mutations returns a list of possible mutations for the given node.
func (m *MutatorElse) Mutations(node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil
	}
	// we ignore else ifs and nil blocks
	_, ok = n.Else.(*ast.IfStmt)
	if ok || n.Else == nil {
		return nil
	}

	old := n.Else

	return []mutator.Mutation{
		mutator.Mutation{
			Change: func() {
				n.Else = astutil.CreateNoopOfStatement(old)
			},
			Reset: func() {
				n.Else = old
			},
		},
	}
}
