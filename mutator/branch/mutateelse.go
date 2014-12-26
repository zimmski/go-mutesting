package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/mutator"
)

type MutatorElse struct{}

func NewMutatorElse() *MutatorElse {
	return &MutatorElse{}
}

func init() {
	mutator.Register("branch/else", func() mutator.Mutator {
		return NewMutatorElse()
	})
}

func (m *MutatorElse) check(node ast.Node) (*ast.IfStmt, bool) {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil, false
	}

	// if it is an else if or the block is nil
	_, ok = n.Else.(*ast.IfStmt)
	if ok || n.Else == nil {
		return nil, false
	}

	return n, true
}

func (m *MutatorElse) Check(node ast.Node) bool {
	_, ok := m.check(node)

	return ok
}

func (m *MutatorElse) Mutate(node ast.Node, changed chan bool) {
	n, ok := m.check(node)
	if !ok {
		changed <- false

		return
	}

	old := n.Else
	n.Else = &ast.EmptyStmt{
		Semicolon: n.Else.Pos(),
	}

	changed <- true
	<-changed

	n.Else = old

	changed <- true
	<-changed
}
