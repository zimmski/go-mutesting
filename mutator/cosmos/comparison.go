package cosmos

import (
	"go/ast"
	"go/types"

	"github.com/osmosis-labs/go-mutesting/mutator"
)

func init() {
	mutator.Register("cosmos/comparison", MutatorComparisonCosmos)
}

var comparisonMutations = map[string]string{
	"GT":  "LTE",
	"LT":  "GTE",
	"GTE": "LT",
	"LTE": "GT",
}

// MutatorComparisonCosmos implements a mutator to change Cosmos SDK comparisons.
func MutatorComparisonCosmos(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	// ensure node has a valid SDK comparison operator
	if _, ok := comparisonMutations[n.Name]; !ok {
		return nil
	}

	o := n.Name
	r, ok := comparisonMutations[n.Name]
	if !ok {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Name = r
			},
			Reset: func() {
				n.Name = o
			},
		},
	}
}
