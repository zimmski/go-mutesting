package mutator

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockMutator(pkg *types.Package, info *types.Info, node ast.Node) []Mutation {
	// Do nothing

	return nil
}

func TestMockMutator(t *testing.T) {
	// Mock is not registered
	for _, name := range List() {
		if name == "mock" {
			assert.Fail(t, "mock should not be in the mutator list yet")
		}
	}

	m, err := New("mock")
	assert.Nil(t, m)
	assert.NotNil(t, err)

	// Register mock
	Register("mock", mockMutator)

	// Mock is registered
	found := false
	for _, name := range List() {
		if name == "mock" {
			found = true

			break
		}
	}
	assert.True(t, found)

	m, err = New("mock")
	assert.NotNil(t, m)
	assert.Nil(t, err)

	// Register mock a second time
	caught := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		Register("mock", mockMutator)
	}()
	assert.True(t, caught)

	// Register nil function
	caught = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		Register("mockachino", nil)
	}()
	assert.True(t, caught)
}
