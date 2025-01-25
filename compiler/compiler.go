package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/vm"
)

type Compiler struct {
	frames              *ds.Stack[*Frame]
	constants           []vm.Value
	functionsIdx        []int
	labelCounter        int
	loopContext         *ds.Stack[loopContext]
	optimizationEnabled bool
}

type loopContext struct {
	loopStart int
	loopEnd   int
}

type CompilerOption func(*Compiler)

func WithOptimization(enabled bool) CompilerOption {
	return func(c *Compiler) {
		c.optimizationEnabled = enabled
	}
}

func New(opts ...CompilerOption) *Compiler {
	nilValue := vm.Value{}
	nilValue.SetNil()

	c := &Compiler{
		frames:      &ds.Stack[*Frame]{},
		loopContext: &ds.Stack[loopContext]{},
		constants: []vm.Value{
			nilValue,
		},
		optimizationEnabled: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Compiler) Compile(p []ir.Statement) (*Frame, []vm.Value, error) {
	c.pushFrame()

	for _, s := range p {
		instructions, err := c.compileStmt(s)
		if err != nil {
			return nil, nil, err
		}
		c.addInstructions(instructions)
	}

	mainFrame := c.popFrame()
	mainFrame.Instructions = append(mainFrame.Instructions, opcode.OP_HALT)
	if c.optimizationEnabled {
		mainFrame.Instructions = c.optimize(mainFrame.Instructions)
	}
	mainFrame.Instructions = c.handleLabels(mainFrame.Instructions)

	for _, f := range c.functionsIdx {
		fn := c.constants[f].GetFunction()
		if c.optimizationEnabled {
			fn.Instructions = c.optimize(fn.Instructions)
		}
		fn.Instructions = c.handleLabels(fn.Instructions)
	}

	return mainFrame, c.constants, nil
}

func (c *Compiler) handleLabels(instructions []opcode.OpCode) []opcode.OpCode {
	var newInstructions []opcode.OpCode

	labels := map[opcode.OpCode]opcode.OpCode{} // label idx => instruction idx

	for _, instr := range opcode.OpCodeIterator(instructions) {
		if instr.Op == opcode.OP_LABEL {
			labels[instr.Args[0]] = opcode.OpCode(instr.Addr) + 2
		}
	}

	for _, instr := range opcode.OpCodeIterator(instructions) {
		switch instr.Op {
		case opcode.OP_JUMP:
			instr.Args[0] = labels[instr.Args[0]]

		case opcode.OP_JUMPF:
			instr.Args[0] = labels[instr.Args[0]]
		}

		newInstructions = append(newInstructions, instr.Op)
		newInstructions = append(newInstructions, instr.Args...)
	}

	return newInstructions
}

func (c *Compiler) compileStmt(s ir.Statement) ([]opcode.OpCode, error) {
	switch s := s.(type) {
	case ir.BlockStmt:
		c.beginBlock()

		instructions := []opcode.OpCode{}
		for _, stmt := range s.Statements {
			stmtInstructions, err := c.compileStmt(stmt)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, stmtInstructions...)
		}

		c.endBlock()

		return instructions, nil

	case ir.LetStmt:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_LET)
		instructions = append(instructions, opcode.OpCode(c.defineVar(s.Name)))

		return instructions, nil

	case ir.LoopStmt:
		loopBegin := c.label()
		loopEnd := c.label()

		c.beginLoop(loopBegin, loopEnd)

		var instructions []opcode.OpCode
		body, err := c.compileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		c.endLoop()

		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(loopBegin))
		instructions = append(instructions, body...)
		instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(loopBegin))
		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(loopEnd))

		return instructions, nil

	case ir.BreakStmt:
		loop := c.currentLoop()
		return []opcode.OpCode{
			opcode.OP_JUMP,
			opcode.OpCode(loop.loopEnd),
		}, nil

	case ir.ContinueStmt:
		loop := c.currentLoop()
		return []opcode.OpCode{
			opcode.OP_JUMP,
			opcode.OpCode(loop.loopStart),
		}, nil

	case ir.IfStmt:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		body, err := c.compileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		if s.Alternative == nil {
			falseLabel := c.label()

			instructions = append(instructions, expr...)
			instructions = append(instructions, opcode.OP_JUMPF, opcode.OpCode(falseLabel))
			instructions = append(instructions, body...)
			instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(falseLabel))
		} else {
			lbl1 := c.label()
			lbl2 := c.label()

			alternative, err := c.compileStmt(*s.Alternative)
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, expr...)
			instructions = append(instructions, opcode.OP_JUMPF, opcode.OpCode(lbl1))
			instructions = append(instructions, body...)
			instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(lbl2))
			instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(lbl1))
			instructions = append(instructions, alternative...)
			instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(lbl2))
		}

		return instructions, nil

	case ir.ExpressionStmt:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_POP)

		return instructions, nil

	case ir.ReturnStmt:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_RET)

		return instructions, nil

	default:
		panic(fmt.Sprintf("Unimplemented: %T", s))
	}
}

