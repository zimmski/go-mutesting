package astutil

import (
	"go/ast"
	"go/token"
	"go/types"
)

// IdentifiersInStatement returns all identifiers with their found in a statement.
func IdentifiersInStatement(pkg *types.Package, info *types.Info, stmt ast.Stmt) []ast.Expr {
	w := &identifierWalker{
		pkg:  pkg,
		info: info,
	}

	ast.Walk(w, stmt)

	return w.identifiers
}

type identifierWalker struct {
	identifiers []ast.Expr
	pkg         *types.Package
	info        *types.Info
}

func checkForSelectorExpr(node ast.Expr) bool {
	switch n := node.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return checkForSelectorExpr(n.X)
	}

	return false
}

func (w *identifierWalker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		// Ignore the blank identifier
		if n.Name == "_" {
			return nil
		}

		// Ignore keywords
		if token.Lookup(n.Name) != token.IDENT {
			return nil
		}

		// We are only interested in variables
		if obj, ok := w.info.Uses[n]; ok {
			if _, ok := obj.(*types.Var); !ok {
				return nil
			}
		}

		// FIXME instead of manually creating a new node, clone it and trim the node from its comments and position https://github.com/zimmski/go-mutesting/issues/49
		w.identifiers = append(w.identifiers, &ast.Ident{
			Name: n.Name,
		})

		return nil
	case *ast.SelectorExpr:
		if !checkForSelectorExpr(n) {
			return nil
		}

		// Check if we need to instantiate the expression
		initialize := false
		if n.Sel != nil {
			if obj, ok := w.info.Uses[n.Sel]; ok {
				t := obj.Type()

				switch t.Underlying().(type) {
				case *types.Array, *types.Map, *types.Slice, *types.Struct:
					initialize = true
				}
			}
		}

		if initialize {
			// FIXME we need to clone the node and trim comments and position recursively https://github.com/zimmski/go-mutesting/issues/49
			w.identifiers = append(w.identifiers, &ast.CompositeLit{
				Type: n,
			})
		} else {
			// FIXME we need to clone the node and trim comments and position recursively https://github.com/zimmski/go-mutesting/issues/49
			w.identifiers = append(w.identifiers, n)
		}

		return nil
	}

	return w
}

// Functions returns all found functions.
func Functions(n ast.Node) []*ast.FuncDecl {
	w := &functionWalker{}

	ast.Walk(w, n)

	return w.functions
}

type functionWalker struct {
	functions []*ast.FuncDecl
}

func (w *functionWalker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.functions = append(w.functions, n)

		return nil
	}

	return w
}
