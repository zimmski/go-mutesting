package statement

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/mutator"
)

// MutatorRemoveStatement implements a mutator to remove statements
type MutatorRemoveStatement struct{}

// NewMutatorRemoveStatement returns a new instance of a MutatorRemoveStatement mutator
func NewMutatorRemoveStatement() *MutatorRemoveStatement {
	return &MutatorRemoveStatement{}
}

func init() {
	mutator.Register(MutatorRemoveStatement{}.String(), func() mutator.Mutator {
		return NewMutatorRemoveStatement()
	})
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

func (m *MutatorRemoveStatement) check(node ast.Node) uint {
	var count uint
	var l []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	default:
		return 0
	}

	for _, ni := range l {
		if m.checkStatement(ni) {
			count++
		}
	}

	return count
}

// Check validates how often a node can be mutated by a mutator
func (m *MutatorRemoveStatement) Check(node ast.Node) uint {
	count := m.check(node)

	return count
}

// Mutate mutates a given node if it can be mutated by the mutator.
// It first checks if the given node can be mutated by the mutator. If the node cannot be mutated, false is send into the given control channel and the method returns. If the node can be mutated, the current state of the node is saved. Afterwards the node is mutated, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the original state of the node is restored, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the method returns which finishes the mutation process.
func (m *MutatorRemoveStatement) Mutate(node ast.Node, changed chan bool) {
	count := m.check(node)
	if count == 0 {
		changed <- false

		return
	}

	var l []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	}

	for i, ni := range l {
		if m.checkStatement(ni) {
			old := l[i]
			l[i] = &ast.EmptyStmt{
				Semicolon: old.Pos(),
			}

			changed <- true
			<-changed

			l[i] = old

			changed <- true
			<-changed
		}
	}
}

// String implements the String method of the Stringer interface
func (m MutatorRemoveStatement) String() string {
	return "statement/remove"
}
