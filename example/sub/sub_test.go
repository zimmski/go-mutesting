package sub

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestBaz(t *testing.T) {
	Equal(t, baz(), 2)
}
