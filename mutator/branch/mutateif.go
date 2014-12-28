package branch

import (
	"go/ast"

	"github.com/zimmski/go-mutesting/mutator"
)

type MutatorIf struct{}

func NewMutatorIf() *MutatorIf {
	return &MutatorIf{}
}

func init() {
	mutator.Register(MutatorIf{}.String(), func() mutator.Mutator {
		return NewMutatorIf()
	})
}

func (m *MutatorIf) check(node ast.Node) (*ast.IfStmt, bool) {
	n, ok := node.(*ast.IfStmt)

	return n, ok
}

func (m *MutatorIf) Check(node ast.Node) bool {
	_, ok := m.check(node)

	return ok
}

func (m *MutatorIf) Mutate(node ast.Node, changed chan bool) {
	n, ok := m.check(node)
	if !ok {
		changed <- false

		return
	}

	old := n.Body.List
	n.Body.List = make([]ast.Stmt, 0)

	changed <- true
	<-changed

	n.Body.List = old

	changed <- true
	<-changed
}

func (m MutatorIf) String() string {
	return "branch/if"
}
