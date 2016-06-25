package expression

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register(MutatorRemoveTerm{}.String(), func() mutator.Mutator {
		return NewMutatorRemoveTerm()
	})
}

// MutatorRemoveTerm implements a mutator to remove expression terms
type MutatorRemoveTerm struct{}

// NewMutatorRemoveTerm returns a new instance of a MutatorRemoveTerm mutator
func NewMutatorRemoveTerm() *MutatorRemoveTerm {
	return &MutatorRemoveTerm{}
}

// String implements the String method of the Stringer interface
func (m MutatorRemoveTerm) String() string {
	return "expression/remove"
}

// Mutations returns a list of possible mutations for the given node.
func (m *MutatorRemoveTerm) Mutations(node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}
	if n.Op != token.LAND && n.Op != token.LOR {
		return nil
	}

	var r *ast.Ident

	switch n.Op {
	case token.LAND:
		r = ast.NewIdent("true")
	case token.LOR:
		r = ast.NewIdent("false")
	}

	x := n.X
	y := n.Y

	return []mutator.Mutation{
		mutator.Mutation{
			Change: func() {
				n.X = r
			},
			Reset: func() {
				n.X = x
			},
		},
		mutator.Mutation{
			Change: func() {
				n.Y = r
			},
			Reset: func() {
				n.Y = y
			},
		},
	}
}
