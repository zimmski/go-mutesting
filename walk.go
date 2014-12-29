package mutesting

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/zimmski/go-mutesting/mutator"
)

// CountWalk returns the number of corresponding mutations for a given mutator.
// It traverses the AST of the given node and calls the method Check of the given mutator for every node and sums up the returned counts. After completion of the traversal the final counter is returned.
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

// Visit implements the Visit method of the ast.Visitor interface
func (w *countWalk) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}

	if n := w.mutator.Check(node); n > 0 {
		w.count += n
	}

	return w
}

// MutateWalk mutates the given node with the given mutator returning a channel to control the mutation steps.
// It traverses the AST of the given node and calls the method Check of the given mutator to verify that a node can be mutated by the mutator. If a node can be mutated the method Mutate of the given mutator is executed with the node and the control channel. After completion of the traversal the control channel is closed.
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

// Visit implements the Visit method of the ast.Visitor interface
func (w *mutateWalk) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}

	if w.mutator.Check(node) > 0 {
		w.mutator.Mutate(node, w.changed)
	}

	return w
}

// PrintWalk traverses the AST of the given node and prints every node to STDOUT.
func PrintWalk(node ast.Node) {
	w := &printWalk{
		level: 0,
	}

	ast.Walk(w, node)
}

type printWalk struct {
	level int
}

// Visit implements the Visit method of the ast.Visitor interface
func (w *printWalk) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		w.level++

		fmt.Printf("%s(%p)%#v\n", strings.Repeat("\t", w.level), node, node)
	} else {
		w.level--
	}

	return w
}
