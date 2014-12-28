package mutesting

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

// ParseFile parses the content of the given file and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseFile(f string) (*ast.File, *token.FileSet, error) {
	src, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, nil, err
	}

	return ParseSource(src)
}

// ParseSource parses the given source and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseSource(src interface{}) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments|parser.AllErrors)

	return f, fset, err
}
