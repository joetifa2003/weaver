package parser

import (
	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
	"github.com/joetifa2003/weaver/internal/pargo/lexer"
)

func program() pargo.Parser[ast.Program] {
	return pargo.Map(pargo.Many(stmt()), func(stmts []ast.Statement) (ast.Program, error) {
		return ast.Program{Statements: stmts}, nil
	})
}

func Parse(src string) (ast.Program, error) {
	p, err := pargo.Parse(program(), lexer.New(), src)
	if err != nil {
		return ast.Program{}, err
	}

	return p, nil
}
