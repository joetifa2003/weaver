package parser

import (
	"errors"
	"strconv"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func intExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.TokenType(TT_INT),
		func(s string) (ast.Expr, error) {
			val, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

			return ast.IntExpr{Value: val}, nil
		},
	)
}

func floatExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.TokenType(TT_FLOAT),
		func(lhs string) (ast.Expr, error) {
			val, err := strconv.ParseFloat(lhs, 64)
			if err != nil {
				return nil, err
			}

			return ast.FloatExpr{Value: val}, nil
		},
	)
}

func booleanExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.OneOf(
			pargo.Exactly("true"),
			pargo.Exactly("false"),
		),
		func(s string) (ast.Expr, error) {
			if s != "true" && s != "false" {
				return nil, errors.New("invalid boolean") // TODO: handle errors in a better way
			}

			return ast.BoolExpr{Value: s == "true"}, nil
		},
	)
}

func stringExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.TokenType(TT_STRING),
		func(val string) (ast.Expr, error) {
			return ast.StringExpr{Value: val}, nil
		},
	)
}

func binaryExpr(operand pargo.Parser[ast.Expr], op string) pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.ManySep(operand, pargo.Exactly(op)),
		func(exprs []ast.Expr) (ast.Expr, error) {
			if len(exprs) == 1 {
				return exprs[0], nil
			}

			return ast.BinaryExpr{Operands: exprs, Operator: op}, nil
		},
	)
}

func identExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.TokenType(TT_IDENT),
		func(s string) (ast.Expr, error) {
			return ast.IdentExpr{Name: s}, nil
		},
	)
}

func atom() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		intExpr(),
		floatExpr(),
		booleanExpr(),
		stringExpr(),
		identExpr(),
	)
}

func expr() pargo.Parser[ast.Expr] {
	return lessThanExpr()
}

func lessThanExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		addExpr(),
		"<",
	)
}

func addExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		subExpr(),
		"+",
	)
}

func subExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		mulExpr(),
		"-",
	)
}

func mulExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		divExpr(),
		"*",
	)
}

func divExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		atom(),
		"/",
	)
}
