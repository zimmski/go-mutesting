package mutator

import (
	"fmt"
	"go/ast"
	"sort"
)

// Mutator defines a mutator for mutation testing
type Mutator interface {
	// Check validates if a given node can be mutated by the mutator
	Check(node ast.Node) bool
	// Mutate mutates a given node if it can be mutated by the mutator.
	// It first checks if the given node can be mutated by the mutator. If the node cannot be mutated, false is send into the given control channel and the method returns. If the node can be mutated, the current state of the node is saved. Afterwards the node is mutated, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the original state of the node is restored, true is send into the given control channel and the method waits on the channel to continue the process. After receiving a value from the channel the method returns which finishes the mutation process.
	Mutate(node ast.Node, changed chan bool)
}

var mutatorLookup = make(map[string]func() Mutator)

// New returns a new mutator instance given the registered name of the mutator.
// The error return argument is not nil, if the name does not exist in the registered mutator list.
func New(name string) (Mutator, error) {
	mutatorFunc, ok := mutatorLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown mutator %q", name)
	}

	return mutatorFunc(), nil
}

// List returns a list of all registered mutator names.
func List() []string {
	keyMutatorLookup := make([]string, 0, len(mutatorLookup))

	for key := range mutatorLookup {
		keyMutatorLookup = append(keyMutatorLookup, key)
	}

	sort.Strings(keyMutatorLookup)

	return keyMutatorLookup
}

// Register registers a mutator instance function with the given name.
func Register(name string, mutatorFunc func() Mutator) {
	if mutatorFunc == nil {
		panic("mutator function is nil")
	}

	if _, ok := mutatorLookup[name]; ok {
		panic(fmt.Sprintf("mutator %q already registered", name))
	}

	mutatorLookup[name] = mutatorFunc
}
