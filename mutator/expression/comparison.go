package expression

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("expression/comparison", MutatorComparison)
}

var comparisonMutations = map[token.Token]token.Token{
	token.LSS: token.LEQ,
	token.LEQ: token.LSS,
	token.GTR: token.GEQ,
	token.GEQ: token.GTR,
}

// MutatorComparison implements a mutator to change comparisons.
func MutatorComparison(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	o := n.Op
	r, ok := comparisonMutations[n.Op]
	if !ok {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Op = r
			},
			Reset: func() {
				n.Op = o
			},
		},
	}
}
