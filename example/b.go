package example

import (
	"net/http"
)

func fooB() (a A, b http.Header) {
	a, b = A{}, http.Header{}

	return a, b
}
