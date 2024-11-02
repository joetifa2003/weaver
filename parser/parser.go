package parser

import (
	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func program() pargo.Parser[ast.Program] {
	return pargo.Map(
		pargo.ManyAll(stmt()),
		func(stmts []ast.Statement) (ast.Program, error) {
			return ast.Program{Statements: stmts}, nil
		})
}

func Parse(src string) (ast.Program, error) {
	p, err := pargo.Parse(program(), newLexer(), src)
	if err != nil {
		return ast.Program{}, err
	}

	return p, nil
}
