package astutil

import (
	"go/ast"
	"go/token"
	"go/types"
)

// CreateNoopOfStatement creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatement(pkg *types.Package, info *types.Info, stmt ast.Stmt) ast.Stmt {
	return CreateNoopOfStatements(pkg, info, []ast.Stmt{stmt})
}

// CreateNoopOfStatements creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatements(pkg *types.Package, info *types.Info, stmts []ast.Stmt) ast.Stmt {
	var ids []ast.Expr
	for _, stmt := range stmts {
		ids = append(ids, IdentifiersInStatement(pkg, info, stmt)...)
	}

	if len(ids) == 0 {
		return &ast.EmptyStmt{
			Semicolon: token.NoPos,
		}
	}

	lhs := make([]ast.Expr, len(ids))
	for i := range ids {
		lhs[i] = ast.NewIdent("_")
	}

	return &ast.AssignStmt{
		Lhs: lhs,
		Rhs: ids,
		Tok: token.ASSIGN,
	}
}
