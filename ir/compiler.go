package ir

import (
	"errors"
	"fmt"

	"github.com/joetifa2003/weaver/ast"
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(p ast.Program) ([]Statement, error) {
	var res []Statement

	for _, stmt := range p.Statements {
		stmtIr, err := c.CompileStmt(stmt)
		if err != nil {
			return nil, err
		}
		res = append(res, stmtIr)
	}

	return res, nil
}

func (c *Compiler) CompileStmt(s ast.Statement) (Statement, error) {
	switch s := s.(type) {
	case ast.BlockStmt:
		var res []Statement
		for _, stmt := range s.Statements {
			stmtIr, err := c.CompileStmt(stmt)
			if err != nil {
				return nil, err
			}
			res = append(res, stmtIr)
		}
		return BlockStmt{res}, nil

	case ast.LetStmt:
		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		return LetStmt{
			Name: s.Name,
			Expr: expr,
		}, nil

	case ast.IfStmt:
		cond, err := c.CompileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		body, err := c.CompileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		if s.Alternative == nil {
			return IfStmt{
				Condition: cond,
				Body:      body,
			}, nil
		}

		alternative, err := c.CompileStmt(*s.Alternative)
		if err != nil {
			return nil, err
		}

		return IfStmt{
			Condition:   cond,
			Body:        body,
			Alternative: &alternative,
		}, nil

	case ast.ReturnStmt:
		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		return ReturnStmt{Expr: expr}, nil

	case ast.ExprStmt:
		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		return ExpressionStmt{Expr: expr}, nil

	case ast.WhileStmt:
		cond, err := c.CompileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		body, err := c.CompileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		return LoopStmt{
			BlockStmt{
				[]Statement{
					IfStmt{
						Condition: UnaryExpr{
							Expr:     cond,
							Operator: UnaryOpNot,
						},
						Body: BreakStmt{},
					},
					body,
				},
			},
		}, nil

	case ast.ForStmt:
		init, err := c.CompileStmt(s.InitStmt)
		if err != nil {
			return nil, err
		}

		cond, err := c.CompileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		incr, err := c.CompileExpr(s.Increment)
		if err != nil {
			return nil, err
		}

		body, err := c.CompileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		return BlockStmt{
			Statements: []Statement{
				init,
				LoopStmt{
					BlockStmt{
						[]Statement{
							IfStmt{
								Condition: UnaryExpr{
									Expr:     cond,
									Operator: UnaryOpNot,
								},
								Body: BreakStmt{},
							},
							body,
							ExpressionStmt{
								incr,
							},
						},
					},
				},
			},
		}, nil

	case ast.MatchStmt:
		res := BlockStmt{}

		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		res.Statements = append(res.Statements,
			LetStmt{
				Name: "__$e",
				Expr: expr,
			},
		)

		if len(s.Cases) == 0 {
			return res, nil
		}

		currentCase, err := c.compileMatchCase(s.Cases[0], "__$e")
		if err != nil {
			return nil, err
		}

		for _, m := range s.Cases[1:] {
			ifStmt, err := c.compileMatchCase(m, "__$e")
			if err != nil {
				return nil, err
			}

			ifStmt.Alternative = stmtPointer(currentCase)
			currentCase = ifStmt
		}

		res.Statements = append(res.Statements, currentCase)

		return res, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", s))
	}
}

func stmtPointer(s Statement) *Statement {
	return &s
}

func (c *Compiler) compileMatchCase(m ast.MatchCase, exprName string) (IfStmt, error) {
	var res IfStmt

	expr, err := c.CompileExpr(m.Expr)
	if err != nil {
		return res, err
	}

	body, err := c.CompileStmt(m.Body)
	if err != nil {
		return res, err
	}

	res.Condition = BinaryExpr{
		Operator: BinaryOpEq,
		Operands: []Expr{
			IdentExpr{Name: exprName},
			expr,
		},
	}

	res.Body = body

	return res, nil
}

func (c *Compiler) CompileExpr(e ast.Expr) (Expr, error) {
	switch e := e.(type) {
	case ast.IntExpr:
		return IntExpr{Value: e.Value}, nil

	case ast.FloatExpr:
		return FloatExpr{Value: e.Value}, nil

	case ast.BoolExpr:
		return BoolExpr{Value: e.Value}, nil

	case ast.StringExpr:
		return StringExpr{Value: e.Value}, nil

	case ast.AssignExpr:
		assignee, err := c.CompileExpr(e.Assignee)
		if err != nil {
			return nil, err
		}

		expr, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		switch e := assignee.(type) {
		case IdentExpr:
			return VarAssignExpr{
				Name:  e.Name,
				Value: expr,
			}, nil
		case PostFixExpr:
			assignee := e.Ops[:len(e.Ops)-1]
			if err != nil {
				return nil, err
			}

			idx, ok := e.Ops[len(e.Ops)-1].(IndexOp)
			if !ok {
				return nil, errors.New("invalid lhs of assignment")
			}

			return IdxAssignExpr{
				Assignee: PostFixExpr{Ops: assignee, Expr: e.Expr},
				Index:    idx.Index,
				Value:    expr,
			}, nil

		default:
			return nil, errors.New("invalid lhs of assignment")
		}

	case ast.ArrayExpr:
		var res []Expr
		for _, expr := range e.Exprs {
			expr, err := c.CompileExpr(expr)
			if err != nil {
				return nil, err
			}
			res = append(res, expr)
		}
		return ArrayExpr{Exprs: res}, nil

	case ast.ObjectExpr:
		res := map[string]Expr{}

		for key, value := range e.KVs {
			expr, err := c.CompileExpr(value)
			if err != nil {
				return nil, err
			}
			res[key] = expr
		}

		return ObjectExpr{KVs: res}, nil

	case ast.IdentExpr:
		return IdentExpr{Name: e.Name}, nil

	case ast.BinaryExpr:
		var res []Expr
		for _, expr := range e.Operands {
			expr, err := c.CompileExpr(expr)
			if err != nil {
				return nil, err
			}
			res = append(res, expr)
		}

		if e.Operator == "|" {
			return c.compilePipeExpr(res)
		}

		return BinaryExpr{
			Operator: c.getBinaryOp(e.Operator),
			Operands: res,
		}, nil

	case ast.UnaryExpr:
		expr, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		return UnaryExpr{
			Operator: c.getUnaryOp(e.Operator),
			Expr:     expr,
		}, nil

	case ast.PostFixExpr:
		expr, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		var res []PostFixOp

		for _, op := range e.Ops {
			switch op := op.(type) {
			case ast.DotOp:
				res = append(res, IndexOp{
					Index: StringExpr{Value: op.Index},
				})

			case ast.IndexOp:
				idx, err := c.CompileExpr(op.Index)
				if err != nil {
					return nil, err
				}

				res = append(res, IndexOp{
					Index: idx,
				})

			case ast.CallOp:
				var args []Expr
				for _, arg := range op.Args {
					expr, err := c.CompileExpr(arg)
					if err != nil {
						return nil, err
					}

					args = append(args, expr)
				}

				res = append(res, CallOp{
					Args: args,
				})
			}
		}

		return PostFixExpr{
			Expr: expr,
			Ops:  res,
		}, nil

	case ast.FunctionExpr:
		body, err := c.CompileStmt(e.Body)
		if err != nil {
			return nil, err
		}

		return FunctionExpr{
			Params: e.Params,
			Body:   body,
		}, nil

	case ast.LambdaExpr:
		body, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		return FunctionExpr{
			Params: e.Params,
			Body:   ReturnStmt{Expr: body},
		}, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", e))
	}
}

func (c *Compiler) getBinaryOp(op string) BinaryOp {
	switch op {
	case "+":
		return BinaryOpAdd
	case "-":
		return BinaryOpSub
	case "*":
		return BinaryOpMul
	case "/":
		return BinaryOpDiv
	case "%":
		return BinaryOpMod
	case "==":
		return BinaryOpEq
	case "!=":
		return BinaryOpNeq
	case "<":
		return BinaryOpLt
	case "<=":
		return BinaryOpLte
	case ">":
		return BinaryOpGt
	case ">=":
		return BinaryOpGte
	case "or":
		return BinaryOpOr
	case "and":
		return BinaryOpAnd
	default:
		panic(fmt.Sprintf("unimplemented operator %s", op))
	}
}

func (c *Compiler) getUnaryOp(op string) UnaryOp {
	switch op {
	case "!":
		return UnaryOpNot
	case "-":
		return UnaryOpNegate
	default:
		panic(fmt.Sprintf("unimplemented operator %s", op))
	}
}

var pipeErr = errors.New("right operand of pipe must be a call expression")

func (c *Compiler) compilePipeExpr(exprs []Expr) (Expr, error) {
	left := exprs[0]

	for _, right := range exprs[1:] {
		right, ok := right.(PostFixExpr)
		if !ok {
			return nil, pipeErr
		}

		lastOp, ok := right.Ops[len(right.Ops)-1].(CallOp)
		if !ok {
			return nil, pipeErr
		}

		lastOp.Args = append([]Expr{left}, lastOp.Args...)
		right.Ops[len(right.Ops)-1] = lastOp

		left = right
	}

	return left, nil
}
