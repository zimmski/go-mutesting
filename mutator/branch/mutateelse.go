package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register("branch/else", MutatorElse)
}

// MutatorElse implements a mutator for else branches.
func MutatorElse(node ast.Node) []mutator.Mutation {
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
