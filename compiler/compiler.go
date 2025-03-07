package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/builtin"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/vm"
)

type Compiler struct {
	constants           []vm.Value
	optimizationEnabled bool
	loopContext         *ds.Stack[loopContext]
	functionsIdx        []int
	labelCounter        int
	reg                 *builtin.Registry
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

func New(reg *builtin.Registry, opts ...CompilerOption) *Compiler {
	nilValue := vm.Value{}
	nilValue.SetNil()

	c := &Compiler{
		loopContext: &ds.Stack[loopContext]{},
		constants: []vm.Value{
			nilValue,
		},
		optimizationEnabled: true,
		reg:                 reg,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Compiler) Compile(p ir.Program) ([]opcode.OpCode, int, []vm.Value, error) {
	var instructions []opcode.OpCode

	for _, s := range p.Statements {
		r, err := c.compileStmt(s)
		if err != nil {
			return nil, 0, nil, err
		}
		instructions = append(instructions, r...)
	}
	instructions = append(instructions, opcode.OP_HALT)
	if c.optimizationEnabled {
		instructions = c.optimize(instructions)
	}
	instructions = c.handleLabels(instructions)

	for _, f := range c.functionsIdx {
		fn := c.constants[f].GetFunction()
		if c.optimizationEnabled {
			fn.Instructions = c.optimize(fn.Instructions)
		}
		fn.Instructions = c.handleLabels(fn.Instructions)
	}

	return instructions, p.VarCount, c.constants, nil
}

func (c *Compiler) handleLabels(instructions []opcode.OpCode) []opcode.OpCode {
	newInstructions := make([]opcode.OpCode, 0, len(instructions))

	labels := map[opcode.OpCode]opcode.OpCode{} // label idx => instruction idx

	instrIdx := 0
	for _, instr := range opcode.OpCodeIterator(instructions) {
		if instr.Op == opcode.OP_LABEL {
			labels[instr.Args[0]] = opcode.OpCode(instrIdx)
			continue
		}

		instrIdx += 1 + len(instr.Args)
	}

	for _, instr := range opcode.OpCodeIterator(instructions, opcode.OP_LABEL) {
		switch instr.Op {
		case opcode.OP_JUMP:
			instr.Args[0] = labels[instr.Args[0]]

		case opcode.OP_PJUMP_F:
			instr.Args[0] = labels[instr.Args[0]]

		case opcode.OP_PJUMP_T:
			instr.Args[0] = labels[instr.Args[0]]

		case opcode.OP_JUMP_T:
			instr.Args[0] = labels[instr.Args[0]]

		case opcode.OP_JUMP_F:
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
		instructions := []opcode.OpCode{}
		for _, stmt := range s.Statements {
			stmtInstructions, err := c.compileStmt(stmt)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, stmtInstructions...)
		}

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
		var instructions []opcode.OpCode

		if s.IncrementStatement != nil {
			incr, err := c.compileStmt(s.IncrementStatement)
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, incr...)
		}

		loop := c.currentLoop()

		instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(loop.loopStart))

		return instructions, nil

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
			instructions = append(instructions, opcode.OP_PJUMP_F, opcode.OpCode(falseLabel))
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
			instructions = append(instructions, opcode.OP_PJUMP_F, opcode.OpCode(lbl1))
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

	default:
		panic(fmt.Sprintf("Unimplemented: %T", s))
	}
}

