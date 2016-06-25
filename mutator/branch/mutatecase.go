package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/case", MutatorCase)
}

// MutatorCase implements a mutator for case clauses.
func MutatorCase(node ast.Node) []mutator.Mutation {
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
