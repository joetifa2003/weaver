package ir

import (
	"errors"
	"fmt"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
)

type Compiler struct {
	frames      *ds.Stack[*frame]
	loopContext *ds.Stack[loopContext]
}

func NewCompiler() *Compiler {
	return &Compiler{
		frames:      ds.NewStack[*frame](),
		loopContext: ds.NewStack[loopContext](),
	}
}

type loopType int

const (
	loopTypeWhile loopType = iota
	loopTypeFor
)

type loopContext struct {
	loopType           loopType
	incrementStatement Statement
}

func (c *Compiler) pushFrame() *frame {
	f := NewFrame(c.currentFrame())
	c.frames.Push(f)
	return f
}

func (c *Compiler) popFrame() *frame {
	return c.frames.Pop()
}

func (c *Compiler) currentFrame() *frame {
	return c.frames.Peek()
}

func (c *Compiler) Compile(p ast.Program) (Program, error) {
	c.pushFrame()

	for _, stmt := range p.Statements {
		stmtIr, err := c.CompileStmt(stmt)
		if err != nil {
			return Program{}, err
		}
		c.currentFrame().pushStmt(stmtIr)
	}

	res := c.popFrame().export()

	return Program{
		VarCount:   res.VarCount,
		Statements: res.Body,
	}, nil
}

