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
func Mutator(t *testing.T, m mutator.Mutator, testFile string, count int) {
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

	// mutate a non relevant node
	assert.Nil(t, m.Mutations(f))

	// count the actual mutations
	n := mutesting.CountWalk(f, m)
	assert.Equal(t, count, n)

	// mutate all relevant nodes -> test whole mutation process
	changed := mutesting.MutateWalk(f, m)

	for i := 0; i < count; i++ {
		assert.True(t, <-changed)

		buf := new(bytes.Buffer)
		err = printer.Fprint(buf, fset, f)
		assert.Nil(t, err)

		changedFilename := fmt.Sprintf("%s.%d.go", testFile, i)
		changedFile, err := ioutil.ReadFile(changedFilename)
		assert.Nil(t, err)

		assert.Equal(t, string(changedFile), buf.String(), fmt.Sprintf("For change file %q", changedFilename))

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
