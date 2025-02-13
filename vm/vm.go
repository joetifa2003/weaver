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
		case opcode.OP_ADD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--
			left.Add(right, &v.stack[v.sp])

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

			left.Mul(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_MOD:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Mod(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_EQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.Equal(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_NEQ:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.NotEqual(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_LT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.LessThan(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_LTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.LessThanEqual(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_GT:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.GreaterThan(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_GTE:
			right := v.stack[v.sp]
			left := v.stack[v.sp-1]
			v.sp--

			left.GreaterThanEqual(right, &v.stack[v.sp])

			v.curFrame.ip++

		case opcode.OP_NOT:
			v.stack[v.sp].SetBool(!v.stack[v.sp].IsTruthy())
			v.curFrame.ip++

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

		case opcode.OP_INC_LOCAL:
			idx := v.curFrame.instructions[v.curFrame.ip+1]
			i := v.stack[idx].GetInt()
			v.stack[idx].SetInt(i + 1)
			v.sp++
			v.stack[v.sp].SetInt(i + 1)
			v.curFrame.ip += 2

		case opcode.OP_LOAD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.sp++
				v.stack[v.sp] = v.stack[v.curFrame.stackOffset+index]
				v.curFrame.ip += 3

			case opcode.ScopeTypeFree:
				v.sp++
				v.stack[v.sp] = v.curFrame.freeVars[index]
				v.curFrame.ip += 3

			case opcode.ScopeTypeGlobal:
				v.sp++
				v.stack[v.sp] = v.stack[index]
				v.curFrame.ip += 3

			case opcode.ScopeTypeConst:
				v.sp++
				v.stack[v.sp] = v.constants[index]
				v.curFrame.ip += 3
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

		case opcode.OP_STORE_FREE:
			index := int(v.curFrame.instructions[v.curFrame.ip+1])
			val := v.stack[v.sp]
			v.curFrame.freeVars[index] = val
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

		case opcode.OP_HALT:
			return

		case opcode.OP_LOAD_LOAD_ADD:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Add(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Add(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Add(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Add(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_ADD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Add(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Add(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Add(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Add(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_SUB:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Sub(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Sub(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Sub(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Sub(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_SUB:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Sub(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Sub(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Sub(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Sub(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MUL:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Mul(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Mul(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Mul(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Mul(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MUL:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Mul(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Mul(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Mul(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Mul(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_DIV:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Div(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Div(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Div(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Div(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_DIV:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Div(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Div(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Div(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Div(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MOD:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Mod(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Mod(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Mod(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Mod(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MOD:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Mod(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Mod(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Mod(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Mod(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LT:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.LessThan(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.LessThan(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.LessThan(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.LessThan(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LT:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].LessThan(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].LessThan(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].LessThan(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].LessThan(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LTE:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.LessThanEqual(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.LessThanEqual(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.LessThanEqual(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.LessThanEqual(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LTE:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].LessThanEqual(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].LessThanEqual(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].LessThanEqual(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].LessThanEqual(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GT:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.GreaterThan(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.GreaterThan(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.GreaterThan(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.GreaterThan(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GT:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].GreaterThan(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].GreaterThan(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].GreaterThan(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].GreaterThan(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GTE:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.GreaterThanEqual(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.GreaterThanEqual(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.GreaterThanEqual(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.GreaterThanEqual(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GTE:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].GreaterThanEqual(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].GreaterThanEqual(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].GreaterThanEqual(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].GreaterThanEqual(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_EQ:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.Equal(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.Equal(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.Equal(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.Equal(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_EQ:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].Equal(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].Equal(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].Equal(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].Equal(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_NEQ:
			scope1 := v.curFrame.instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.instructions[v.curFrame.ip+4])

			v.sp++

			var v1 Value

			switch scope1 {
			case opcode.ScopeTypeLocal:
				v1 = v.stack[v.curFrame.stackOffset+index1]
			case opcode.ScopeTypeFree:
				v1 = v.curFrame.freeVars[index1]
			case opcode.ScopeTypeGlobal:
				v1 = v.stack[index1]
			case opcode.ScopeTypeConst:
				v1 = v.constants[index1]
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope1))
			}

			switch scope2 {
			case opcode.ScopeTypeLocal:
				v1.NotEqual(v.stack[v.curFrame.stackOffset+index2], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v1.NotEqual(v.curFrame.freeVars[index2], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v1.NotEqual(v.stack[index2], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v1.NotEqual(v.constants[index2], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope2))
			}

			v.curFrame.ip += 5

		case opcode.OP_LOAD_NEQ:
			scope := v.curFrame.instructions[v.curFrame.ip+1]
			index := int(v.curFrame.instructions[v.curFrame.ip+2])

			switch scope {
			case opcode.ScopeTypeLocal:
				v.stack[v.sp].NotEqual(v.stack[v.curFrame.stackOffset+index], &v.stack[v.sp])
			case opcode.ScopeTypeFree:
				v.stack[v.sp].NotEqual(v.curFrame.freeVars[index], &v.stack[v.sp])
			case opcode.ScopeTypeGlobal:
				v.stack[v.sp].NotEqual(v.stack[index], &v.stack[v.sp])
			case opcode.ScopeTypeConst:
				v.stack[v.sp].NotEqual(v.constants[index], &v.stack[v.sp])
			default:
				panic(fmt.Sprintf("unimplemented scope %d", scope))
			}

			v.curFrame.ip += 3

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
