package parser

import (
	"errors"
	"strconv"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func binaryExpr(operand pargo.Parser[ast.Expr], op string) pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.SomeSep(operand, pargo.Exactly(op)),
		func(exprs []ast.Expr) (ast.Expr, error) {
			if len(exprs) == 1 {
				return exprs[0], nil
			}

			return ast.BinaryExpr{Operands: exprs, Operator: op}, nil
		},
	)
}

func expr() pargo.Parser[ast.Expr] {
	return orExpr()
}

func orExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		andExpr(),
		"||",
	)
}

func andExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		pipeExpr(),
		"&&",
	)
}

func pipeExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		equalityExpr(),
		"|",
	)
}

func equalityExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		nequalExpr(),
		"==",
	)
}

func nequalExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		lessThanExpr(),
		"!=",
	)
}

func lessThanExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		lessThanEqualExpr(),
		"<",
	)
}

func lessThanEqualExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		greaterThanEqualExpr(),
		"<=",
	)
}

func greaterThanEqualExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		greaterThanExpr(),
		">=",
	)
}

func greaterThanExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		addExpr(),
		">",
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
		modExpr(),
		"-",
	)
}

func modExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		mulExpr(),
		"%",
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
		notExpr(),
		"/",
	)
}

func notExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence2(
			pargo.Exactly("!"),
			pargo.Lazy(notExpr),
			func(_ string, expr ast.Expr) ast.Expr {
				return ast.UnaryExpr{Operator: "!", Expr: expr}
			},
		),
		callExpr(),
	)
}

func callExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence4(
			atom(),
			pargo.Exactly("("),
			pargo.ManySep(pargo.Lazy(expr), pargo.Exactly(",")),
			pargo.Exactly(")"),
			func(callee ast.Expr, _ string, args []ast.Expr, _ string) ast.Expr {
				return ast.CallExpr{Callee: callee, Args: args}
			},
		),
		atom(),
	)
}

func atom() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		intExpr(),
		floatExpr(),
		booleanExpr(),
		stringExpr(),
		assignExpr(),
		identExpr(),
		functionExpr(),
		lambdaExpr(),
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

func assignExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence3(
		pargo.TokenType(TT_IDENT),
		pargo.Exactly("="),
		pargo.Lazy(expr),
		func(name string, _ string, expr ast.Expr) ast.Expr {
			return ast.AssignExpr{Name: name, Expr: expr}
		},
	)
}

func functionExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence4(
		pargo.Exactly("|"),
		pargo.ManySep(pargo.TokenType(TT_IDENT), pargo.Exactly(",")),
		pargo.Exactly("|"),
		blockStmt(),
		func(_ string, params []string, _ string, body ast.Statement) ast.Expr {
			return ast.FunctionExpr{Params: params, Body: body}
		},
	)
}

func lambdaExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence4(
		pargo.Exactly("|"),
		pargo.ManySep(pargo.TokenType(TT_IDENT), pargo.Exactly(",")),
		pargo.Exactly("|"),
		pargo.Lazy(expr),
		func(_ string, params []string, _ string, expr ast.Expr) ast.Expr {
			return ast.FunctionExpr{
				Params: params,
				Body: ast.ReturnStmt{
					Expr: expr,
				},
			}
		},
	)
}

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
