package parser

import (
	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func varDeclStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence4(
		pargo.TokenType(TT_IDENT),
		pargo.Exactly(":="),
		expr(),
		pargo.Optional(pargo.Exactly(";")),
		func(name string, _ string, expr ast.Expr, _ *string) ast.Statement {
			return ast.LetStmt{Name: name, Expr: expr}
		},
	)
}

func blockStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		pargo.ManyUntil(pargo.Lazy(stmt), pargo.Exactly("}")),
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
	return pargo.OneOf(
		pargo.Sequence6(
			pargo.Exactly("for"),
			pargo.Lazy(stmt),
			pargo.Lazy(expr),
			pargo.Exactly(";"),
			pargo.Lazy(expr),
			blockStmt(),
			func(_ string, initStmt ast.Statement, condition ast.Expr, _ string, increment ast.Expr, body ast.Statement) ast.Statement {
				return ast.ForStmt{
					InitStmt:  initStmt,
					Condition: condition,
					Increment: increment,
					Body:      body,
				}
			},
		),
		pargo.Sequence5(
			pargo.Exactly("for"),
			pargo.TokenType(TT_IDENT),
			pargo.Exactly("in"),
			rangeAtom(),
			pargo.Lazy(stmt),
			func(_ string, variable string, _ string, r rangeType, body ast.Statement) ast.Statement {
				return ast.ForRangeStmt{
					Variable: variable,
					Start:    r.start,
					End:      r.end,
					Body:     body,
				}
			},
		),
	)
}

type rangeType struct {
	start ast.Expr
	end   ast.Expr
}

func rangeAtom() pargo.Parser[rangeType] {
	return pargo.Sequence3(
		pargo.Lazy(expr),
		pargo.Exactly(".."),
		pargo.Lazy(expr),
		func(start ast.Expr, _ string, end ast.Expr) rangeType {
			return rangeType{start, end}
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

func exprStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence2(
		expr(),
		pargo.Optional(pargo.Exactly(";")),
		func(expr ast.Expr, _ *string) ast.Statement {
			return ast.ExprStmt{Expr: expr}
		},
	)
}

func continueStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence(
		func(_ string) ast.Statement {
			return ast.ContinueStmt{}
		},
		pargo.Exactly("continue"),
	)
}

func breakStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence(
		func(_ string) ast.Statement {
			return ast.BreakStmt{}
		},
		pargo.Exactly("break"),
	)
}

func matchStmt() pargo.Parser[ast.Statement] {
	return pargo.Sequence5(
		pargo.Exactly("match"),
		expr(),
		pargo.Exactly("{"),
		pargo.ManySep(matchCase(), pargo.Exactly(",")),
		pargo.Exactly("}"),
		func(_ string, expr ast.Expr, _ string, cases []ast.MatchCase, _ string) ast.Statement {
			return ast.MatchStmt{Expr: expr, Cases: cases}
		},
	)
}

func matchCase() pargo.Parser[ast.MatchCase] {
	return pargo.Sequence4(
		matchCondition(),
		pargo.Optional(
			pargo.Sequence2(
				pargo.Exactly("if"),
				pargo.Lazy(expr),
				func(_ string, cond ast.Expr) ast.Expr {
					return cond
				},
			),
		),
		pargo.Exactly("=>"),
		pargo.Lazy(stmt),
		func(cond ast.MatchCaseCondition, extraCond *ast.Expr, _ string, body ast.Statement) ast.MatchCase {
			return ast.MatchCase{Condition: cond, Body: body, ExtraCond: extraCond}
		},
	)
}

// matchConditionBase parses a single, non-OR condition.
func matchConditionBase() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.OneOf(
		matchCaseTypeError(),
		matchCaseTypeNumber(),
		matchCaseTypeString(),
		matchRangeCondition(),
		matchCaseInt(),
		matchCaseFloat(),
		matchCaseString(),
		matchCaseArray(),
		matchCaseObject(),
		matchCaseIdent(), // Keep ident last as it's the most general
	)
}

// matchCondition parses one or more base conditions separated by '|'
func matchCondition() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		pargo.SomeSep(matchConditionBase(), pargo.Exactly("|")),
		func(conditions []ast.MatchCaseCondition) (ast.MatchCaseCondition, error) {
			if len(conditions) == 1 {
				return conditions[0], nil
			}
			return ast.MatchCaseOr{Conditions: conditions}, nil
		},
	)
}

func matchRangeCondition() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Sequence3(
		pargo.OneOf(
			floatExpr(),
			intExpr(),
		),
		pargo.Exactly(".."),
		pargo.OneOf(
			floatExpr(),
			intExpr(),
		),
		func(left ast.Expr, _ string, right ast.Expr) ast.MatchCaseCondition {
			return ast.MatchCaseRange{Begin: left, End: right}
		},
	)
}

