package osutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	src := "copy.go"
	dst := "copy.go.tmp"

	err := CopyFile(src, dst)
	assert.Nil(t, err)

	s, err := ioutil.ReadFile(src)
	assert.Nil(t, err)

	d, err := ioutil.ReadFile(dst)
	assert.Nil(t, err)

	assert.Equal(t, s, d)

	err = os.Remove(dst)
	assert.Nil(t, err)
}