func (c *Compiler) CompileStmt(s ast.Statement) (Statement, error) {
	switch s := s.(type) {
	case ast.BreakStmt:
		return BreakStmt{}, nil

	case ast.ContinueStmt:
		if c.loopContext.Peek().loopType != loopTypeFor {
			return ContinueStmt{}, nil
		}

		return ContinueStmt{
			IncrementStatement: c.loopContext.Peek().incrementStatement,
		}, nil

	case ast.BlockStmt:
		b := c.currentFrame().pushBlock()
		for _, stmt := range s.Statements {
			stmtIr, err := c.CompileStmt(stmt)
			if err != nil {
				return nil, err
			}
			b.pushStmt(stmtIr)
		}
		return c.currentFrame().popBlock().export(), nil

	case ast.LetStmt:
		v := c.currentFrame().define(s.Name)

		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		return ExpressionStmt{
			Expr: v.assign(expr),
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

	case ast.ExprStmt:
		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		return ExpressionStmt{Expr: expr}, nil

	case ast.WhileStmt:
		b := c.currentFrame().pushBlock()
		cond, err := c.CompileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		c.loopContext.Push(loopContext{
			loopType: loopTypeWhile,
		})

		body, err := c.CompileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		c.loopContext.Pop()

		b.pushStmt(IfStmt{
			Condition: UnaryExpr{
				Expr:     cond,
				Operator: UnaryOpNot,
			},
			Body: BreakStmt{},
		})

		b.pushStmt(body)

		res := c.currentFrame().popBlock().export()

		return LoopStmt{res}, nil

	case ast.ForStmt:
		outer := c.currentFrame().pushBlock()

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

		c.loopContext.Push(loopContext{
			loopType: loopTypeFor,
			incrementStatement: ExpressionStmt{
				Expr: incr,
			},
		})

		body, err := c.CompileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		c.loopContext.Pop()

		outer.pushStmt(init)

		inner := c.currentFrame().pushBlock()
		inner.pushStmt(IfStmt{
			Condition: UnaryExpr{
				Expr:     cond,
				Operator: UnaryOpNot,
			},
			Body: BreakStmt{},
		})
		inner.pushStmt(body)
		inner.pushStmt(ExpressionStmt{
			incr,
		})
		inner = c.currentFrame().popBlock()

		outer.pushStmt(LoopStmt{inner.export()})

		res := c.currentFrame().popBlock().export()

		return res, nil

	case ast.MatchStmt:
		outer := c.currentFrame().pushBlock()

		expr, err := c.CompileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		exprVar := c.currentFrame().define("")

		outer.pushStmt(exprVar.assignStmt(expr))

		exprIdent := exprVar.load()

		if len(s.Cases) == 0 {
			return c.currentFrame().popBlock().export(), nil
		}

		inner := c.currentFrame().pushBlock()

		currentCase, err := c.compileMatchCase(s.Cases[len(s.Cases)-1], exprIdent)
		if err != nil {
			return nil, err
		}
		inner.freeAll()

		for i := len(s.Cases) - 2; i >= 0; i-- {
			m := s.Cases[i]
			ifStmt, err := c.compileMatchCase(m, exprIdent)
			if err != nil {
				return nil, err
			}

			ifStmt.Alternative = stmtPointer(currentCase)
			currentCase = ifStmt
			inner.freeAll()
		}

		inner.pushStmt(currentCase)

		inner = c.currentFrame().popBlock()

		outer.pushStmt(inner.export())

		outer = c.currentFrame().popBlock()

		return outer.export(), nil

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
	case ast.MatchCaseRange:
		begin, err := c.CompileExpr(cond.Begin)
		if err != nil {
			return nil, err
		}

		end, err := c.CompileExpr(cond.End)
		if err != nil {
			return nil, err
		}

		return BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				BinaryExpr{
					BinaryOpOr,
					[]Expr{
						BinaryExpr{
							BinaryOpEq,
							[]Expr{
								PostFixExpr{
									BuiltInExpr{"type"},
									[]PostFixOp{
										CallOp{
											Args: []Expr{
												expr,
											},
										},
									},
								},
								StringExpr{"number"},
							},
						},
						BinaryExpr{
							BinaryOpEq,
							[]Expr{
								PostFixExpr{
									BuiltInExpr{"type"},
									[]PostFixOp{
										CallOp{
											Args: []Expr{
												expr,
											},
										},
									},
								},
								StringExpr{"number"},
							},
						},
					},
				},
				BinaryExpr{
					BinaryOpGt,
					[]Expr{
						expr,
						begin,
					},
				},
				BinaryExpr{
					BinaryOpLt,
					[]Expr{
						expr,
						end,
					},
				},
			},
		}, nil

	case ast.MatchCaseInt:
		return BinaryExpr{
			BinaryOpAnd,
			[]Expr{
				BinaryExpr{
					BinaryOpEq,
					[]Expr{
						PostFixExpr{
							BuiltInExpr{"type"},
							[]PostFixOp{
								CallOp{
									Args: []Expr{
										expr,
									},
								},
							},
						},
						StringExpr{"number"},
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
							BuiltInExpr{"type"},
							[]PostFixOp{
								CallOp{
									Args: []Expr{
										expr,
									},
								},
							},
						},
						StringExpr{"number"},
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
							BuiltInExpr{"type"},
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
		v := c.currentFrame().define(cond.Name)
		return orTrue(v.assign(expr)), nil

	case ast.MatchCaseObject:
		isObject := BinaryExpr{
			BinaryOpEq,
			[]Expr{
				PostFixExpr{
					BuiltInExpr{"type"},
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
					BuiltInExpr{"len"},
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
			v := c.currentFrame().define("")

			res.Operands = append(res.Operands,
				orTrue(
					v.assign(PostFixExpr{
						expr,
						[]PostFixOp{
							IndexOp{Index: StringExpr{Value: key}},
						},
					}),
				),
			)
			child, err := c.compileMatchCondition(cond, v.load())
			if err != nil {
				return nil, err
			}

			res.Operands = append(res.Operands, child)
		}

		return res, nil

	case ast.MatchCaseArray:
		isArray := BinaryExpr{
			BinaryOpEq,
			[]Expr{
				PostFixExpr{
					BuiltInExpr{"type"},
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
					BuiltInExpr{"len"},
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
			v := c.currentFrame().define("")

			res.Operands = append(res.Operands,
				orTrue(
					v.assign(PostFixExpr{
						expr,
						[]PostFixOp{
							IndexOp{Index: IntExpr{Value: i}},
						},
					}),
				),
			)
			child, err := c.compileMatchCondition(cond, v.load())
			if err != nil {
				return nil, err
			}

			res.Operands = append(res.Operands, child)
		}

		return res, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", cond))
	}
}

func orTrue(value Expr) Expr {
	return BinaryExpr{
		BinaryOpOr,
		[]Expr{
			value,
			BoolExpr{Value: true},
		},
	}
}

func (c *Compiler) CompileExpr(e ast.Expr) (Expr, error) {
	switch e := e.(type) {
	case ast.ReturnExpr:
		expr, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}
		return ReturnExpr{Expr: expr}, nil

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
		case LoadExpr:
			return VarAssignExpr{
				Var:   e.Var,
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
		v, err := c.currentFrame().resolve(e.Name)
		if err == nil {
			return v.load(), nil
		}

		return BuiltInExpr{e.Name}, nil

	case ast.ModuleLoadExpr:
		return ModuleLoadExpr{Name: e.Name, Load: e.Load}, nil

	case ast.BinaryExpr:
		var res []Expr
		for _, expr := range e.Operands {
			expr, err := c.CompileExpr(expr)
			if err != nil {
				return nil, err
			}
			res = append(res, expr)
		}

		if e.Operator == ast.BinaryOpPipe {
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
			default:
				panic(fmt.Sprintf("unimplemented postfix op %T", op))
			}
		}

		return PostFixExpr{
			Expr: expr,
			Ops:  res,
		}, nil

	case ast.VarIncrementExpr:
		v, err := c.currentFrame().resolve(e.Name)
		if err != nil {
			return nil, err
		}

		return VarIncrementExpr{
			Var: v.export(),
		}, nil

	case ast.VarDecrementExpr:
		v, err := c.currentFrame().resolve(e.Name)
		if err != nil {
			return nil, err
		}

		return VarDecrementExpr{
			Var: v.export(),
		}, nil

	case ast.FunctionExpr:
		frame := c.pushFrame()

		for _, param := range e.Params {
			frame.define(param)
		}

		body, err := c.CompileStmt(e.Body)
		if err != nil {
			return nil, err
		}

		frame.pushStmt(body)

		frameExpr := c.popFrame().export()
		frameExpr.ParamsCount = len(e.Params)

		return frameExpr, nil

	case ast.LambdaExpr:
		return c.CompileExpr(
			ast.FunctionExpr{
				Params: e.Params,
				Body: ast.BlockStmt{
					Statements: []ast.Statement{
						ast.ExprStmt{
							Expr: ast.ReturnExpr{
								Expr: e.Expr,
							},
						},
					},
				},
			},
		)

	case ast.NilExpr:
		return NilExpr{}, nil

	case ast.TernaryExpr:
		cond, err := c.CompileExpr(e.Expr)
		if err != nil {
			return nil, err
		}
		trueExpr, err := c.CompileExpr(e.TrueExpr)
		if err != nil {
			return nil, err
		}
		falseExpr, err := c.CompileExpr(e.FalseExpr)
		if err != nil {
			return nil, err
		}

		return IfExpr{
			Condition: cond,
			TrueExpr:  trueExpr,
			FalseExpr: falseExpr,
		}, nil

	default:
		panic(fmt.Sprintf("unimplemented %T", e))
	}
}

func (c *Compiler) getBinaryOp(op ast.BinaryOp) BinaryOp {
	switch op {
	case ast.BinaryOpAdd:
		return BinaryOpAdd
	case ast.BinaryOpSub:
		return BinaryOpSub
	case ast.BinaryOpMul:
		return BinaryOpMul
	case ast.BinaryOpDiv:
		return BinaryOpDiv
	case ast.BinaryOpMod:
		return BinaryOpMod
	case ast.BinaryOpEq:
		return BinaryOpEq
	case ast.BinaryOpNeq:
		return BinaryOpNeq
	case ast.BinaryOpLt:
		return BinaryOpLt
	case ast.BinaryOpLte:
		return BinaryOpLte
	case ast.BinaryOpGt:
		return BinaryOpGt
	case ast.BinaryOpGte:
		return BinaryOpGte
	case ast.BinaryOpOr:
		return BinaryOpOr
	case ast.BinaryOpAnd:
		return BinaryOpAnd
	case ast.BinaryOpPipe:
		// nothing, handled in compilePipeExpr
		return 0
	default:
		panic(fmt.Sprintf("unimplemented operator %s", op))
	}
}

func (c *Compiler) getUnaryOp(op ast.UnaryOp) UnaryOp {
	switch op {
	case ast.UnaryOpNot:
		return UnaryOpNot
	case ast.UnaryOpNegate:
		return UnaryOpNegate
	default:
		panic(fmt.Sprintf("unimplemented operator %s", op))
	}
}

var errPipe = errors.New("right operand of pipe must be a call expression")

func (c *Compiler) compilePipeExpr(exprs []Expr) (Expr, error) {
	left := exprs[0]

	for _, right := range exprs[1:] {
		right, ok := right.(PostFixExpr)
		if !ok {
			return nil, errPipe
		}

		lastOp, ok := right.Ops[len(right.Ops)-1].(CallOp)
		if !ok {
			return nil, errPipe
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
