package vm

import (
	"fmt"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

const MaxStack = 1024 * 4

const MaxCallStack = 1024

type Frame struct {
	*compiler.Frame
	ip int
}

type VM struct {
	constants []value.Value
	stack     [MaxStack]value.Value
	sp        int

	frames [MaxCallStack]*Frame
	fp     int
}

func New(constants []value.Value, mainFrame *compiler.Frame) *VM {
	frames := [MaxCallStack]*Frame{}
	frames[0] = &Frame{mainFrame, 0}
	return &VM{
		constants: constants,
		frames:    frames,
		stack:     [MaxStack]value.Value{},
		sp:        -1,
		fp:        0,
	}
}

func (v *VM) Run() {
	v.initializeFrame()

	for v.currentFrame().ip < len(v.currentFrame().Instructions) {
		switch v.currentFrame().Instructions[v.currentFrame().ip] {
		case opcode.OP_CONSTANT:
			v.incrementIP()
			index := v.currentInstruction()
			v.sp++
			v.stack[v.sp] = v.constants[index]
			v.incrementIP()

		case opcode.OP_LET:
			v.incrementIP()
			index := v.currentInstruction()
			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.incrementIP()

		case opcode.OP_LOAD:
			v.incrementIP()
			index := v.currentInstruction()
			val := v.stack[index]
			v.sp++
			v.stack[v.sp] = val
			v.incrementIP()

		case opcode.OP_ASSIGN:
			v.incrementIP()
			index := v.currentInstruction()
			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.incrementIP()

		case opcode.OP_ECHO:
			value := v.stack[v.sp]
			v.sp--
			fmt.Println(value.String())
			v.incrementIP()

		case opcode.OP_JUMP:
			v.incrementIP()
			v.currentFrame().ip += int(v.currentInstruction())

		case opcode.OP_JUMPF:
			v.incrementIP()
			operand := v.stack[v.sp]
			v.sp--
			offset := int(v.currentInstruction())

			if !operand.IsTruthy() {
				v.currentFrame().ip += offset
			} else {
				v.incrementIP()
			}

		case opcode.OP_LT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case value.ValueTypeInt:
				switch right.VType {
				case value.ValueTypeInt:
					v.stack[v.sp].SetBool(left.GetInt() < right.GetInt())
				default:
					panic(fmt.Sprintf("illegal operation %s < %s", left, right))
				}
			default:
				panic(fmt.Sprintf("illegal operation %s < %s", left, right))
			}

			v.incrementIP()

		case opcode.OP_ADD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case value.ValueTypeInt:
				switch right.VType {
				case value.ValueTypeInt:
					v.stack[v.sp].SetInt(left.GetInt() + right.GetInt())
				case value.ValueTypeFloat:
					v.stack[v.sp].SetFloat(float64(left.GetInt()) + right.GetFloat())
				default:
					panic(fmt.Sprintf("illegal operation %s + %s", left, right))
				}
			case value.ValueTypeFloat:
				switch right.VType {
				case value.ValueTypeInt:
					v.stack[v.sp].SetFloat(left.GetFloat() + float64(right.GetInt()))
				case value.ValueTypeFloat:
					v.stack[v.sp].SetFloat(left.GetFloat() + right.GetFloat())
				default:
					panic(fmt.Sprintf("illegal operation %s + %s", left, right))
				}
			default:
				panic(fmt.Sprintf("illegal operation %s + %s", left, right))
			}

			v.incrementIP()

		case opcode.OP_MUL:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			v.stack[v.sp].SetInt(left.GetInt() * right.GetInt())

			v.incrementIP()

		default:
			panic(fmt.Sprintf("unimplemented %s", v.currentInstruction()))
		}
	}
}

func (v *VM) currentFrame() *Frame {
	return v.frames[v.fp]
}

func (v *VM) incrementIP() {
	v.currentFrame().ip++
}

func (v *VM) currentInstruction() opcode.OpCode {
	return v.currentFrame().Instructions[v.currentFrame().ip]
}

func (v *VM) initializeFrame() {
	f := v.currentFrame()
	for range f.Vars {
		v.sp++
		v.stack[v.sp] = value.Value{}
	}
}
