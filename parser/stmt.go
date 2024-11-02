package parser

import (
	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func varDeclStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence3(
		pargo.TokenType(TT_IDENT),
		pargo.Exactly(":="),
		expr(),
		func(name string, _ string, expr ast.Expr) ast.Statement {
			return ast.LetStmt{Name: name, Expr: expr}
		},
	)
}

func echoStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence2(
		pargo.Exactly("echo"),
		expr(),
		func(_ string, expr ast.Expr) ast.Statement {
			return ast.EchoStmt{Expr: expr}
		},
	)
}

func blockStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		pargo.Many(pargo.Lazy(stmt)),
		pargo.Exactly("}"),
		func(_ string, stmts []ast.Statement, _ string) ast.Statement {
			return ast.BlockStmt{Statements: stmts}
		},
	)
}

func stmt() pargo.Parser[ast.Statement] {
	return pargo.OneOf(
		blockStmt(),
		echoStmt(),
		varDeclStmt(),
	)
}
