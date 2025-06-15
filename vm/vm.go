package vm

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/joetifa2003/weaver/opcode"
)

const MaxStack = 1024

const MaxCallStack = 1024

type Frame struct {
	Instructions []opcode.OpCode
	FreeVars     []Value
	HaltAfter    bool
	NumVars      int
	Path         string
	Constants    []Value

	ip          int
	stackOffset int
	returnAddr  int
	hasTry      bool
}

type VM struct {
	Executor *Executor
	Ctx      context.Context

	stack     [MaxStack]Value
	callStack [MaxCallStack]Frame
	curFrame  *Frame
	reg       *Registry
	ctxCancel context.CancelFunc
	running   *atomic.Bool
	sp        int
	fp        int
}

func New(executor *Executor) *VM {
	running := atomic.Bool{}
	running.Store(true)
	ctx, cancel := context.WithCancel(context.Background())
	vm := &VM{
		Executor:  executor,
		sp:        -1,
		fp:        -1,
		running:   &running,
		Ctx:       ctx,
		ctxCancel: cancel,
	}

	return vm
}

func (v *VM) Resurrect() {
	v.running.Store(true)
	v.sp = -1
	v.fp = -1
	v.curFrame = nil
	v.Ctx, v.ctxCancel = context.WithCancel(context.Background())
}

func (v *VM) Stop() {
	v.ctxCancel()
	v.running.Store(false)
}

var scopeGettersDeref = [4]func(v *VM, idx int) *Value{
	opcode.ScopeTypeConst: func(v *VM, idx int) *Value {
		return v.curFrame.Constants[idx].deref()
	},
	opcode.ScopeTypeGlobal: func(v *VM, idx int) *Value {
		return v.stack[idx].deref()
	},
	opcode.ScopeTypeLocal: func(v *VM, idx int) *Value {
		return v.stack[v.curFrame.stackOffset+idx].deref()
	},
	opcode.ScopeTypeFree: func(v *VM, idx int) *Value {
		return v.curFrame.FreeVars[idx].deref()
	},
}

var scopeGetters = [4]func(v *VM, idx int) *Value{
	opcode.ScopeTypeConst: func(v *VM, idx int) *Value {
		return &v.curFrame.Constants[idx]
	},
	opcode.ScopeTypeGlobal: func(v *VM, idx int) *Value {
		return &v.stack[idx]
	},
	opcode.ScopeTypeLocal: func(v *VM, idx int) *Value {
		return &v.stack[v.curFrame.stackOffset+idx]
	},
	opcode.ScopeTypeFree: func(v *VM, idx int) *Value {
		return &v.curFrame.FreeVars[idx]
	},
}

func (v *VM) CurrentFrame() *Frame {
	return v.curFrame
}

func (v *VM) GetStackVal(idx int) Value {
	return v.stack[idx]
}

