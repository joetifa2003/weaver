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
			return ast.ForStmt{
				InitStmt:  initStmt,
				Condition: condition,
				Increment: increment,
				Body:      body,
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
	return pargo.Sequence3(
		matchCondition(),
		pargo.Exactly("=>"),
		pargo.Lazy(stmt),
		func(cond ast.MatchCaseCondition, _ string, body ast.Statement) ast.MatchCase {
			return ast.MatchCase{Condition: cond, Body: body}
		},
	)
}

func matchCondition() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.OneOf(
		matchCaseInt(),
		matchCaseFloat(),
		matchCaseString(),
		matchCaseArray(),
		matchCaseObject(),
	)
}

func matchCaseInt() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		intExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			intExpr := expr.(ast.IntExpr)
			return ast.MatchCaseInt{Value: intExpr.Value}, nil
		},
	)
}

func matchCaseFloat() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		floatExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			floatExpr := expr.(ast.FloatExpr)
			return ast.MatchCaseFloat{Value: floatExpr.Value}, nil
		},
	)
}

func matchCaseString() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Map(
		stringExpr(),
		func(expr ast.Expr) (ast.MatchCaseCondition, error) {
			stringExpr := expr.(ast.StringExpr)
			return ast.MatchCaseString{Value: stringExpr.Value}, nil
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
	value V
}

func matchCaseObject() pargo.Parser[ast.MatchCaseCondition] {
	return pargo.Sequence3(
		pargo.Exactly("{"),
		pargo.ManySep(
			pargo.Sequence3(
				pargo.TokenType(TT_IDENT),
				pargo.Exactly(":"),
				pargo.Lazy(matchCondition),
				func(key string, _ string, value ast.MatchCaseCondition) keyValue[string, ast.MatchCaseCondition] {
					return keyValue[string, ast.MatchCaseCondition]{key, value}
				},
			),
			pargo.Exactly(","),
		),
		pargo.Exactly("}"),
		func(_ string, kvs []keyValue[string, ast.MatchCaseCondition], _ string) ast.MatchCaseCondition {
			m := map[string]ast.MatchCaseCondition{}
			for _, kv := range kvs {
				m[kv.key] = kv.value
			}
			return ast.MatchCaseObject{KVs: m}
		},
	)
}

func stmt() pargo.Parser[ast.Statement] {
	return pargo.OneOf(
		varDeclStmt(),
		blockStmt(),
		whileStmt(),
		ifStmt(),
		returnStmt(),
		forStmt(),
		continueStmt(),
		breakStmt(),
		matchStmt(),

		// keep this at the end
		exprStmt(),
	)
}
