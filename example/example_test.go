package example

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	Equal(t, foo(), 16)
}
