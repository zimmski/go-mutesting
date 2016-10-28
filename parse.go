package mutesting

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
)

// ParseFile parses the content of the given file and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseFile(file string) (*ast.File, *token.FileSet, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	return ParseSource(data)
}

// ParseSource parses the given source and returns the corresponding ast.File node and its file set for positional information.
// If a fatal error is encountered the error return argument is not nil.
func ParseSource(data interface{}) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()

	src, err := parser.ParseFile(fset, "", data, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	return src, fset, err
}

// ParseAndTypeCheckFile parses and type-checks the given file, and returns everything interesting about the file.
// If a fatal error is encountered the error return argument is not nil.
func ParseAndTypeCheckFile(file string) (*ast.File, *token.FileSet, *types.Package, *types.Info, error) {
	src, fset, err := ParseFile(file)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not open file %q: %v", file, err)
	}

	dir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not absolute the file path of %q: %v", file, err)
	}

	buildPkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not create build package of %q: %v", file, err)
	}

	conf := types.Config{
		Importer: importer.Default(),
	}

	info := &types.Info{
		Uses: make(map[*ast.Ident]types.Object),
	}

	pkg, err := conf.Check(buildPkg.ImportPath, fset, []*ast.File{src}, info) // TODO query the import path without the additional go/build.ImportDirt step
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not type-check file %q: %v", file, err)
	}

	return src, fset, pkg, info, nil
}
