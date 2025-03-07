package vm

import (
	"fmt"

	"github.com/joetifa2003/weaver/opcode"
)

const MaxStack = 1024

const MaxCallStack = 1024

type Frame struct {
	instructions []opcode.OpCode
	freeVars     []Value
	ip           int
	numVars      int
	stackOffset  int
	returnAddr   int
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
		haltAfter:    true,
		returnAddr:   vars,
	}, 0)

	return vm
}

func (v *VM) RunFunction(f Value, args ...Value) Value {
	fn := f.GetFunction()
	v.sp++
	retAddr := v.sp
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
		stackOffset:  retAddr + 1,
		returnAddr:   retAddr,
	}, len(args))
	v.Run()

	retVal := v.stack[v.sp]
	v.sp--

	return retVal
}

var scopeGettersDeref = [4]func(v *VM, idx int) *Value{
	opcode.ScopeTypeConst: func(v *VM, idx int) *Value {
		return v.constants[idx].deref()
	},
	opcode.ScopeTypeGlobal: func(v *VM, idx int) *Value {
		return v.stack[idx].deref()
	},
	opcode.ScopeTypeLocal: func(v *VM, idx int) *Value {
		return v.stack[v.curFrame.stackOffset+idx].deref()
	},
	opcode.ScopeTypeFree: func(v *VM, idx int) *Value {
		return v.curFrame.freeVars[idx].deref()
	},
}

var scopeGetters = [4]func(v *VM, idx int) *Value{
	opcode.ScopeTypeConst: func(v *VM, idx int) *Value {
		return &v.constants[idx]
	},
	opcode.ScopeTypeGlobal: func(v *VM, idx int) *Value {
		return &v.stack[idx]
	},
	opcode.ScopeTypeLocal: func(v *VM, idx int) *Value {
		return &v.stack[v.curFrame.stackOffset+idx]
	},
	opcode.ScopeTypeFree: func(v *VM, idx int) *Value {
		return &v.curFrame.freeVars[idx]
	},
}