func (c *Compiler) compileExpr(e ir.Expr) ([]opcode.OpCode, error) {
	switch e := e.(type) {
	case ir.BinaryExpr:
		var instructions []opcode.OpCode

		for _, operand := range e.Operands {
			expr, err := c.compileExpr(operand)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
		}
		for range len(e.Operands) - 1 {
			instructions = append(instructions, c.binOperatorOpcode(e.Operator))
		}

		return instructions, nil

	case ir.PostFixExpr:
		var instructions []opcode.OpCode

		lhs, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, lhs...)

		for _, op := range e.Ops {
			switch op := op.(type) {
			case ir.IndexOp:
				expr, err := c.compileExpr(op.Index)
				if err != nil {
					return nil, err
				}
				instructions = append(instructions, expr...)
				instructions = append(instructions, opcode.OP_INDEX)

			case ir.CallOp:
				for _, arg := range op.Args {
					expr, err := c.compileExpr(arg)
					if err != nil {
						return nil, err
					}
					instructions = append(instructions, expr...)
				}
				instructions = append(instructions, opcode.OP_CALL, opcode.OpCode(len(op.Args)))
			}
		}

		return instructions, nil

	case ir.UnaryExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, c.unaryOperatorOpcode(e.Operator))

		return instructions, nil

	case ir.FunctionExpr:
		c.pushFrame()

		for _, param := range e.Params {
			c.defineVar(param)
		}

		body, err := c.compileStmt(e.Body)
		if err != nil {
			return nil, err
		}
		c.addInstructions(body)

		// default return
		c.addInstructions([]opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(0),
			opcode.OP_RET,
		})

		frame := c.popFrame()

		fnValue := vm.Value{}
		fnValue.SetFunction(vm.FunctionValue{
			NumVars:      len(frame.Vars),
			Instructions: frame.Instructions,
		})

		constant := c.defineConstant(fnValue)
		c.functionsIdx = append(c.functionsIdx, constant)

		var instructions []opcode.OpCode

		for _, freeVar := range frame.FreeVars {
			instructions = append(instructions, freeVar.Parent.load()...)
		}
		instructions = append(instructions,
			opcode.OP_FUNC,
			opcode.OpCode(constant),
			opcode.OpCode(len(frame.FreeVars)),
		)

		return instructions, nil

	case ir.IdentExpr:
		v, err := c.resolveVar(e.Name)
		if err == nil {
			return v.load(), nil
		}

		if f, ok := builtInFunctions[e.Name]; ok {
			fVal := vm.Value{}
			fVal.SetNativeFunction(f)

			return []opcode.OpCode{
				opcode.OP_CONST,
				opcode.OpCode(c.defineConstant(fVal)),
			}, nil
		}

		return nil, err

	case ir.IntExpr:
		value := vm.Value{}
		value.SetInt(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.FloatExpr:
		value := vm.Value{}
		value.SetFloat(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.StringExpr:
		value := vm.Value{}
		value.SetString(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.BoolExpr:
		value := vm.Value{}
		value.SetBool(e.Value)

		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.VarAssignExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Value)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)

		v, err := c.resolveVar(e.Name)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, v.store()...)

		return instructions, nil

	case ir.IdxAssignExpr:
		var instructions []opcode.OpCode

		assignee, err := c.compileExpr(e.Assignee)
		if err != nil {
			return nil, err
		}

		idx, err := c.compileExpr(e.Index)
		if err != nil {
			return nil, err
		}

		val, err := c.compileExpr(e.Value)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, val...)
		instructions = append(instructions, assignee...)
		instructions = append(instructions, idx...)
		instructions = append(instructions, opcode.OP_STORE_IDX)

		return instructions, nil

	case ir.ArrayExpr:
		var instructions []opcode.OpCode

		instructions = append(instructions, opcode.OP_ARRAY)
		for _, expr := range e.Exprs {
			expr, err := c.compileExpr(expr)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
			instructions = append(instructions, opcode.OP_APUSH)
		}

		return instructions, nil

	case ir.ObjectExpr:
		var instructions []opcode.OpCode
		instructions = append(instructions, opcode.OP_OBJ)

		for key, value := range e.KVs {
			exprInstructions, err := c.compileExpr(value)
			if err != nil {
				return nil, err
			}
			keyInstructions, err := c.compileExpr(ir.StringExpr{Value: key})
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, exprInstructions...)
			instructions = append(instructions, keyInstructions...)
			instructions = append(instructions, opcode.OP_OPUSH)
		}

		return instructions, nil
	}

	panic(fmt.Sprintf("unimplemented %T", e))
}

