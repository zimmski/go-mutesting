package astutil

import (
	"go/ast"
)

// IdentifiersInStatement returns all identifiers with their found in a statement.
func IdentifiersInStatement(stmt ast.Stmt) []ast.Expr {
	w := &identifierWalker{}

	ast.Walk(w, stmt)

	return w.identifiers
}

type identifierWalker struct {
	identifiers []ast.Expr
}

func (w *identifierWalker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		if n.Name != "nil" {
			w.identifiers = append(w.identifiers, n)
		}

		return nil
	case *ast.SelectorExpr:
		w.identifiers = append(w.identifiers, n)

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