func (v *VM) Run() Value {
	for {
		switch v.curFrame.instructions[v.curFrame.ip] {
		case opcode.OP_UPGRADE_REF:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			idx := int(v.curFrame.instructions[v.curFrame.ip+2])

			val := scopeGetters[scope](v, idx)
			if val.VType != ValueTypeRef {
				val.SetRef(*val)
			}

			v.sp++
			v.stack[v.sp] = *val
			v.curFrame.ip += 3

		case opcode.OP_ADD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--
			left.Add(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_SUB:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--
			left.Sub(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_MUL:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Mul(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_DIV:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Div(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_MOD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Mod(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_EQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Equal(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_NEQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.NotEqual(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_NEG:
			v.stack[v.sp].Negate(&v.stack[v.sp])
			v.curFrame.ip++

		case opcode.OP_LT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.LessThan(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_LTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.LessThanEqual(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_GT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.GreaterThan(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_GTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.GreaterThanEqual(&right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_NOT:
			v.stack[v.sp].SetBool(!v.stack[v.sp].IsTruthy())
			v.curFrame.ip++

		case opcode.OP_EMPTY_FUNC:
			v.sp++
			v.stack[v.sp].SetFunction(FunctionValue{})
			v.curFrame.ip++

		case opcode.OP_FUNC_LET:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])
			f := v.stack[v.sp].GetFunction()

			emptyFunc := scopeGettersDeref[scope](v, index).GetFunction()
			emptyFunc.FreeVars = f.FreeVars
			emptyFunc.Instructions = f.Instructions
			emptyFunc.NumVars = f.NumVars
			v.curFrame.ip += 3

		case opcode.OP_FUNC:
			constantIndex := int(v.curFrame.instructions[v.curFrame.ip+1])
			freeVarsCount := int(v.curFrame.instructions[v.curFrame.ip+2])
			var freeVars []Value

			for range freeVarsCount {
				freeVars = append(freeVars, v.stack[v.sp])
				v.stack[v.sp].SetNil()
				v.sp--
			}

			fn := *v.constants[constantIndex].GetFunction()
			fn.FreeVars = freeVars
			v.sp++
			v.stack[v.sp].SetFunction(fn)

			v.curFrame.ip += 3

		case opcode.OP_INC:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() + 1)
			v.sp++
			v.stack[v.sp] = *v1
			v.curFrame.ip += 3

		case opcode.OP_INC_POP:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() + 1)
			v.curFrame.ip += 3

		case opcode.OP_DEC:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() - 1)
			v.sp++
			v.stack[v.sp] = *v1
			v.curFrame.ip += 3

		case opcode.OP_DEC_POP:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() - 1)
			v.curFrame.ip += 3

		case opcode.OP_LOAD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.sp++
			v.stack[v.sp] = *scopeGettersDeref[scope](v, index)
			v.curFrame.ip += 3

		case opcode.OP_STORE_FREE:
			index := int(v.curFrame.instructions[v.curFrame.ip+1])
			val := v.stack[v.sp]
			v.curFrame.freeVars[index].deref().Set(val)
			v.sp--
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

		case opcode.OP_JUMP:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])
			v.curFrame.ip = newIp

		case opcode.OP_PJUMP_F:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip += 2
			} else {
				v.curFrame.ip = newIp
			}

			v.sp--

		case opcode.OP_PJUMP_T:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip = newIp
			} else {
				v.curFrame.ip += 2
			}

			v.sp--

		case opcode.OP_JUMP_F:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip += 2
			} else {
				v.curFrame.ip = newIp
			}

		case opcode.OP_JUMP_T:
			newIp := int(v.curFrame.instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip = newIp
			} else {
				v.curFrame.ip += 2
			}

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
			val := v.stack[v.sp-1]
			v.sp--

			val.Index(&index, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_POP:
			v.sp--
			v.curFrame.ip++

		case opcode.OP_STORE_IDX:
			idx := v.stack[v.sp]
			assignee := v.stack[v.sp-1]
			val := v.stack[v.sp-2]

			assignee.SetIndex(&idx, val)

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
					returnAddr:   calleeIdx,
				}
				v.curFrame.ip += 2
				v.pushFrame(frame, numArgs)
			case ValueTypeNativeFunction:
				fn := callee.GetNativeFunction()
				args := v.stack[argsBegin : argsBegin+numArgs]
				r, err := fn(v, args)
				if err != nil {
					panic(err)
				}
				v.sp = calleeIdx
				v.stack[v.sp] = r

				v.curFrame.ip += 2
			default:
				panic(fmt.Sprintf("illegal callee type %s", callee.VType))
			}

		case opcode.OP_RET:
			val := v.stack[v.sp]
			v.sp = v.curFrame.returnAddr
			v.stack[v.sp] = val
			haltAfter := v.curFrame.haltAfter
			v.popFrame()
			if haltAfter {
				return val
			}

		case opcode.OP_HALT:
			return Value{}

		case opcode.OP_LOAD_LOAD_ADD:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Add(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_ADD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Add(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_SUB:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Sub(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_SUB:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Sub(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MUL:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Mul(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MUL:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Mul(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_DIV:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Div(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_DIV:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Div(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MOD:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).Mod(
				scopeGettersDeref[scope2](v, index2),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MOD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Mod(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LT:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				LessThan(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LT:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].LessThan(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LTE:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				LessThanEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LTE:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].LessThanEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GT:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				GreaterThan(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GT:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].GreaterThan(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GTE:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				GreaterThanEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GTE:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].GreaterThanEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_EQ:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Equal(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_EQ:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].Equal(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_NEQ:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				NotEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_NEQ:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			v.stack[v.sp].NotEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		default:
			panic(fmt.Sprintf("unimplemented %d", v.curFrame.instructions[v.curFrame.ip]))
		}
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
	if v.fp >= 0 {
		v.curFrame = v.callStack[v.fp]
	} else {
		v.curFrame = nil
	}
}
