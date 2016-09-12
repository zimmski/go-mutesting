package astutil

import (
	"go/ast"
	"go/token"
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

var blacklistedIdentifiers = map[string]bool{
	// Built-ins
	"append":     true,
	"bool":       true,
	"byte":       true,
	"cap":        true,
	"close":      true,
	"complex":    true,
	"complex128": true,
	"complex64":  true,
	"copy":       true,
	"delete":     true,
	"error":      true,
	"false":      true,
	"float32":    true,
	"float64":    true,
	"imag":       true,
	"int":        true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"int8":       true,
	"iota":       true,
	"len":        true,
	"make":       true,
	"new":        true,
	"nil":        true,
	"panic":      true,
	"print":      true,
	"println":    true,
	"real":       true,
	"recover":    true,
	"rune":       true,
	"string":     true,
	"true":       true,
	"uint":       true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uint8":      true,
	"uintptr":    true,
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

		// Ignore our blacklist identifiers
		if _, ok := blacklistedIdentifiers[n.Name]; ok {
			return nil
		}

		w.identifiers = append(w.identifiers, n)

		return nil
	case *ast.SelectorExpr:
		if checkForSelectorExpr(n) {
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
