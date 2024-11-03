package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

type Compiler struct {
	frames    *ds.Stack[*Frame]
	constants []value.Value
}

func New() *Compiler {
	return &Compiler{
		frames: &ds.Stack[*Frame]{},
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

	return c.popFrame(), c.constants, nil
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

	default:
		panic(fmt.Sprintf("Unimplemented: %T", s))
	}
}

func (c *Compiler) compileExpr(e ast.Expr) ([]opcode.OpCode, error) {
	switch e := e.(type) {
	case ast.BinaryExpr:
		var instructions []opcode.OpCode
		for _, operand := range e.Operands {
			expr, err := c.compileExpr(operand)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
		}

		for range (len(e.Operands) + 1) / 2 {
			switch e.Operator {
			case "+":
				instructions = append(instructions, opcode.OP_ADD)
			case "*":
				instructions = append(instructions, opcode.OP_MUL)
			default:
				panic(fmt.Sprintf("unimplemented operator %s", e.Operator))
			}
		}

		return instructions, nil

	case ast.IdentExpr:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.OpCode(c.resolveVar(e.Name)),
		}, nil

	case ast.IntExpr:
		return []opcode.OpCode{
			opcode.OP_CONSTANT,
			opcode.OpCode(c.defineConstant(value.NewInt(e.Value))),
		}, nil
	}

	panic(fmt.Sprintf("unimplemented %T", e))
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
