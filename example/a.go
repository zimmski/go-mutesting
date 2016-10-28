package example

import (
	"net/http"
)

type A struct{}

func fooA() (a A, b http.Header) {
	_, _, _ = a, b, http.Header{}

	return a, b
}
