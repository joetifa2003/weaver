package vm

import (
	"fmt"

	"github.com/joetifa2003/weaver/opcode"
)

const MaxStack = 1024

const MaxCallStack = 1024

type Frame struct {
	ip int
	// TODO: Embed the function value in the frame
	instructions []opcode.OpCode
	numVars      int
	freeVars     []Value
	stackOffset  int
	haltAfter    bool
}

type VM struct {
	stack     [MaxStack]Value
	callStack [MaxCallStack]*Frame
	constants []Value
	curFrame  *Frame

	sp int
	fp int
}

func New(constants []Value, instructions []opcode.OpCode, vars int) *VM {
	vm := &VM{
		constants: constants,
		sp:        -1,
		fp:        -1,
	}

	vm.pushFrame(&Frame{
		instructions: instructions,
		numVars:      vars,
		ip:           0,
	}, 0)

	return vm
}

func (v *VM) RunFunction(f Value, args ...Value) Value {
	fn := f.GetFunction()
	for _, arg := range args {
		v.sp++
		v.stack[v.sp] = arg
	}
	v.pushFrame(&Frame{
		instructions: fn.Instructions,
		numVars:      fn.NumVars,
		freeVars:     fn.FreeVars,
		ip:           0,
		haltAfter:    true,
		stackOffset:  v.sp - len(args) + 1,
	}, len(args))
	v.Run()

	retVal := v.stack[v.sp]
	v.sp--

	return retVal
}