func (v *VM) Run(frame Frame, args int) bool {
	v.pushFrame(frame, args)

	for {
		if !v.running.Load() {
			return true
		}

		switch v.curFrame.Instructions[v.curFrame.ip] {
		case opcode.OP_TRY:
			v.curFrame.hasTry = !v.curFrame.hasTry
			v.curFrame.ip += 1

		case opcode.OP_UPGRADE_REF:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			idx := int(v.curFrame.Instructions[v.curFrame.ip+2])

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
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])
			f := v.stack[v.sp].GetFunction()

			emptyFunc := scopeGettersDeref[scope](v, index).GetFunction()
			emptyFunc.FreeVars = f.FreeVars
			emptyFunc.Instructions = f.Instructions
			emptyFunc.NumVars = f.NumVars
			emptyFunc.Constants = f.Constants
			v.curFrame.ip += 3

		case opcode.OP_FUNC:
			constantIndex := int(v.curFrame.Instructions[v.curFrame.ip+1])
			freeVarsCount := int(v.curFrame.Instructions[v.curFrame.ip+2])
			var freeVars []Value

			for range freeVarsCount {
				freeVars = append(freeVars, v.stack[v.sp])
				v.stack[v.sp].SetNil()
				v.sp--
			}

			fn := *v.curFrame.Constants[constantIndex].GetFunction()
			fn.FreeVars = freeVars
			v.sp++
			v.stack[v.sp].SetFunction(fn)

			v.curFrame.ip += 3

		case opcode.OP_INC:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() + 1)
			v.sp++
			v.stack[v.sp] = *v1
			v.curFrame.ip += 3

		case opcode.OP_INC_POP:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() + 1)
			v.curFrame.ip += 3

		case opcode.OP_DEC:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() - 1)
			v.sp++
			v.stack[v.sp] = *v1
			v.curFrame.ip += 3

		case opcode.OP_DEC_POP:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])
			v1 := scopeGettersDeref[scope](v, index)
			v1.SetNumber(v1.GetNumber() - 1)
			v.curFrame.ip += 3

		case opcode.OP_LOAD:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.sp++
			v.stack[v.sp] = *scopeGettersDeref[scope](v, index)
			v.curFrame.ip += 3

		case opcode.OP_STORE_FREE:
			index := int(v.curFrame.Instructions[v.curFrame.ip+1])
			val := v.stack[v.sp]
			v.curFrame.FreeVars[index].deref().Set(val)
			v.curFrame.ip += 2

		case opcode.OP_LET:
			index := v.curFrame.stackOffset + int(v.curFrame.Instructions[v.curFrame.ip+1])
			v.stack[index] = v.stack[v.sp]
			v.sp--
			v.curFrame.ip += 2

		case opcode.OP_STORE:
			index := v.curFrame.stackOffset + int(v.curFrame.Instructions[v.curFrame.ip+1])
			isRef := int(v.curFrame.Instructions[v.curFrame.ip+2])

			if isRef == 1 {
				v.stack[index].deref().Set(v.stack[v.sp])
			} else {
				v.stack[index] = v.stack[v.sp]
			}
			v.curFrame.ip += 3

		case opcode.OP_ECHO:
			value := v.stack[v.sp]
			v.sp--
			fmt.Println(value.String())
			v.curFrame.ip += 1

		case opcode.OP_JUMP:
			newIp := int(v.curFrame.Instructions[v.curFrame.ip+1])
			v.curFrame.ip = newIp

		case opcode.OP_PJUMP_F:
			newIp := int(v.curFrame.Instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip += 2
			} else {
				v.curFrame.ip = newIp
			}

			v.sp--

		case opcode.OP_PJUMP_T:
			newIp := int(v.curFrame.Instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip = newIp
			} else {
				v.curFrame.ip += 2
			}

			v.sp--

		case opcode.OP_JUMP_F:
			newIp := int(v.curFrame.Instructions[v.curFrame.ip+1])

			if v.stack[v.sp].IsTruthy() {
				v.curFrame.ip += 2
			} else {
				v.curFrame.ip = newIp
			}

		case opcode.OP_JUMP_T:
			newIp := int(v.curFrame.Instructions[v.curFrame.ip+1])

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
			numArgs := int(v.curFrame.Instructions[v.curFrame.ip+1])
			calleeIdx := v.sp - numArgs
			argsBegin := calleeIdx + 1
			callee := v.stack[calleeIdx]

			switch callee.VType {
			case ValueTypeFunction:
				fn := callee.GetFunction()
				frame := Frame{
					Instructions: fn.Instructions,
					NumVars:      fn.NumVars,
					FreeVars:     fn.FreeVars,
					Constants:    fn.Constants,
					Path:         fn.Path,
					ip:           0,
					stackOffset:  argsBegin,
					returnAddr:   calleeIdx,
				}
				v.curFrame.ip += 2
				v.pushFrame(frame, numArgs)
			case ValueTypeNativeFunction:
				fn := callee.GetNativeFunction()
				args := v.stack[argsBegin : argsBegin+numArgs]
				r, ok := fn(v, args)
				v.sp = calleeIdx
				v.stack[v.sp] = r
				v.curFrame.ip += 2
				if !ok && r.IsError() {
					if !v.raise(r) {
						return false
					}
				}
			default:
				panic(fmt.Sprintf("illegal callee type %s", callee.VType))
			}

		case opcode.OP_RAISE:
			val := v.stack[v.sp]
			if !v.raise(val) {
				return false
			}

		case opcode.OP_RET:
			val := v.stack[v.sp]
			v.sp = v.curFrame.returnAddr
			v.stack[v.sp] = val
			haltAfter := v.curFrame.HaltAfter
			v.popFrame()
			if haltAfter {
				return true
			}

		case opcode.OP_HALT:
			return true

		case opcode.OP_LOAD_LOAD_ADD:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Add(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_ADD:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Add(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_SUB:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Sub(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_SUB:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Sub(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MUL:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Mul(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MUL:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Mul(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_DIV:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Div(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_DIV:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Div(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_MOD:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).Mod(
				scopeGettersDeref[scope2](v, index2),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_MOD:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Mod(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LT:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				LessThan(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LT:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].LessThan(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_LTE:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				LessThanEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_LTE:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].LessThanEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GT:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				GreaterThan(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GT:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].GreaterThan(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_GTE:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				GreaterThanEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_GTE:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].GreaterThanEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_EQ:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				Equal(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_EQ:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].Equal(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		case opcode.OP_LOAD_LOAD_NEQ:
			scope1 := v.curFrame.Instructions[v.curFrame.ip+1]
			index1 := int(v.curFrame.Instructions[v.curFrame.ip+2])
			scope2 := v.curFrame.Instructions[v.curFrame.ip+3]
			index2 := int(v.curFrame.Instructions[v.curFrame.ip+4])

			v.sp++

			scopeGettersDeref[scope1](v, index1).
				NotEqual(
					scopeGettersDeref[scope2](v, index2),
					&v.stack[v.sp],
				)

			v.curFrame.ip += 5

		case opcode.OP_LOAD_NEQ:
			scope := v.curFrame.Instructions[v.curFrame.ip+1]
			index := int(v.curFrame.Instructions[v.curFrame.ip+2])

			v.stack[v.sp].NotEqual(
				scopeGettersDeref[scope](v, index),
				&v.stack[v.sp],
			)

			v.curFrame.ip += 3

		default:
			panic(fmt.Sprintf("unimplemented %d", v.curFrame.Instructions[v.curFrame.ip]))
		}
	}
}

func (v *VM) raise(val Value) bool {
	prevFrame := v.curFrame
	v.popFrame() // pop the current frame

	for v.fp >= 0 {
		if v.callStack[v.fp].hasTry {
			v.sp = prevFrame.returnAddr
			v.stack[v.sp] = val
			return true
		}

		if v.curFrame.HaltAfter {
			v.sp = v.curFrame.returnAddr
			v.stack[v.sp] = val
			return false
		}

		prevFrame = v.curFrame
		v.popFrame()
	}

	return false
}

func (v *VM) pushFrame(f Frame, args int) {
	v.fp++

	// local variables initialization
	for range f.NumVars - args {
		v.sp++
	}

	// TODO: handle case when args > NumVars

	v.callStack[v.fp] = f
	v.curFrame = &v.callStack[v.fp]
}

func (v *VM) popFrame() {
	v.fp--
	if v.fp >= 0 {
		v.curFrame = &v.callStack[v.fp]
	} else {
		v.curFrame = nil
	}
}

func (v *VM) RunFunction(f Value, args ...Value) (Value, bool) {
	fn := f.GetFunction()
	v.sp++
	retAddr := v.sp
	for _, arg := range args {
		v.sp++
		v.stack[v.sp] = arg
	}

	ok := v.Run(Frame{
		Instructions: fn.Instructions,
		NumVars:      fn.NumVars,
		FreeVars:     fn.FreeVars,
		Constants:    fn.Constants,
		Path:         fn.Path,
		HaltAfter:    true,
		ip:           0,
		stackOffset:  retAddr + 1,
		returnAddr:   retAddr,
	}, len(args))

	retVal := v.stack[retAddr]
	v.sp--

	if !ok {
		return retVal, false
	}

	return retVal, true
}
