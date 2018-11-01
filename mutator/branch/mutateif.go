package branch

import (
	"go/ast"
	"go/types"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/if", MutatorIf)
}

// MutatorIf implements a mutator for if and else if branches.
func MutatorIf(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil
	}

	old := n.Body.List

	return []mutator.Mutation{
		{
			Change: func() {
				n.Body.List = []ast.Stmt{
					astutil.CreateNoopOfStatement(pkg, info, n.Body),
				}
			},
			Reset: func() {
				n.Body.List = old
			},
		},
	}
}
