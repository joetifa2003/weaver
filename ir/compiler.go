package ir

import (
	"errors"
	"fmt"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
)

type Compiler struct {
	blocks *ds.Stack[*basicBlock]
}

func NewCompiler() *Compiler {
	return &Compiler{
		blocks: ds.NewStack[*basicBlock](),
	}
}

func (c *Compiler) beginBlock() *basicBlock {
	b := &basicBlock{}
	b.parent = c.currentBlock()
	b.idx = c.blocks.Len()
	c.blocks.Push(b)
	return b
}

func (c *Compiler) endBlock() BlockStmt {
	basic := c.blocks.Pop()
	stmt := basic.BlockStmt

	for _, v := range basic.vars {
		if v.noInit {
			continue
		}

		stmt.Statements = append(
			[]Statement{
				LetStmt{
					Name: v.id(),
					Expr: NilExpr{},
				},
			},
			stmt.Statements...,
		)
	}

	return stmt
}

func (c *Compiler) currentBlock() *basicBlock {
	return c.blocks.Peek()
}

func (c *Compiler) blockAdd(s Statement) {
	b := c.currentBlock()
	b.Statements = append(b.Statements, s)
}

func (c *Compiler) Compile(p ast.Program) ([]Statement, error) {
	c.beginBlock()

	for _, stmt := range p.Statements {
		stmtIr, err := c.CompileStmt(stmt)
		if err != nil {
			return nil, err
		}
		c.blockAdd(stmtIr)
	}

	res := c.endBlock()

	return res.Statements, nil
}

func (c *Compiler) CompileStmt(s ast.Statement) (Statement, error) {
	switch s := s.(type) {
	case ast.BlockStmt:
		c.beginBlock()
		for _, stmt := range s.Statements {
			stmtIr, err := c.CompileStmt(stmt)
			if err != nil {
				return nil, err
			}
			c.blockAdd(stmtIr)
		}
		return c.endBlock(), nil

	case ast.LetStmt:
		v, err := c.currentBlock().allocate(s.Name)
		if err != nil {
			return nil, err
		}

		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		return ExpressionStmt{
			Expr: VarAssignExpr{
				Name:  v.id(),
				Value: expr,
			},
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

		c.beginBlock()

		c.blockAdd(IfStmt{
			Condition: UnaryExpr{
				Expr:     cond,
				Operator: UnaryOpNot,
			},
			Body: BreakStmt{},
		})

		c.blockAdd(body)

		res := c.endBlock()

		return LoopStmt{res}, nil

	case ast.ForStmt:
		c.beginBlock()

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

		c.blockAdd(init)
		c.beginBlock() // inner block
		c.blockAdd(IfStmt{
			Condition: UnaryExpr{
				Expr:     cond,
				Operator: UnaryOpNot,
			},
			Body: BreakStmt{},
		})
		c.blockAdd(body)
		c.blockAdd(ExpressionStmt{
			incr,
		})
		inner := c.endBlock()

		c.blockAdd(LoopStmt{inner})

		res := c.endBlock()

		return res, nil

	case ast.MatchStmt:
		c.beginBlock()

		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		exprVar, err := c.currentBlock().allocate("")
		if err != nil {
			return nil, err
		}

		c.blockAdd(
			ExpressionStmt{
				Expr: VarAssignExpr{
					Name:  exprVar.id(),
					Value: expr,
				},
			},
		)

		exprIdent := IdentExpr{exprVar.id()}

		if len(s.Cases) == 0 {
			return c.endBlock(), nil
		}

		c.beginBlock()

		currentCase, err := c.compileMatchCase(s.Cases[len(s.Cases)-1], exprIdent)
		if err != nil {
			return nil, err
		}

		c.currentBlock().deallocateAll()

		if s.ElseBody != nil {
			elseStmt, err := c.CompileStmt(*s.ElseBody)
			if err != nil {
				return nil, err
			}

			currentCase.Alternative = stmtPointer(elseStmt)
		}

		for i := len(s.Cases) - 2; i >= 0; i-- {
			m := s.Cases[i]
			ifStmt, err := c.compileMatchCase(m, exprIdent)
			if err != nil {
				return nil, err
			}

			ifStmt.Alternative = stmtPointer(currentCase)
			currentCase = ifStmt

			c.currentBlock().deallocateAll()
		}

		c.blockAdd(currentCase)

		inner := c.endBlock()

		c.blockAdd(inner)

		res := c.endBlock()

		return res, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", s))
	}
}

func (c *Compiler) compileMatchCase(m ast.MatchCase, expr Expr) (IfStmt, error) {
	var res IfStmt

	cond, err := c.compileMatchCondition(m.Condition, expr)
	if err != nil {
		return res, err
	}

	if m.ExtraCond != nil {
		extraCond, err := c.CompileExpr(*m.ExtraCond)
		if err != nil {
			return res, err
		}

		cond = BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				cond,
				extraCond,
			},
		}
	}

	body, err := c.CompileStmt(m.Body)
	if err != nil {
		return res, err
	}

	res.Condition = cond
	res.Body = body

	return res, nil
}

