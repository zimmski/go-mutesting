package statement

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

func init() {
	mutator.Register(MutatorRemoveStatement{}.String(), func() mutator.Mutator {
		return NewMutatorRemoveStatement()
	})
}

// MutatorRemoveStatement implements a mutator to remove statements
type MutatorRemoveStatement struct{}

// NewMutatorRemoveStatement returns a new instance of a MutatorRemoveStatement mutator
func NewMutatorRemoveStatement() *MutatorRemoveStatement {
	return &MutatorRemoveStatement{}
}

// String implements the String method of the Stringer interface
func (m MutatorRemoveStatement) String() string {
	return "statement/remove"
}

func (m *MutatorRemoveStatement) checkStatement(node ast.Stmt) bool {
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

// Mutations returns a list of possible mutations for the given node.
func (m *MutatorRemoveStatement) Mutations(node ast.Node) []mutator.Mutation {
	var l []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	}

	var mutations []mutator.Mutation

	for i, ni := range l {
		if m.checkStatement(ni) {
			li := i
			old := l[li]

			mutations = append(mutations, mutator.Mutation{
				Change: func() {
					l[li] = astutil.CreateNoopOfStatement(old)
				},
				Reset: func() {
					l[li] = old
				},
			})
		}
	}

	return mutations
}
