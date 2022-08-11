package branch

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

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
	
	containsErr := strings.Contains(fmt.Sprintf("%v", n.Cond), "err")
	containsNilCheck := strings.Contains(fmt.Sprintf("%v", n.Cond), "!= nil")
	
	if containsErr && containsNilCheck && len(n.Body.List) == 1 {
		return nil
	}
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