func (c *Compiler) compileMatchCondition(cond ast.MatchCaseCondition, expr Expr) (Expr, error) {
	switch cond := cond.(type) {
	case ast.MatchCaseInt:
		return BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						PostFixExpr{
							IdentExpr{"type"},
							[]PostFixOp{
								CallOp{
									Args: []Expr{
										expr,
									},
								},
							},
						},
						StringExpr{"int"},
					},
				},
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						expr,
						IntExpr{cond.Value},
					},
				},
			},
		}, nil

	case ast.MatchCaseFloat:
		return BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						PostFixExpr{
							IdentExpr{"type"},
							[]PostFixOp{
								CallOp{
									Args: []Expr{
										expr,
									},
								},
							},
						},
						StringExpr{"float"},
					},
				},
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						expr,
						FloatExpr{cond.Value},
					},
				},
			},
		}, nil

	case ast.MatchCaseString:
		return BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						PostFixExpr{
							IdentExpr{"type"},
							[]PostFixOp{
								CallOp{
									Args: []Expr{
										expr,
									},
								},
							},
						},
						StringExpr{"string"},
					},
				},
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						expr,
						StringExpr{cond.Value},
					},
				},
			},
		}, nil

	case ast.MatchCaseIdent:
		v, err := c.currentBlock().allocate(cond.Name)
		if err != nil {
			return nil, err
		}
		return assignOrTrue(v.id(), expr), nil

	case ast.MatchCaseObject:
		isObject := BinaryExpr{
			BinaryOpEq,
			[]Expr{
				PostFixExpr{
					IdentExpr{"type"},
					[]PostFixOp{
						CallOp{[]Expr{expr}},
					},
				},
				StringExpr{"object"},
			},
		}

		hasCorrectLength := BinaryExpr{
			BinaryOpGte,
			[]Expr{
				PostFixExpr{
					IdentExpr{"len"},
					[]PostFixOp{
						CallOp{[]Expr{expr}},
					},
				},
				IntExpr{len(cond.KVs)},
			},
		}

		res := BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				isObject,
				hasCorrectLength,
			},
		}

		for key, cond := range cond.KVs {
			v, err := c.currentBlock().allocate("")
			if err != nil {
				return nil, err
			}

			res.Operands = append(res.Operands,
				assignOrTrue(v.id(), PostFixExpr{
					expr,
					[]PostFixOp{
						IndexOp{Index: StringExpr{Value: key}},
					},
				}),
			)
			child, err := c.compileMatchCondition(cond, IdentExpr{v.id()})
			if err != nil {
				return nil, err
			}
			v.deallocate()

			res.Operands = append(res.Operands, child)
		}

		return res, nil

	case ast.MatchCaseArray:
		isArray := BinaryExpr{
			BinaryOpEq,
			[]Expr{
				PostFixExpr{
					IdentExpr{"type"},
					[]PostFixOp{
						CallOp{[]Expr{expr}},
					},
				},
				StringExpr{"array"},
			},
		}

		hasCorrectLength := BinaryExpr{
			BinaryOpGte,
			[]Expr{
				PostFixExpr{
					IdentExpr{"len"},
					[]PostFixOp{
						CallOp{[]Expr{expr}},
					},
				},
				IntExpr{len(cond.Conditions)},
			},
		}

		res := BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				isArray,
				hasCorrectLength,
			},
		}

		for i, cond := range cond.Conditions {
			v, err := c.currentBlock().allocate("")
			if err != nil {
				return nil, err
			}

			res.Operands = append(res.Operands,
				assignOrTrue(v.id(), PostFixExpr{
					expr,
					[]PostFixOp{
						IndexOp{Index: IntExpr{Value: i}},
					},
				}),
			)
			child, err := c.compileMatchCondition(cond, IdentExpr{v.id()})
			if err != nil {
				return nil, err
			}
			v.deallocate()

			res.Operands = append(res.Operands, child)
		}

		return res, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", cond))
	}
}

func assignOrTrue(name string, value Expr) Expr {
	return BinaryExpr{
		BinaryOpOr,
		[]Expr{
			VarAssignExpr{
				Name:  name,
				Value: value,
			},
			BoolExpr{Value: true},
		},
	}
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
		v, err := c.currentBlock().resolve(e.Name)
		if err != nil {
			// TODO: handle the builting functions at ir level
			return IdentExpr{Name: e.Name}, nil
		}
		return IdentExpr{Name: v.id()}, nil

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
		c.beginBlock()

		for i, param := range e.Params {
			v, err := c.currentBlock().allocate(param)
			if err != nil {
				return nil, err
			}

			v.noInit = true
			e.Params[i] = v.id()
		}

		body, err := c.CompileStmt(e.Body)
		if err != nil {
			return nil, err
		}

		c.blockAdd(body)

		inner := c.endBlock()

		return FunctionExpr{
			Params: e.Params,
			Body:   inner,
		}, nil

	case ast.LambdaExpr:
		c.beginBlock()

		for i, param := range e.Params {
			v, err := c.currentBlock().allocate(param)
			if err != nil {
				return nil, err
			}

			v.noInit = true
			e.Params[i] = v.id()
		}

		body, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		c.blockAdd(ReturnStmt{Expr: body})

		inner := c.endBlock()

		return FunctionExpr{
			Params: e.Params,
			Body:   inner,
		}, nil

	case ast.NilExpr:
		return NilExpr{}, nil

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

func stmtPointer(s Statement) *Statement {
	return &s
}
