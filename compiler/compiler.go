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

	case ast.AssignStmt:
		var instructions []opcode.OpCode
		expr, err := c.compileExpr(s.Expr)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_ASSIGN)
		instructions = append(instructions, opcode.OpCode(c.resolveVar(s.Name)))

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

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_JUMPF, opcode.OpCode(len(body)+3))
		instructions = append(instructions, body...)
		instructions = append(instructions, opcode.OP_JUMP, opcode.OpCode(-(len(body) + 2 + len(expr) + 1)))

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

		instructions = append(instructions, expr...)
		instructions = append(instructions, opcode.OP_JUMPF, opcode.OpCode(len(body)+1))
		instructions = append(instructions, body...)

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
		for _, operand := range e.Operands {
			expr, err := c.compileExpr(operand)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, expr...)
		}

		for range len(e.Operands) - 1 {
			switch e.Operator {
			case "+":
				instructions = append(instructions, opcode.OP_ADD)
			case "*":
				instructions = append(instructions, opcode.OP_MUL)
			case "<":
				instructions = append(instructions, opcode.OP_LT)
			case "%":
				instructions = append(instructions, opcode.OP_MOD)
			case "/":
				instructions = append(instructions, opcode.OP_DIV)
			case "==":
				instructions = append(instructions, opcode.OP_EQ)
			default:
				panic(fmt.Sprintf("unimplemented operator %s", e.Operator))
			}
		}

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
			opcode.OP_CONSTANT,
			opcode.OpCode(0),
			opcode.OP_RET,
		})

		frame := c.popFrame()

		fnValue := value.Value{}
		fnValue.SetFunction(value.FunctionValue{
			NumVars:      len(frame.Vars),
			Instructions: frame.Instructions,
		})

		return []opcode.OpCode{
			opcode.OP_CONSTANT,
			opcode.OpCode(c.defineConstant(fnValue)),
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
			opcode.OP_CONSTANT,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.FloatExpr:
		value := value.Value{}
		value.SetFloat(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONSTANT,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

	case ast.StringExpr:
		value := value.Value{}
		value.SetString(e.Value)
		return []opcode.OpCode{
			opcode.OP_CONSTANT,
			opcode.OpCode(c.defineConstant(value)),
		}, nil

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