func (v *VM) Run() {
	for {
		switch v.curFrame.instructions[v.curFrame.ip] {
		case opcode.OP_CONST:
			index := v.curFrame.instructions[v.curFrame.ip+1]
			v.sp++
			v.stack[v.sp] = v.constants[index]
			v.curFrame.ip += 2

		case opcode.OP_FUNC:
			constantIndex := int(v.curFrame.instructions[v.curFrame.ip+1])
			freeVarsCount := int(v.curFrame.instructions[v.curFrame.ip+2])
			var freeVars []Value

			for range freeVarsCount {
				freeVars = append(freeVars, v.stack[v.sp])
				v.sp--
			}

			fn := *v.constants[constantIndex].GetFunction()
			fn.FreeVars = freeVars
			v.sp++
			v.stack[v.sp].SetFunction(fn)

			v.curFrame.ip += 3

		case opcode.OP_LOAD:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])
			val := v.stack[index]
			v.sp++
			v.stack[v.sp] = val
			v.curFrame.ip += 2

		case opcode.OP_LOAD_FREE:
			index := int(v.curFrame.instructions[v.curFrame.ip+1])
			val := v.curFrame.freeVars[index]

			v.sp++
			v.stack[v.sp] = val
			v.curFrame.ip += 2

		case opcode.OP_LET:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])
			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.curFrame.ip += 2

		case opcode.OP_STORE:
			index := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])

			v.stack[index] = v.stack[v.sp]
			v.curFrame.ip += 2

		case opcode.OP_ECHO:
			value := v.stack[v.sp]
			v.sp--
			fmt.Println(value.String())
			v.curFrame.ip += 1

		case opcode.OP_LABEL:
			v.curFrame.ip += 2

		case opcode.OP_JUMP:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])
			v.curFrame.ip = newIp

		case opcode.OP_JUMPF:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])
			operand := v.stack[v.sp]
			v.sp--

			if operand.IsTruthy() {
				v.curFrame.ip += 2
			} else {
				v.curFrame.ip = newIp
			}

		case opcode.OP_LT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case ValueTypeInt:
				switch right.VType {
				case ValueTypeInt:
					v.stack[v.sp].SetBool(left.GetInt() < right.GetInt())
				default:
					panic(fmt.Sprintf("illegal operation %s < %s", left, right))
				}
			default:
				panic(fmt.Sprintf("illegal operation %s < %s", left, right))
			}

			v.curFrame.ip++

		case opcode.OP_LTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case ValueTypeInt:
				switch right.VType {
				case ValueTypeInt:
					v.stack[v.sp].SetBool(left.GetInt() <= right.GetInt())
				default:
					panic(fmt.Sprintf("illegal operation %s < %s", left, right))
				}
			default:
				panic(fmt.Sprintf("illegal operation %s < %s", left, right))
			}

			v.curFrame.ip++

		case opcode.OP_GT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case ValueTypeInt:
				switch right.VType {
				case ValueTypeInt:
					v.stack[v.sp].SetBool(left.GetInt() > right.GetInt())
				default:
					panic(fmt.Sprintf("illegal operation %s < %s", left, right))
				}
			default:
				panic(fmt.Sprintf("illegal operation %s < %s", left, right))
			}

			v.curFrame.ip++

		case opcode.OP_GTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			switch left.VType {
			case ValueTypeInt:
				switch right.VType {
				case ValueTypeInt:
					v.stack[v.sp].SetBool(left.GetInt() > right.GetInt())
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
			left.Add(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_OR:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--
			if left.IsTruthy() {
				v.stack[v.sp] = left
			} else {
				v.stack[v.sp] = right
			}

			v.curFrame.ip++

		case opcode.OP_AND:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			if left.IsTruthy() {
				v.stack[v.sp] = right
			} else {
				v.stack[v.sp] = left
			}

			v.curFrame.ip++

		case opcode.OP_OBJ:
			v.sp++
			v.stack[v.sp].SetObject(map[string]Value{})
			v.curFrame.ip++

		case opcode.OP_OPUSH:
			key := v.stack[v.sp]
			v.sp--
			value := v.stack[v.sp]
			v.sp--
			obj := v.stack[v.sp].GetObject()
			obj[key.GetString()] = value

			v.curFrame.ip++

		case opcode.OP_ARRAY:
			v.sp++
			v.stack[v.sp].SetArray(nil)

			v.curFrame.ip++

		case opcode.OP_APUSH:
			val := v.stack[v.sp]
			v.sp--
			arr := v.stack[v.sp].GetArray()
			*arr = append(*arr, val)

			v.curFrame.ip++

		case opcode.OP_INDEX:
			index := v.stack[v.sp]
			arr := v.stack[v.sp-1]
			v.sp--

			switch arr.VType {
			case ValueTypeArray:
				val := (*arr.GetArray())[index.GetInt()]
				v.stack[v.sp] = val
			case ValueTypeObject:
				val := arr.GetObject()[index.GetString()]
				v.stack[v.sp] = val
			}

			v.curFrame.ip++

		case opcode.OP_NOT:
			right := v.stack[v.sp]
			v.stack[v.sp].SetBool(!right.IsTruthy())
			v.curFrame.ip++

		case opcode.OP_SUB:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--
			left.Sub(right, &v.stack[v.sp])

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

		case opcode.OP_NEQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			v.stack[v.sp].SetBool(left.GetInt() == right.GetInt())

			v.curFrame.ip++

		case opcode.OP_POP:
			v.sp--
			v.curFrame.ip++

		case opcode.OP_STORE_IDX:
			idx := v.stack[v.sp]
			assignee := v.stack[v.sp-1]
			val := v.stack[v.sp-2]

			switch assignee.VType {
			case ValueTypeArray:
				(*assignee.GetArray())[idx.GetInt()] = val
			case ValueTypeObject:
				assignee.GetObject()[idx.String()] = val
			}

			v.sp -= 2
			v.stack[v.sp] = assignee
			v.curFrame.ip++

		case opcode.OP_CALL:
			// stack state
			// callee        <- return address
			// args1         <- stackOffset
			// args2
			// ...
			// argsN
			numArgs := int(v.curFrame.instructions[v.curFrame.ip+1])
			calleeIdx := v.sp - numArgs
			argsBegin := calleeIdx + 1
			callee := v.stack[calleeIdx]

			switch callee.VType {
			case ValueTypeFunction:
				fn := callee.GetFunction()
				frame := &Frame{
					instructions: fn.Instructions,
					numVars:      fn.NumVars,
					freeVars:     fn.FreeVars,
					ip:           0,
					stackOffset:  argsBegin,
				}
				v.curFrame.ip += 2
				v.pushFrame(frame, numArgs)
			case ValueTypeNativeFunction:
				fn := callee.GetNativeFunction()
				args := v.stack[argsBegin : argsBegin+numArgs]
				r := fn(v, args...)
				v.sp = calleeIdx
				v.stack[v.sp] = r

				v.curFrame.ip += 2
			}

		case opcode.OP_RET:
			val := v.stack[v.sp]
			v.sp = v.curFrame.stackOffset - 1
			v.stack[v.sp] = val
			heltAfter := v.curFrame.haltAfter
			v.popFrame()
			if heltAfter {
				return
			}

		case opcode.OP_CONST_LET:
			constantIdx := int(v.curFrame.instructions[v.curFrame.ip+1])
			variableIdx := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+2])
			v.stack[variableIdx] = v.constants[constantIdx]
			v.curFrame.ip += 3

		case opcode.OP_LOAD_CONST_ADD_LET:
			constantIdx := int(v.curFrame.instructions[v.curFrame.ip+1])
			variableIdx := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+2])
			variableIdx2 := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+3])
			v.constants[constantIdx].Add(v.stack[variableIdx], &v.stack[variableIdx2])
			v.curFrame.ip += 4

		case opcode.OP_LOAD_CONST_ADD:
			constantIdx := int(v.curFrame.instructions[v.curFrame.ip+1])
			variableIdx := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+2])
			v.sp++
			v.constants[constantIdx].Add(v.stack[variableIdx], &v.stack[v.sp])
			v.curFrame.ip += 3

		case opcode.OP_CONST_ADD:
			constantIdx := int(v.curFrame.instructions[v.curFrame.ip+1])
			right := v.constants[constantIdx]
			left := v.stack[v.sp]
			left.Add(right, &v.stack[v.sp])
			v.curFrame.ip += 2

		case opcode.OP_LOAD_LOAD_LT:
			rightIdx := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+1])
			leftIdx := v.curFrame.stackOffset + int(v.curFrame.instructions[v.curFrame.ip+2])
			v.sp++
			v.stack[v.sp].SetBool(v.stack[leftIdx].GetInt() < v.stack[rightIdx].GetInt())
			v.curFrame.ip += 3

		case opcode.OP_HALT:
			return

		default:
			panic(fmt.Sprintf("unimplemented %d", v.curFrame.instructions[v.curFrame.ip]))
		}
	}
}

func (v *VM) getFunctionFrame(f Value, numArgs int) *Frame {
	fn := f.GetFunction()
	return &Frame{
		instructions: fn.Instructions,
		numVars:      fn.NumVars,
		ip:           0,
		stackOffset:  v.sp - numArgs + 1,
	}
}

func (v *VM) pushFrame(f *Frame, args int) {
	v.fp++

	for range f.numVars - args {
		v.sp++
		val := Value{}
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
