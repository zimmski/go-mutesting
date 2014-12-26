package mutator

import (
	"fmt"
	"go/ast"
	"sort"
)

type Mutator interface {
	Check(node ast.Node) bool
	Mutate(node ast.Node, changed chan bool)
}

var mutatorLookup = make(map[string]func() Mutator)

func New(name string) (Mutator, error) {
	mutatorFunc, ok := mutatorLookup[name]
	if !ok {
		return nil, fmt.Errorf("unknown mutator %q", name)
	}

	return mutatorFunc(), nil
}

func List() []string {
	keyMutatorLookup := make([]string, 0, len(mutatorLookup))

	for key := range mutatorLookup {
		keyMutatorLookup = append(keyMutatorLookup, key)
	}

	sort.Strings(keyMutatorLookup)

	return keyMutatorLookup
}

func Register(name string, mutatorFunc func() Mutator) {
	if mutatorFunc == nil {
		panic("mutator function is nil")
	}

	if _, ok := mutatorLookup[name]; ok {
		panic(fmt.Sprintf("mutator %q already registered", name))
	}

	mutatorLookup[name] = mutatorFunc
}
