package parser

import (
	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func varDeclStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence4(
		identifier(),
		pargo.Exactly(":"),
		pargo.Exactly("="),
		expr(),
		func(name string, _ string, _ string, expr ast.Expr) ast.Statement {
			return ast.LetStmt{Name: name, Expr: expr}
		},
	)
}