func matchCaseInt() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		intExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			intExpr := expr.(ast.IntExpr)
			return ast.MatchCaseInt(intExpr), nil
		},
	)
}

func matchCaseFloat() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		floatExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			floatExpr := expr.(ast.FloatExpr)
			return ast.MatchCaseFloat(floatExpr), nil
		},
	)
}

func matchCaseString() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		stringExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			stringExpr := expr.(ast.StringExpr)
			return ast.MatchCaseString(stringExpr), nil
		},
	)
}

func matchCaseArray() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Sequence3(
		pargo.Exactly("["),
		pargo.ManySep(pargo.Lazy(matchCondition), pargo.Exactly(",")),
		pargo.Exactly("]"),
		func(_ string, conditions []ast.MatchCaseCondition, _ string) ast.MatchCaseCondition {
			return ast.MatchCaseArray{Conditions: conditions}
		},
	)
}

type keyValue[K comparable, V any] struct {
	key   K
	value *V
}

func matchCaseObject() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		pargo.ManySep(
			pargo.OneOf(
				pargo.Sequence3(
					pargo.TokenType(TT_IDENT),
					pargo.Exactly(":"),
					pargo.Lazy(matchCondition),
					func(key string, _ string, value ast.MatchCaseCondition) keyValue[string, ast.MatchCaseCondition] {
						return keyValue[string, ast.MatchCaseCondition]{key, &value}
					},
				),
				pargo.Map(pargo.TokenType(TT_IDENT), func(ident string) (keyValue[string, ast.MatchCaseCondition], error) {
					return keyValue[string, ast.MatchCaseCondition]{ident, nil}, nil
				}),
			),
			pargo.Exactly(","),
		),
		pargo.Exactly("}"),
		func(_ string, kvs []keyValue[string, ast.MatchCaseCondition], _ string) ast.MatchCaseCondition {
			m := map[string]*ast.MatchCaseCondition{}
			for _, kv := range kvs {
				m[kv.key] = kv.value
			}
			return ast.MatchCaseObject{KVs: m}
		},
	)
}

func matchCaseType(typ string, f func(childCond *ast.MatchCaseCondition) ast.MatchCaseCondition) pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Sequence4(
		pargo.Exactly(typ),
		pargo.Exactly("("),
		pargo.Optional(pargo.Lazy(matchCondition)),
		pargo.Exactly(")"),
		func(_ string, _ string, cond *ast.MatchCaseCondition, _ string) ast.MatchCaseCondition {
			return f(cond)
		},
	)
}

func matchCaseTypeError() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.OneOf(
		// error() with no arguments
		pargo.Sequence3(
			pargo.Exactly("error"),
			pargo.Exactly("("),
			pargo.Exactly(")"),
			func(kw, lp, rp string) ast.MatchCaseCondition {
				return ast.MatchCaseTypeError{}
			},
		),
		// error(msg) with one argument
		pargo.Sequence4(
			pargo.Exactly("error"),
			pargo.Exactly("("),
			pargo.Lazy(matchCondition),
			pargo.Exactly(")"),
			func(kw, lp string, msgCond ast.MatchCaseCondition, rp string) ast.MatchCaseCondition {
				return ast.MatchCaseTypeError{
					Message: msgCond,
				}
			},
		),
		// error(msg, details) with two arguments
		pargo.Sequence6(
			pargo.Exactly("error"),
			pargo.Exactly("("),
			pargo.Lazy(matchCondition),
			pargo.Exactly(","),
			pargo.Lazy(matchCondition),
			pargo.Exactly(")"),
			func(kw, lp string, msgCond ast.MatchCaseCondition, comma string, detailsCond ast.MatchCaseCondition, rp string) ast.MatchCaseCondition {
				return ast.MatchCaseTypeError{
					Message: msgCond,
					Data:    detailsCond,
				}
			},
		),
	)
}

func matchCaseTypeNumber() pargo.Parser[ast.MatchCaseCondition] {
	return matchCaseType("number", func(cond *ast.MatchCaseCondition) ast.MatchCaseCondition {
		return ast.MatchCaseTypeNumber{Cond: cond}
	})
}

func matchCaseTypeString() pargo.Parser[ast.MatchCaseCondition] {
	return matchCaseType("string", func(cond *ast.MatchCaseCondition) ast.MatchCaseCondition {
		return ast.MatchCaseTypeString{Cond: cond}
	})
}

func matchCaseIdent() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		pargo.TokenType(TT_IDENT),
		func(s string) (ast.MatchCaseCondition, error) {
			return ast.MatchCaseIdent{Name: s}, nil
		},
	)
}

func stmt() pargo.Parser[ast.Statement] {
	return pargo.OneOf(
		varDeclStmt(),
		blockStmt(),
		whileStmt(),
		ifStmt(),
		forStmt(),
		continueStmt(),
		breakStmt(),
		matchStmt(),

		// keep this at the end
		exprStmt(),
	)
}
