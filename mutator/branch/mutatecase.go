package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/mutator"
)

// MutatorCase implements a mutator for case
type MutatorCase struct{}

// NewMutatorCase returns a new instance of a MutatorCase mutator
func NewMutatorCase() *MutatorCase {
	return &MutatorCase{}
}

func init() {
	mutator.Register(MutatorCase{}.String(), func() mutator.Mutator {
		return NewMutatorCase()
	})
}

func (m *MutatorCase) check(node ast.Node) (*ast.CaseClause, bool) {
	n, ok := node.(*ast.CaseClause)

	return n, ok
}

// Check validates how often a node can be mutated by a mutator
func (m *MutatorCase) Check(node ast.Node) uint {
	_, ok := m.check(node)
	if !ok {
		return 0
	}

	return 1
}

// Mutate mutates a given node if it can be mutated by the mutator.
// It first checks if the given node can be mutated by the mutator. If the node cannot be mutated, false is send into the given control channel and the method returns. If the node can be mutated, the current state of the node is saved. Afterwards the node is mutated, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the original state of the node is restored, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the method returns which finishes the mutation process.
func (m *MutatorCase) Mutate(node ast.Node, changed chan bool) {
	n, ok := m.check(node)
	if !ok {
		changed <- false

		return
	}

	old := n.Body
	n.Body = make([]ast.Stmt, 0)

	changed <- true
	<-changed

	n.Body = old

	changed <- true
	<-changed
}

// String implements the String method of the Stringer interface
func (m MutatorCase) String() string {
	return "branch/case"
}
