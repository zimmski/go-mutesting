package example

import (
	"net/http"
)

type A struct{}

func fooA() (a A, b http.Header) {
	a, b = A{}, http.Header{}

	return a, b
}
