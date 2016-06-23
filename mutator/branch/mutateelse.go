package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/astutil"
	"github.com/zimmski/go-mutesting/mutator"
)

// MutatorElse implements a mutator for else branches
type MutatorElse struct{}

// NewMutatorElse returns a new instance of a MutatorElse mutator
func NewMutatorElse() *MutatorElse {
	return &MutatorElse{}
}

func init() {
	mutator.Register(MutatorElse{}.String(), func() mutator.Mutator {
		return NewMutatorElse()
	})
}

func (m *MutatorElse) check(node ast.Node) (*ast.IfStmt, bool) {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil, false
	}

	// we ignore else ifs and nil blocks
	_, ok = n.Else.(*ast.IfStmt)
	if ok || n.Else == nil {
		return nil, false
	}

	return n, true
}

// Check validates how often a node can be mutated by a mutator
func (m *MutatorElse) Check(node ast.Node) uint {
	_, ok := m.check(node)
	if !ok {
		return 0
	}

	return 1
}

// Mutate mutates a given node if it can be mutated by the mutator.
// It first checks if the given node can be mutated by the mutator. If the node cannot be mutated, false is send into the given control channel and the method returns. If the node can be mutated, the current state of the node is saved. Afterwards the node is mutated, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the original state of the node is restored, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the method returns which finishes the mutation process.
func (m *MutatorElse) Mutate(node ast.Node, changed chan bool) {
	n, ok := m.check(node)
	if !ok {
		changed <- false

		return
	}

	old := n.Else
	n.Else = astutil.CreateNoopOfStatement(old)

	changed <- true
	<-changed

	n.Else = old

	changed <- true
	<-changed
}

// String implements the String method of the Stringer interface
func (m MutatorElse) String() string {
	return "branch/else"
}
