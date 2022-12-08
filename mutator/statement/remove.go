package statement

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	mutesting "github.com/osmosis-labs/go-mutesting"
	"github.com/osmosis-labs/go-mutesting/astutil"
	"github.com/osmosis-labs/go-mutesting/mutator"
)

func init() {
	mutator.Register("statement/remove", MutatorRemoveStatement)
}

func checkRemoveStatement(node ast.Stmt) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		if n.Tok != token.DEFINE {
			return true
		}
	case *ast.ExprStmt, *ast.IncDecStmt:
		return true
	}

	return false
}

// MutatorRemoveStatement implements a mutator to remove statements.
func MutatorRemoveStatement(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
	var l []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	}

	// If the statement block only has one item in it (i.e. AST sub-tree is only 1 level deep),
	// we check to see if that statement is a panic.
	//
	// Since we only target AST leaves here, all isolated panics that are nested in conditionals
	// are still included in mutations in earlier conditional mutations (which run before statement
	// mutations), and only those that are nested in error checks are filtered (this is due to an
	// AST quirk where it does not treat error check panics as conditionals and instead directly 
	// labels them as statements).
	if len(l) == 1 {
		if containsPanic := strings.Contains(mutesting.GetNodeASTString(node), "panic"); containsPanic {
			return nil
		}
	}

	var mutations []mutator.Mutation

	for i, ni := range l {
		if checkRemoveStatement(ni) {
			li := i
			old := l[li]

			mutations = append(mutations, mutator.Mutation{
				Change: func() {
					l[li] = astutil.CreateNoopOfStatement(pkg, info, old)
				},
				Reset: func() {
					l[li] = old
				},
			})
		}
	}

	return mutations
}
