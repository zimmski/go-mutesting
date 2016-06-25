package mutator

import (
	"fmt"
	"go/ast"
	"sort"
)

// Mutator defines a mutator for mutation testing
type Mutator interface {
	fmt.Stringer

	// Mutations returns a list of possible mutations for the given node.
	Mutations(node ast.Node) []Mutation
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
