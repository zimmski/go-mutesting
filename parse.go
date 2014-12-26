package mutesting

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func ParseFile(f string) (*ast.File, *token.FileSet, error) {
	src, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, nil, err
	}

	return ParseSource(src)
}

func ParseSource(src interface{}) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, 0)

	return f, fset, err
}
