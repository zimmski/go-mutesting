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

var blacklistedIdentifiers = map[string]bool{
	// blank identifier
	"_": true,
	// builtin - can be used as identifier but are unlikely to be in practice
	// (except perhaps with panic, defer, recover, print, prinln?)
	"bool":       true,
	"true":       true,
	"false":      true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
	"string":     true,
	"int":        true,
	"uint":       true,
	"uintptr":    true,
	"byte":       true,
	"rune":       true,
	"iota":       true,
	"nil":        true,
	"append":     true,
	"copy":       true,
	"delete":     true,
	"len":        true,
	"cap":        true,
	"make":       true,
	"new":        true,
	"complex":    true,
	"real":       true,
	"imag":       true,
	"close":      true,
	"panic":      true,
	"recover":    true,
	"print":      true,
	"println":    true,
	"error":      true,
	// reserved keywords - cannot be used as identifier.
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case":        true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
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
		if _, ok := blacklistedIdentifiers[n.Name]; !ok {
			w.identifiers = append(w.identifiers, n)
		}

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
