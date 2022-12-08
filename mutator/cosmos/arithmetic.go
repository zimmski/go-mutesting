package cosmos

import (
	"go/ast"
	"go/types"

	"github.com/osmosis-labs/go-mutesting/mutator"
)

func init() {
	mutator.Register("cosmos/arithmetic", MutatorArithmeticCosmos)
}

var arithmeticMutations = map[string]string{
	"Add": "Sub",
	"Sub": "Add",
	"Mul": "Quo",
	"Quo": "Mul",
}

// MutatorArithmeticCosmos implements a mutator to change Cosmos SDK arithmetic.
func MutatorArithmeticCosmos(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	// ensure node has a valid SDK arithmetic operator
	if _, ok := arithmeticMutations[n.Name]; !ok {
		return nil
	}

	o := n.Name
	r, ok := arithmeticMutations[n.Name]
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
