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

func whileStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence3(
		pargo.Exactly("while"),
		expr(),
		blockStmt(),
		func(_ string, condition ast.Expr, statement ast.Statement) ast.Statement {
			return ast.WhileStmt{Condition: condition, Body: statement}
		},
	)
}

func forStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence7(
		pargo.Exactly("for"),
		pargo.Lazy(stmt),
		pargo.Exactly(";"),
		pargo.Lazy(expr),
		pargo.Exactly(";"),
		pargo.Lazy(expr),
		blockStmt(),
		func(_ string, initStmt ast.Statement, _ string, condition ast.Expr, _ string, increment ast.Expr, body ast.Statement) ast.Statement {
			return ast.BlockStmt{
				Statements: []ast.Statement{
					initStmt,
					ast.WhileStmt{
						Condition: condition,
						Body: ast.BlockStmt{
							Statements: []ast.Statement{
								body,
								ast.ExprStmt{
									Expr: increment,
								},
							},
						},
					},
				},
			}
		},
	)
}

func ifStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence4(
		pargo.Exactly("if"),
		expr(),
		blockStmt(),
		pargo.Optional(
			pargo.Sequence2(
				pargo.Exactly("else"),
				pargo.OneOf(blockStmt(), pargo.Lazy(ifStmt)),
				func(_ string, alternative ast.Statement) ast.Statement {
					return alternative
				},
			),
		),
		func(_ string, condition ast.Expr, body ast.Statement, alternative *ast.Statement) ast.Statement {
			return ast.IfStmt{Condition: condition, Body: body, Alternative: alternative}
		},
	)
}

func returnStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence2(
		pargo.Exactly("return"),
		expr(),
		func(_ string, expr ast.Expr) ast.Statement {
			return ast.ReturnStmt{Expr: expr}
		},
	)
}

func exprStmt() pargo.Parser[ast.Statement] {
	return pargo.Map(
		expr(),
		func(expr ast.Expr) (ast.Statement, error) {
			return ast.ExprStmt{Expr: expr}, nil
		},
	)
}

func stmt() pargo.Parser[ast.Statement] {
	return pargo.OneOf(
		varDeclStmt(),
		blockStmt(),
		echoStmt(),
		whileStmt(),
		ifStmt(),
		returnStmt(),
		forStmt(),

		// keep this at the end
		exprStmt(),
	)
}
