package astutil

import (
	"go/ast"
	"go/token"
)

// CreateNoopOfStatement creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatement(stmt ast.Stmt) ast.Stmt {
	return CreateNoopOfStatements([]ast.Stmt{stmt})
}

// CreateNoopOfStatements creates a syntactically safe noop statement out of a given statement.
func CreateNoopOfStatements(stmts []ast.Stmt) ast.Stmt {
	var ids []ast.Expr
	for _, stmt := range stmts {
		ids = append(ids, IdentifiersInStatement(stmt)...)
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
