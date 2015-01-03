package expression

import (
	"go/ast"
	"go/token"

	"github.com/zimmski/go-mutesting/mutator"
)

// MutatorRemoveTerm implements a mutator to remove expression terms
type MutatorRemoveTerm struct{}

// NewMutatorRemoveTerm returns a new instance of a MutatorRemoveTerm mutator
func NewMutatorRemoveTerm() *MutatorRemoveTerm {
	return &MutatorRemoveTerm{}
}

func init() {
	mutator.Register(MutatorRemoveTerm{}.String(), func() mutator.Mutator {
		return NewMutatorRemoveTerm()
	})
}

func (m *MutatorRemoveTerm) check(node ast.Node) (*ast.BinaryExpr, bool) {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}

	if n.Op != token.LAND && n.Op != token.LOR {
		return nil, false
	}

	return n, true
}

// Check validates how often a node can be mutated by a mutator
func (m *MutatorRemoveTerm) Check(node ast.Node) uint {
	_, ok := m.check(node)
	if !ok {
		return 0
	}

	return 2
}

// Mutate mutates a given node if it can be mutated by the mutator.
// It first checks if the given node can be mutated by the mutator. If the node cannot be mutated, false is send into the given control channel and the method returns. If the node can be mutated, the current state of the node is saved. Afterwards the node is mutated, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the original state of the node is restored, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the method returns which finishes the mutation process.
func (m *MutatorRemoveTerm) Mutate(node ast.Node, changed chan bool) {
	n, ok := m.check(node)
	if !ok {
		changed <- false

		return
	}

	var r *ast.Ident

	switch n.Op {
	case token.LAND:
		r = ast.NewIdent("true")
	case token.LOR:
		r = ast.NewIdent("false")
	}

	x := n.X
	y := n.Y

	n.X = r

	changed <- true
	<-changed

	n.X = x

	changed <- true
	<-changed

	n.Y = r

	changed <- true
	<-changed

	n.Y = y

	changed <- true
	<-changed
}

// String implements the String method of the Stringer interface
func (m MutatorRemoveTerm) String() string {
	return "expression/remove"
}