func (c *Compiler) compileExpr(e ir.Expr) ([]opcode.OpCode, error) {
	switch e := e.(type) {
	case ir.ReturnExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_RET)

		return instructions, nil

	case ir.BinaryExpr:
		var instructions []opcode.OpCode

		if e.Operator == ir.BinaryOpAnd {
			return c.compileAndExpr(e)
		} else if e.Operator == ir.BinaryOpOr {
			return c.compileOrExpr(e)
		}

		firstExpr, err := c.compileExpr(e.Operands[0])
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, firstExpr...)

		for _, operand := range e.Operands[1:] {
			expr, err := c.compileExpr(operand)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
			instructions = append(instructions, c.binOperatorOpcode(e.Operator))
		}

		return instructions, nil

	case ir.IndexExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)

		idx, err := c.compileExpr(e.Index)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, idx...)
		instructions = append(instructions, opcode.OP_INDEX)

		return instructions, nil

	case ir.CallExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)

		for _, arg := range e.Args {
			arg, err := c.compileExpr(arg)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, arg...)
		}

		instructions = append(instructions, opcode.OP_CALL, opcode.OpCode(len(e.Args)))

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

	case ir.VarIncrementExpr:
		return []opcode.OpCode{
			opcode.OP_INC,
			opcode.OpCode(e.Var.Scope),
			opcode.OpCode(e.Var.Idx),
		}, nil

	case ir.VarDecrementExpr:
		return []opcode.OpCode{
			opcode.OP_DEC,
			opcode.OpCode(e.Var.Scope),
			opcode.OpCode(e.Var.Idx),
		}, nil

	case ir.FrameExpr:
		var frameBodyInstructions []opcode.OpCode
		body, err := c.compileStmt(ir.BlockStmt{
			Statements: e.Body,
		})
		if err != nil {
			return nil, err
		}
		frameBodyInstructions = append(frameBodyInstructions, body...)
		// default return
		frameBodyInstructions = append(frameBodyInstructions,
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(0),
			opcode.OP_RET,
		)

		fnValue := vm.Value{}
		fnValue.SetFunction(vm.FunctionValue{
			NumVars:      e.VarCount,
			Instructions: frameBodyInstructions,
		})

		constant := c.defineConstant(fnValue)
		c.functionsIdx = append(c.functionsIdx, constant)

		var instructions []opcode.OpCode

		for _, freeVar := range helpers.ReverseIter(e.FreeVars) {
			instructions = append(instructions, c.loadVar(freeVar, opcode.OP_UPGRADE_REF)...)
		}
		instructions = append(instructions,
			opcode.OP_FUNC,
			opcode.OpCode(constant),
			opcode.OpCode(len(e.FreeVars)),
		)

		return instructions, nil

	case ir.LoadExpr:
		return c.loadVar(e.Var, opcode.OP_LOAD), nil

	case ir.BuiltInExpr:
		val, ok := c.reg.ResolveFunc(e.Name)
		if !ok {
			return nil, fmt.Errorf("unknown built-in function %s", e.Name)
		}

		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(val)),
		}, nil

	case ir.ModuleLoadExpr:
		val, ok := c.reg.ResolveModule(e.Name)
		if !ok {
			return nil, fmt.Errorf("unknown module %s", e.Name)
		}

		mod := val.GetModule()
		fn, ok := mod[e.Load]
		if !ok {
			return nil, fmt.Errorf("unknown module %s:%s", e.Name, e.Load)
		}

		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(fn)),
		}, nil

	case ir.IntExpr:
		value := vm.Value{}
		value.SetNumber(float64(e.Value))
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.FloatExpr:
		value := vm.Value{}
		value.SetNumber(e.Value)
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.StringExpr:
		value := vm.Value{}
		value.SetString(e.Value)
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.NilExpr:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(0),
		}, nil

	case ir.BoolExpr:
		value := vm.Value{}
		value.SetBool(e.Value)

		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeConst,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ir.VarAssignExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Value)
		if err != nil {
			return nil, err
		}

		switch e.Value.(type) {
		case ir.FrameExpr:
			instructions = append(instructions, opcode.OP_EMPTY_FUNC)
			instructions = append(instructions, c.storeVar(e.Var)...)
			instructions = append(instructions, opcode.OP_POP)
			instructions = append(instructions, expr...)
			instructions = append(instructions, opcode.OP_FUNC_LET, opcode.OpCode(e.Var.Scope), opcode.OpCode(e.Var.Idx))
			return instructions, nil
		default:
			instructions = append(instructions, expr...)
			instructions = append(instructions, c.storeVar(e.Var)...)
			return instructions, nil
		}

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

	case ir.IfExpr:
		var instructions []opcode.OpCode

		falseLabel := c.label()
		trueLabel := c.label()

		cond, err := c.compileExpr(e.Condition)
		if err != nil {
			return nil, err
		}
		trueExpr, err := c.compileExpr(e.TrueExpr)
		if err != nil {
			return nil, err
		}
		falseExpr, err := c.compileExpr(e.FalseExpr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, cond...)
		instructions = append(instructions, opcode.OP_PJUMP_F, opcode.OpCode(falseLabel))
		instructions = append(instructions, trueExpr...)
		instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(trueLabel))
		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(falseLabel))
		instructions = append(instructions, falseExpr...)
		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(trueLabel))

		return instructions, nil
	}

	panic(fmt.Sprintf("unimplemented %T", e))
}

