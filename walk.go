package mutesting

import (
	"fmt"
	"go/ast"

	"github.com/zimmski/go-mutesting/mutator"
)

func CountWalk(node ast.Node, m mutator.Mutator) uint {
	w := &countWalk{
		count:   0,
		mutator: m,
	}

	ast.Walk(w, node)

	return w.count
}

type countWalk struct {
	count   uint
	mutator mutator.Mutator
}

func (w *countWalk) Visit(node ast.Node) ast.Visitor {
	if w.mutator.Check(node) {
		w.count++
	}

	return w
}

func MutateWalk(node ast.Node, m mutator.Mutator) chan bool {
	w := &mutateWalk{
		changed: make(chan bool),
		mutator: m,
	}

	go func() {
		ast.Walk(w, node)

		close(w.changed)
	}()

	return w.changed
}

type mutateWalk struct {
	changed chan bool
	mutator mutator.Mutator
}

func (w *mutateWalk) Visit(node ast.Node) ast.Visitor {
	if w.mutator.Check(node) {
		w.mutator.Mutate(node, w.changed)
	}

	return w
}

func PrintWalk(node ast.Node) {
	w := &printWalk{}

	ast.Walk(w, node)
}

type printWalk struct{}

func (w *printWalk) Visit(node ast.Node) ast.Visitor {
	fmt.Printf("%#v\n", node)

	return w
}
