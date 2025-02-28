package parser

import (
	"errors"
	"strconv"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func binaryExpr(
	operand pargo.Parser[ast.Expr],
	separator pargo.Parser[string],
	op ast.BinaryOp,
) pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.SomeSep(operand, separator),
		func(exprs []ast.Expr) (ast.Expr, error) {
			if len(exprs) == 1 {
				return exprs[0], nil
			}

			return ast.BinaryExpr{Operands: exprs, Operator: op}, nil
		},
	)
}

func expr() pargo.Parser[ast.Expr] {
	return ternaryExpr()
}

type ternaryBody struct {
	trueExpr  ast.Expr
	falseExpr ast.Expr
}

func ternaryExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence2(
			orExpr(),
			pargo.Optional(
				pargo.Sequence4(
					pargo.Exactly("?"),
					pargo.Lazy(expr),
					pargo.Exactly(":"),
					pargo.Lazy(expr),
					func(_ string, trueExpr ast.Expr, _ string, falseExpr ast.Expr) ternaryBody {
						return ternaryBody{trueExpr, falseExpr}
					},
				),
			),
			func(expr ast.Expr, body *ternaryBody) ast.Expr {
				if body == nil {
					return expr
				}

				return ast.TernaryExpr{
					Expr:      expr,
					TrueExpr:  body.trueExpr,
					FalseExpr: body.falseExpr,
				}
			},
		),
		orExpr(),
	)
}

func orExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		andExpr(),
		pargo.Sequence2(
			pargo.Exactly("|"), pargo.Exactly("|"),
			func(_ string, _ string) string { return "" },
		),
		ast.BinaryOpOr,
	)
}

func andExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		pipeExpr(),
		pargo.Exactly(string(ast.BinaryOpAnd)),
		ast.BinaryOpAnd,
	)
}

func pipeExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		equalityExpr(),
		pargo.Exactly(string(ast.BinaryOpPipe)),
		ast.BinaryOpPipe,
	)
}

func equalityExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		nequalExpr(),
		pargo.Exactly(string(ast.BinaryOpEq)),
		ast.BinaryOpEq,
	)
}

func nequalExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		lessThanExpr(),
		pargo.Exactly(string(ast.BinaryOpNeq)),
		ast.BinaryOpNeq,
	)
}

func lessThanExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		lessThanEqualExpr(),
		pargo.Exactly(string(ast.BinaryOpLt)),
		ast.BinaryOpLt,
	)
}

func lessThanEqualExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		greaterThanEqualExpr(),
		pargo.Exactly(string(ast.BinaryOpLte)),
		ast.BinaryOpLte,
	)
}

func greaterThanEqualExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		greaterThanExpr(),
		pargo.Exactly(string(ast.BinaryOpGte)),
		ast.BinaryOpGte,
	)
}

func greaterThanExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		addExpr(),
		pargo.Exactly(string(ast.BinaryOpGt)),
		ast.BinaryOpGt,
	)
}

func addExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		subExpr(),
		pargo.Exactly(string(ast.BinaryOpAdd)),
		ast.BinaryOpAdd,
	)
}

func subExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		modExpr(),
		pargo.Exactly(string(ast.BinaryOpSub)),
		ast.BinaryOpSub,
	)
}

func modExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		mulExpr(),
		pargo.Exactly(string(ast.BinaryOpMod)),
		ast.BinaryOpMod,
	)
}

func mulExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		divExpr(),
		pargo.Exactly(string(ast.BinaryOpMul)),
		ast.BinaryOpMul,
	)
}

func divExpr() pargo.Parser[ast.Expr] {
	return binaryExpr(
		unaryExpr(),
		pargo.Exactly(string(ast.BinaryOpDiv)),
		ast.BinaryOpDiv,
	)
}

func unaryExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence2(
			pargo.Exactly(string(ast.UnaryOpNot)),
			pargo.Lazy(unaryExpr),
			func(_ string, expr ast.Expr) ast.Expr {
				return ast.UnaryExpr{Operator: ast.UnaryOpNot, Expr: expr}
			},
		),
		pargo.Sequence2(
			pargo.Exactly(string(ast.UnaryOpNegate)),
			pargo.Lazy(unaryExpr),
			func(_ string, expr ast.Expr) ast.Expr {
				return ast.UnaryExpr{Operator: ast.UnaryOpNegate, Expr: expr}
			},
		),
		assignExpr(),
	)
}

func assignExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence3(
			postFixExpr(),
			pargo.Exactly("="),
			pargo.Lazy(expr),
			func(assignee ast.Expr, _ string, expr ast.Expr) ast.Expr {
				return ast.AssignExpr{Assignee: assignee, Expr: expr}
			},
		),
		postFixExpr(),
	)
}

func postFixExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence2(
		increment(),
		pargo.Many(
			pargo.OneOf(
				postFixIndexOp(),
				postFixCallOp(),
				postFixDotOp(),
			),
		),
		func(expr ast.Expr, ops []ast.PostFixOp) ast.Expr {
			if len(ops) == 0 {
				return expr
			}

			return ast.PostFixExpr{Expr: expr, Ops: ops}
		},
	)
}

func postFixIndexOp() pargo.Parser[ast.PostFixOp] {
	return pargo.Sequence3(
		pargo.Exactly("["),
		pargo.Lazy(expr),
		pargo.Exactly("]"),
		func(_ string, expr ast.Expr, _ string) ast.PostFixOp {
			return ast.IndexOp{Index: expr}
		},
	)
}

func postFixCallOp() pargo.Parser[ast.PostFixOp] {
	return pargo.Sequence3(
		pargo.Exactly("("),
		pargo.ManySep(pargo.Lazy(expr), pargo.Exactly(",")),
		pargo.Exactly(")"),
		func(_ string, args []ast.Expr, _ string) ast.PostFixOp {
			return ast.CallOp{Args: args}
		},
	)
}

func postFixDotOp() pargo.Parser[ast.PostFixOp] {
	return pargo.Sequence2(
		pargo.Exactly("."),
		pargo.TokenType(TT_IDENT),
		func(_ string, ident string) ast.PostFixOp {
			return ast.DotOp{
				Index: ident,
			}
		},
	)
}

func increment() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence2(
			identExpr(),
			pargo.Exactly("++"),
			func(expr ast.Expr, _ string) ast.Expr {
				return ast.VarIncrementExpr{Name: expr.(ast.IdentExpr).Name}
			},
		),
		pargo.Sequence2(
			identExpr(),
			pargo.Exactly("--"),
			func(expr ast.Expr, _ string) ast.Expr {
				return ast.VarDecrementExpr{Name: expr.(ast.IdentExpr).Name}
			},
		),
		moduleIndexExpr(),
	)
}

func moduleIndexExpr() pargo.Parser[ast.Expr] {
	return pargo.OneOf(
		pargo.Sequence3(
			pargo.TokenType(TT_IDENT),
			pargo.Exactly(":"),
			pargo.TokenType(TT_IDENT),
			func(name string, _ string, val string) ast.Expr {
				return ast.ModuleLoadExpr{Name: name, Load: val}
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
		nilExpr(),
		identExpr(),
		functionExpr(),
		lambdaExpr(),
		arrayExpr(),
		objectExpr(),
		parenExpr(),
	)
}

func parenExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence3(
		pargo.Exactly("("),
		pargo.Lazy(expr),
		pargo.Exactly(")"),
		func(_ string, expr ast.Expr, _ string) ast.Expr {
			return expr
		},
	)
}

func nilExpr() pargo.Parser[ast.Expr] {
	return pargo.Map(
		pargo.Exactly("nil"),
		func(_ string) (ast.Expr, error) {
			return ast.NilExpr{}, nil
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

func paramList() pargo.Parser[[]string] {
	return pargo.Sequence3(
		pargo.Exactly("|"),
		pargo.ManySep(pargo.TokenType(TT_IDENT), pargo.Exactly(",")),
		pargo.Exactly("|"),
		func(_ string, params []string, _ string) []string {
			return params
		},
	)
}

func functionExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence2(
		paramList(),
		blockStmt(),
		func(params []string, body ast.Statement) ast.Expr {
			return ast.FunctionExpr{Params: params, Body: body}
		},
	)
}

func lambdaExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence2(
		paramList(),
		pargo.Lazy(expr),
		func(params []string, expr ast.Expr) ast.Expr {
			return ast.LambdaExpr{
				Params: params,
				Expr:   expr,
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

func arrayExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence3(
		pargo.Exactly("["),
		pargo.ManySep(pargo.Lazy(expr), pargo.Exactly(",")),
		pargo.Exactly("]"),
		func(_ string, exprs []ast.Expr, _ string) ast.Expr {
			return ast.ArrayExpr{Exprs: exprs}
		},
	)
}

func objectExpr() pargo.Parser[ast.Expr] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		objectKV(),
		pargo.Exactly("}"),
		func(_ string, kv map[string]ast.Expr, _ string) ast.Expr {
			return ast.ObjectExpr{KVs: kv}
		},
	)
}

type kv struct {
	key   string
	value ast.Expr
}

func objectKV() pargo.Parser[map[string]ast.Expr] {
	return pargo.Map(
		pargo.ManySep(
			pargo.Sequence3(
				pargo.OneOf(
					pargo.TokenType(TT_IDENT),
					pargo.TokenType(TT_STRING),
				),
				pargo.Exactly(":"),
				pargo.Lazy(expr),
				func(key string, _ string, expr ast.Expr) kv {
					return kv{key, expr}
				},
			),
			pargo.Exactly(","),
		),
		func(kvs []kv) (map[string]ast.Expr, error) {
			m := map[string]ast.Expr{}
			for _, kv := range kvs {
				m[kv.key] = kv.value
			}
			return m, nil
		},
	)
}
