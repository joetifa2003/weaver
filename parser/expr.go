package parser

import (
	"errors"
	"strconv"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func intExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.Concat(pargo.Some(digit())),

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
	return pargo.Sequence3(
		pargo.Concat(pargo.Some(digit())),
		pargo.Exactly("."),
		pargo.Concat(pargo.Some(digit())),
		func(lhs string, _ string, rhs string) ast.Expr {
			val, err := strconv.ParseFloat(lhs+"."+rhs, 64)
			if err != nil {
				panic(err) // TODO: make sequence return error
			}

			return ast.FloatExpr{Value: val}
		},
	)
}

func booleanExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		identifier(),
		func(s string) (ast.Expr, error) {
			if s != "true" && s != "false" {
				return nil, errors.New("invalid boolean") // TODO: handle errors in a better way
			}

			return ast.BoolExpr{Value: s == "true"}, nil
		},
	)
}

func stringExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence3(
		pargo.Exactly(`"`),
		pargo.Concat(pargo.Many(pargo.Except(`"`))),
		pargo.Exactly(`"`),
		func(_ string, val string, _ string) ast.Expr {
			return ast.StringExpr{Value: val}
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

func atom() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		intExpr(),
		floatExpr(),
		booleanExpr(),
		stringExpr(),
	)
}

func expr() pargo.Parser[ast.Expr] {
	return addExpr()
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
