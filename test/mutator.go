package test

import (
	"bytes"
	"fmt"
	"go/printer"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zimmski/go-mutesting"
	"github.com/zimmski/go-mutesting/mutator"
)

// Mutator tests a mutator.
// It mutates the given original file with the given mutator. Every mutation is then validated with the given changed file. The mutation overall count is validated with the given count.
func Mutator(t *testing.T, m mutator.Mutator, testFile string, count uint) {
	// test if mutator is not nil
	assert.NotNil(t, m)

	// test if the mutator is properly registered
	mn, err := mutator.New(m.String())
	assert.Nil(t, err)
	assert.NotNil(t, mn)

	// read and parse the original source code
	originalFile, err := ioutil.ReadFile(testFile)
	assert.Nil(t, err)

	f, fset, err := mutesting.ParseSource(originalFile)
	assert.Nil(t, err)

	// count the actual mutations
	n := mutesting.CountWalk(f, m)
	assert.Equal(t, n, count)

	// mutate a non relevant node
	changed := make(chan bool)

	go m.Mutate(f, changed)
	assert.False(t, <-changed)

	// mutate all relevant nodes -> test whole mutation process
	changed = mutesting.MutateWalk(f, m)

	for i := uint(0); i < count; i++ {
		assert.True(t, <-changed)

		buf := new(bytes.Buffer)
		err = printer.Fprint(buf, fset, f)
		assert.Nil(t, err)

		changedFile, err := ioutil.ReadFile(fmt.Sprintf("%s.%d.go", testFile, i))
		assert.Nil(t, err)

		assert.Equal(t, string(changedFile), buf.String())

		changed <- true

		assert.True(t, <-changed)

		buf = new(bytes.Buffer)
		err = printer.Fprint(buf, fset, f)
		assert.Nil(t, err)

		assert.Equal(t, string(originalFile), buf.String())

		changed <- true
	}

	_, ok := <-changed
	assert.False(t, ok)
}
