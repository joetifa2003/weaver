package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

type Compiler struct {
	Frames    []*Frame
	Constants []value.Value
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) pushFrame() {
	c.Frames = append(c.Frames, newFrame())
}

func (c *Compiler) popFrame() *Frame {
	res := c.Frames[len(c.Frames)-1]
	c.Frames = c.Frames[:len(c.Frames)-1]

	return res
}

func (c *Compiler) currentFrame() *Frame {
	return c.Frames[len(c.Frames)-1]
}

func (c *Compiler) allocRegister() *Reg {
	return c.currentFrame().allocReg()
}

func (c *Compiler) allocVar(name string) *Reg {
	return c.currentFrame().allocVar(name)
}

func (c *Compiler) defineConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)

	return len(c.Constants) - 1
}

func (c *Compiler) addInstructions(instructions []opcode.OpCode) {
	c.currentFrame().addInstructions(instructions)
}

func (c *Compiler) Compile(p *ast.Program) error {
	c.pushFrame()

	for _, s := range p.Statements {
		instructions, err := c.compileStmt(s)
		if err != nil {
			return err
		}

		c.addInstructions(instructions)
	}

	return nil
}

func (c *Compiler) compileStmt(s ast.Statement) ([]opcode.OpCode, error) {
	switch s := s.(type) {
	case *ast.Echo:
		expr, reg, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}

		var instructions []opcode.OpCode

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_ECHO, opcode.OpCode(reg.Index))

		return instructions, nil

	case *ast.Let:
		expr, reg, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		letReg := c.allocVar(s.Name)

		var instructions []opcode.OpCode
		instructions = append(instructions, expr...)
		instructions = append(instructions,
			opcode.OP_LET,
			opcode.OpCode(reg.Index),
			opcode.OpCode(letReg.Index),
		)
		return instructions, nil

	default:
		panic(fmt.Sprintf("Unimplemented: %T", s))
	}
}

func (c *Compiler) compileExpr(e interface{}) ([]opcode.OpCode, *Reg, error) {
	switch e := e.(type) {
	case *ast.Expr:
		return c.compileExpr(e.Equality)

	case *ast.Equality:
		lhs, reg, err := c.compileExpr(e.Left)
		if err != nil {
			return nil, reg, err
		}
		if e.Right == nil {
			return lhs, reg, nil
		}

	case *ast.Comparison:
		lhs, reg, err := c.compileExpr(e.Left)
		if err != nil {
			return nil, reg, err
		}
		if e.Right == nil {
			return lhs, reg, nil
		}

	case *ast.Addition:
		lhs, reg, err := c.compileExpr(e.Left)
		if err != nil {
			return nil, reg, err
		}
		if e.Right == nil {
			return lhs, reg, nil
		}

	case *ast.Multiplication:
		lhs, reg, err := c.compileExpr(e.Left)
		if err != nil {
			return nil, reg, err
		}
		if e.Right == nil {
			return lhs, reg, nil
		}

	case *ast.Unary:
		lhs, reg, err := c.compileExpr(e.Atom)
		if err != nil {
			return nil, reg, err
		}
		if e.Unary == nil {
			return lhs, reg, nil
		}

	case ast.Atom:
		switch a := e.(type) {
		case *ast.Number:
			reg := c.allocRegister()
			return []opcode.OpCode{
				opcode.OP_CONSTANT,
				opcode.OpCode(c.defineConstant(value.NewInt(a.Value))),
				opcode.OpCode(reg.Index),
			}, reg, nil
		}
	}

	panic("unimplemented")
}
