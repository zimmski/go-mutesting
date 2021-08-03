// +build tagged

package example

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestGreaterThan(t *testing.T) {
	Equal(t, gt(2,1), true)
}

func TestLessThan(t *testing.T) {
	Equal(t, gt(1,2), false)
}

func TestEqual(t *testing.T) {
	Equal(t, gt(2,2), false)
}
