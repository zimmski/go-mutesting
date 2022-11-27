package branch

import (
	"fmt"
	"go/ast"
	"go/types"

	"strings"

	"github.com/osmosis-labs/go-mutesting/astutil"
	"github.com/osmosis-labs/go-mutesting/mutator"
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
	
	// We filter conditionals that take the form `err != nil` { return }, since
	// these mutations are almost always false positives and make up a significant
	// portion of the noise in mutation testing results.
	//
	// We would also like to filter such statements when they have a nested panic
	// instead of a return, but those mutations are treated as statements, not conditionals,
	// so we include separate logic to handle them in the statement mutator.
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