func (c *Compiler) binOperatorOpcode(operator ir.BinaryOp) opcode.OpCode {
	switch operator {
	case ir.BinaryOpAdd:
		return opcode.OP_ADD

	case ir.BinaryOpSub:
		return opcode.OP_SUB

	case ir.BinaryOpMul:
		return opcode.OP_MUL

	case ir.BinaryOpMod:
		return opcode.OP_MOD

	case ir.BinaryOpDiv:
		return opcode.OP_DIV

	case ir.BinaryOpEq:
		return opcode.OP_EQ

	case ir.BinaryOpNeq:
		return opcode.OP_NEQ

	case ir.BinaryOpLt:
		return opcode.OP_LT

	case ir.BinaryOpLte:
		return opcode.OP_LTE

	case ir.BinaryOpGt:
		return opcode.OP_GT

	case ir.BinaryOpGte:
		return opcode.OP_GTE

	case ir.BinaryOpOr:
		return opcode.OP_OR

	case ir.BinaryOpAnd:
		return opcode.OP_AND

	default:
		panic(fmt.Sprintf("unimplemented operator %s", operator))
	}
}

func (c *Compiler) unaryOperatorOpcode(operator ir.UnaryOp) opcode.OpCode {
	switch operator {
	case ir.UnaryOpNot:
		return opcode.OP_NOT

	default:
		panic(fmt.Sprintf("unimplemented operator %s", operator))
	}
}

func (c *Compiler) defineVar(name string) int {
	return c.frames.Peek().defineVar(name)
}

func (c *Compiler) resolveVar(name string) (*Var, error) {
	return c.frames.Peek().resolve(name)
}

func (c *Compiler) defineConstant(v vm.Value) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Compiler) pushFrame() {
	c.frames.Push(NewFrame(c.frames.Peek()))
}

func (c *Compiler) popFrame() *Frame {
	return c.frames.Pop()
}

func (c *Compiler) addInstructions(instructions []opcode.OpCode) {
	c.frames.Peek().addInstructions(instructions)
}

func (c *Compiler) beginBlock() {
	c.frames.Peek().beginBlock()
}

func (c *Compiler) endBlock() {
	c.frames.Peek().endBlock()
}

func (c *Compiler) label() int {
	cc := c.labelCounter
	c.labelCounter++
	return cc
}

func (c *Compiler) beginLoop(begin int, end int) {
	c.loopContext.Push(loopContext{
		loopStart: begin,
		loopEnd:   end,
	})
}

func (c *Compiler) endLoop() {
	c.loopContext.Pop()
}

func (c *Compiler) currentLoop() loopContext {
	return c.loopContext.Peek()
}
