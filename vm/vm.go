package vm

import (
	"fmt"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

const MaxStack = 1024

const MaxCallStack = 1024

type Frame struct {
	ip           int
	instructions []opcode.OpCode
	numVars      int
	stackOffset  int
}

type VM struct {
	stack     [MaxStack]value.Value
	callStack [MaxCallStack]*Frame
	constants []value.Value
	curFrame  *Frame

	sp int
	fp int
}

func New(constants []value.Value, mainFrame *compiler.Frame) *VM {
	vm := &VM{
		constants: constants,
		sp:        -1,
		fp:        -1,
	}

	vm.pushFrame(&Frame{
		instructions: mainFrame.Instructions,
		numVars:      len(mainFrame.Vars),
		ip:           0,
	}, 0)

	return vm
}

func (v *VM) Run() {
	for {
		switch v.curFrame.instructions[v.curFrame.ip] {
		case opcode.OP_CONSTANT:
			index := v.curFrame.instructions[v.curFrame.ip+1]
			v.sp++
			v.stack[v.sp] = v.constants[index]
			v.curFrame.ip += 2

		case opcode.OP_LET:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])
			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.curFrame.ip += 2

		case opcode.OP_LOAD:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])
			val := v.stack[index]
			v.sp++
			v.stack[v.sp] = val
			v.curFrame.ip += 2

		case opcode.OP_ASSIGN:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])

			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.curFrame.ip += 2

		case opcode.OP_ECHO:
			value := v.stack[v.sp]
			v.sp--
			fmt.Println(value.String())
			v.curFrame.ip += 1

		case opcode.OP_JUMP:
			v.curFrame.ip++
			v.curFrame.ip += int(v.curFrame.instructions[v.curFrame.ip])

		case opcode.OP_JUMPF:
			v.curFrame.ip++
			operand := v.stack[v.sp]
			v.sp--
			offset := int(v.curFrame.instructions[v.curFrame.ip])

			if operand.IsTruthy() {
				v.curFrame.ip++
			} else {
				v.curFrame.ip += offset
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

			v.curFrame.ip++

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

			v.curFrame.ip++

		case opcode.OP_MUL:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			v.stack[v.sp].SetInt(left.GetInt() * right.GetInt())

			v.curFrame.ip++

		case opcode.OP_MOD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			v.stack[v.sp].SetInt(left.GetInt() % right.GetInt())

			v.curFrame.ip++

		case opcode.OP_EQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			v.stack[v.sp].SetBool(left.GetInt() == right.GetInt())

			v.curFrame.ip++

		case opcode.OP_POP:
			v.sp--
			v.curFrame.ip++

		case opcode.OP_CALL:
			numArgs := int(v.curFrame.instructions[v.curFrame.ip+1])
			callee := v.stack[v.sp]
			v.sp--

			fn := callee.GetFunction()
			frame := &Frame{
				instructions: fn.Instructions,
				numVars:      fn.NumVars,
				ip:           0,
				stackOffset:  v.sp - numArgs + 1,
			}

			v.curFrame.ip += 2

			v.pushFrame(frame, numArgs)

		case opcode.OP_RET:
			val := v.stack[v.sp]
			v.sp = v.curFrame.stackOffset
			v.stack[v.sp] = val
			v.popFrame()

		case opcode.OP_HALT:
			return

		default:
			panic(fmt.Sprintf("unimplemented %s", v.curFrame.instructions[v.curFrame.ip]))
		}
	}
}

func (v *VM) pushFrame(f *Frame, args int) {
	v.fp++

	for range f.numVars - args {
		v.sp++
		val := value.Value{}
		val.SetNil()
		v.stack[v.sp] = val
	}

	v.callStack[v.fp] = f
	v.curFrame = f
}

func (v *VM) popFrame() {
	v.fp--
	v.curFrame = v.callStack[v.fp]
}
