package cosmos

import (
	"go/ast"
	"go/types"

	"github.com/osmosis-labs/go-mutesting/mutator"
)

func init() {
	mutator.Register("cosmos/getters", MutatorGetterCosmos)
}

var getterMutations = map[string]string{
	"GetToken0": "GetToken1",
	"GetToken1": "GetToken0",
}

// MutatorGetterCosmos implements a mutator to change Cosmos SDK getters.
func MutatorGetterCosmos(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.Ident)
	if !ok {
		return nil
	}

	// ensure node has a valid SDK comparison operator
	if _, ok := getterMutations[n.Name]; !ok {
		return nil
	}

	o := n.Name
	r, ok := getterMutations[n.Name]
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