func (c *Compiler) compileAndExpr(e ir.BinaryExpr) ([]opcode.OpCode, error) {
	var instructions []opcode.OpCode
	endLabel := c.label()

	for i, operand := range e.Operands {
		expr, err := c.compileExpr(operand)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		if i != len(e.Operands)-1 {
			instructions = append(instructions, opcode.OP_JUMP_F, opcode.OpCode(endLabel))
			instructions = append(instructions, opcode.OP_POP)
		}
	}

	instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(endLabel))

	return instructions, nil
}

func (c *Compiler) compileOrExpr(e ir.BinaryExpr) ([]opcode.OpCode, error) {
	var instructions []opcode.OpCode
	endLabel := c.label()

	for i, operand := range e.Operands {
		expr, err := c.compileExpr(operand)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		if i != len(e.Operands)-1 {
			instructions = append(instructions, opcode.OP_JUMP_T, opcode.OpCode(endLabel))
			instructions = append(instructions, opcode.OP_POP)
		}
	}

	instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(endLabel))

	return instructions, nil
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

	case ir.BinaryOpAnd:
		// nothing, handled in compileAndExpr
		return 0

	case ir.BinaryOpOr:
		// nothing, handled in compileOrExpr
		return 0

	default:
		panic(fmt.Sprintf("unimplemented operator %d", operator))
	}
}

func (c *Compiler) unaryOperatorOpcode(operator ir.UnaryOp) opcode.OpCode {
	switch operator {
	case ir.UnaryOpNot:
		return opcode.OP_NOT

	case ir.UnaryOpNegate:
		return opcode.OP_NEG

	default:
		panic(fmt.Sprintf("unimplemented operator %d", operator))
	}
}

func (c *Compiler) defineConstant(v vm.Value) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
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

func (c *Compiler) loadVar(v ir.Var, op opcode.OpCode) []opcode.OpCode {
	switch v.Scope {
	case ir.VarScopeLocal:
		return []opcode.OpCode{
			op,
			opcode.ScopeTypeLocal,
			opcode.OpCode(v.Idx),
		}
	case ir.VarScopeFree:
		return []opcode.OpCode{
			op,
			opcode.ScopeTypeFree,
			opcode.OpCode(v.Idx),
		}
	default:
		panic(fmt.Sprintf("unknown scope %d", v.Scope))
	}
}

func (c *Compiler) storeVar(v ir.Var) []opcode.OpCode {
	switch v.Scope {
	case ir.VarScopeLocal:
		return []opcode.OpCode{
			opcode.OP_STORE,
			opcode.OpCode(v.Idx),
		}
	case ir.VarScopeFree:
		return []opcode.OpCode{
			opcode.OP_STORE_FREE,
			opcode.OpCode(v.Idx),
		}
	default:
		panic(fmt.Sprintf("unknown scope %d", v.Scope))
	}
}
