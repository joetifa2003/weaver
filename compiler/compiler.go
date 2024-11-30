package compiler

import (
	"errors"
	"fmt"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

type Compiler struct {
	frames       *ds.Stack[*Frame]
	constants    []value.Value
	functionsIdx []int
	labelCounter int
}

func New() *Compiler {
	nilValue := value.Value{}
	nilValue.SetNil()
	return &Compiler{
		frames: &ds.Stack[*Frame]{},
		constants: []value.Value{
			nilValue,
		},
	}
}

func (c *Compiler) Compile(p ast.Program) (*Frame, []value.Value, error) {
	c.pushFrame()

	for _, s := range p.Statements {
		instructions, err := c.compileStmt(s)
		if err != nil {
			return nil, nil, err
		}
		c.addInstructions(instructions)
	}

	mainFrame := c.popFrame()
	mainFrame.Instructions = append(mainFrame.Instructions, opcode.OP_HALT)
	mainFrame.Instructions = c.optimize(mainFrame.Instructions)
	mainFrame.Instructions = c.handleLabels(mainFrame.Instructions)

	for _, f := range c.functionsIdx {
		fn := c.constants[f].GetFunction()
		fn.Instructions = c.optimize(fn.Instructions)
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

func (c *Compiler) compileStmt(s ast.Statement) ([]opcode.OpCode, error) {
	switch s := s.(type) {
	case ast.BlockStmt:
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

	case ast.EchoStmt:
		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		var instructions []opcode.OpCode
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_ECHO)

		return instructions, nil

	case ast.LetStmt:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_LET)
		instructions = append(instructions, opcode.OpCode(c.defineVar(s.Name)))

		return instructions, nil

	case ast.WhileStmt:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(s.Condition)
		if err != nil {
			return nil, err
		}

		body, err := c.compileStmt(s.Body)
		if err != nil {
			return nil, err
		}

		loopLabel := c.label()
		falseLabel := c.label()

		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(loopLabel))
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_JUMPF, opcode.OpCode(falseLabel))
		instructions = append(instructions, body...)
		instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(loopLabel))
		instructions = append(instructions, opcode.OP_LABEL, opcode.OpCode(falseLabel))

		return instructions, nil

	case ast.IfStmt:
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

	case ast.ExprStmt:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_POP)

		return instructions, nil

	case ast.ReturnStmt:
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

func (c *Compiler) compileExpr(e ast.Expr) ([]opcode.OpCode, error) {
	switch e := e.(type) {
	case ast.BinaryExpr:
		var instructions []opcode.OpCode

		if e.Operator == "|" {
			return c.compilePipeExpr(e)
		}

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

	case ast.UnaryExpr:
		var instructions []opcode.OpCode

		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, c.unaryOperatorOpcode(e.Operator))

		return instructions, nil

	case ast.FunctionExpr:
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

		fnValue := value.Value{}
		fnValue.SetFunction(value.FunctionValue{
			NumVars:      len(frame.Vars),
			Instructions: frame.Instructions,
		})

		constant := c.defineConstant(fnValue)
		c.functionsIdx = append(c.functionsIdx, constant)

		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(constant),
		}, nil

	case ast.IdentExpr:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.OpCode(c.resolveVar(e.Name)),
		}, nil

	case ast.IntExpr:
		value := value.Value{}
		value.SetInt(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.FloatExpr:
		value := value.Value{}
		value.SetFloat(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.StringExpr:
		value := value.Value{}
		value.SetString(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.BoolExpr:
		value := value.Value{}
		value.SetBool(e.Value)

		return []opcode.OpCode{
			opcode.OP_CONST,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.AssignExpr:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(e.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_STORE)
		instructions = append(instructions, opcode.OpCode(c.resolveVar(e.Name)))

		return instructions, nil

	case ast.CallExpr:
		var instructions []opcode.OpCode

		for _, arg := range e.Args {
			expr, err := c.compileExpr(arg)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
		}

		callee, err := c.compileExpr(e.Callee)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, callee...)
		instructions = append(instructions, opcode.OP_CALL, opcode.OpCode(len(e.Args)))

		return instructions, nil

	case ast.ArrayExpr:
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

	case ast.ObjectExpr:
		var instructions []opcode.OpCode
		instructions = append(instructions, opcode.OP_OBJ)

		for key, value := range e.KVs {
			exprInstructions, err := c.compileExpr(value)
			if err != nil {
				return nil, err
			}
			keyInstructions, err := c.compileExpr(ast.StringExpr{Value: key})
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

func (c *Compiler) compilePipeExpr(e ast.BinaryExpr) ([]opcode.OpCode, error) {
	left := e.Operands[0]
	for _, right := range e.Operands[1:] {
		switch right := right.(type) {
		case ast.CallExpr:
			right.Args = append([]ast.Expr{left}, right.Args...)
			left = right
		default:
			return nil, errors.New("right operand of pipe must be a call expression")
		}
	}

	return c.compileExpr(left)
}

func (c *Compiler) binOperatorOpcode(operator string) opcode.OpCode {
	switch operator {
	case "+":
		return opcode.OP_ADD

	case "-":
		return opcode.OP_SUB

	case "*":
		return opcode.OP_MUL

	case "%":
		return opcode.OP_MOD

	case "/":
		return opcode.OP_DIV

	case "==":
		return opcode.OP_EQ

	case "!=":
		return opcode.OP_NEQ

	case "<":
		return opcode.OP_LT

	case "<=":
		return opcode.OP_LTE

	case ">":
		return opcode.OP_GT

	case ">=":
		return opcode.OP_GTE

	case "or":
		return opcode.OP_OR

	case "and":
		return opcode.OP_AND

	case "!":
		return opcode.OP_NOT

	default:
		panic(fmt.Sprintf("unimplemented operator %s", operator))
	}
}

func (c *Compiler) unaryOperatorOpcode(operator string) opcode.OpCode {
	switch operator {
	case "!":
		return opcode.OP_NOT

	default:
		panic(fmt.Sprintf("unimplemented operator %s", operator))
	}
}

func (c *Compiler) defineVar(name string) int {
	return c.frames.Peek().defineVar(name)
}

func (c *Compiler) resolveVar(name string) int {
	return c.frames.Peek().resolve(name)
}

func (c *Compiler) defineConstant(v value.Value) int {
	c.constants = append(c.constants, v)
	return len(c.constants) - 1
}

func (c *Compiler) pushFrame() {
	c.frames.Push(NewFrame())
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
