package vm

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/joetifa2003/weaver/opcode"
)

type ValueType uint8

const (
	ValueTypeNil ValueType = iota
	ValueTypeNumber
	ValueTypeString
	ValueTypeObject
	ValueTypeBool
	ValueTypeFunction
	ValueTypeArray
	ValueTypeNativeFunction
	ValueTypeModule
	ValueTypeNativeObject
	ValueTypeRef
	ValueTypeError
	ValueTypeTask
	ValueTypeLock
	ValueTypeChannel
	ValueTypeTime
	ValueTypeIterator
)

func (t ValueType) Is(other ...ValueType) bool {
	if len(other) == 0 {
		return true
	}

	for _, o := range other {
		if t == o {
			return true
		}
	}

	return false
}

func (t ValueType) String() string {
	switch t {
	// TODO: Value zero value should be nil
	case ValueTypeNumber:
		return "number"
	case ValueTypeString:
		return "string"
	case ValueTypeObject:
		return "object"
	case ValueTypeBool:
		return "bool"
	case ValueTypeFunction:
		return "function"
	case ValueTypeNil:
		return "nil"
	case ValueTypeArray:
		return "array"
	case ValueTypeNativeFunction:
		return "native function"
	case ValueTypeNativeObject:
		return "native object"
	case ValueTypeError:
		return "error"
	case ValueTypeLock:
		return "lock"
	case ValueTypeTask:
		return "task"
	case ValueTypeChannel:
		return "channel"
	case ValueTypeTime:
		return "time"
	case ValueTypeIterator:
		return "iterator"
	default:
		panic(fmt.Sprintf("unimplemented %d", t))
	}
}

// Value poor mans union/enum.
type Value struct {
	VType        ValueType
	primitive    [8]byte
	nonPrimitive unsafe.Pointer
}

func interpret[T any](b *[8]byte) *T {
	return (*T)(unsafe.Pointer(b))
}

func (v *Value) SetRef(r Value) {
	v.VType = ValueTypeRef
	v.nonPrimitive = unsafe.Pointer(&r)
}

func (v *Value) deref() *Value {
	switch v.VType {
	case ValueTypeRef:
		return (*Value)(v.nonPrimitive)
	default:
		return v
	}
}

func (v *Value) Set(other Value) {
	v.VType = other.VType
	v.nonPrimitive = other.nonPrimitive
	v.primitive = other.primitive
}

func (v *Value) SetTime(t time.Time) {
	v.VType = ValueTypeTime
	v.nonPrimitive = unsafe.Pointer(&t)
}

func (v *Value) GetTime() time.Time {
	return *(*time.Time)(v.nonPrimitive)
}

func (v *Value) SetIter(iter iter.Seq[Value]) {
	v.VType = ValueTypeIterator
	v.nonPrimitive = unsafe.Pointer(&iter)
}

func (v *Value) GetIter() iter.Seq[Value] {
	return *(*iter.Seq[Value])(v.nonPrimitive)
}

type Task struct {
	C     chan Value
	Value Value
}

func (v *Value) SetTask(c chan Value) {
	v.VType = ValueTypeTask
	task := &Task{
		C: c,
	}
	v.nonPrimitive = unsafe.Pointer(task)
}

func (v *Value) GetTask() *Task {
	return (*Task)(v.nonPrimitive)
}

type Lock struct {
	*sync.Mutex
	lock   NativeFunction
	unlock NativeFunction
}

func (v *Value) SetLock(l *sync.Mutex) {
	v.VType = ValueTypeLock
	v.nonPrimitive = unsafe.Pointer(&Lock{
		Mutex: l,
		lock: func(v *VM, args NativeFunctionArgs) Value {
			if len(args) > 0 {
				fnArg, ok := args.Get(0, ValueTypeFunction)
				if !ok {
					return fnArg
				}

				l.Lock()
				defer l.Unlock()
				v.RunFunction(fnArg)

				return Value{}
			}

			l.Lock()
			return Value{}
		},
		unlock: func(v *VM, args NativeFunctionArgs) Value {
			l.Unlock()
			return Value{}
		},
	})
}

func (v *Value) GetLock() *Lock {
	return (*Lock)(v.nonPrimitive)
}

func (v *Value) SetNumber(f float64) {
	v.VType = ValueTypeNumber
	*interpret[float64](&v.primitive) = f
}

func (v *Value) SetChannel(c chan Value) {
	v.VType = ValueTypeChannel
	v.nonPrimitive = unsafe.Pointer(&c)
}

func (v *Value) GetChannel() chan Value {
	return *(*chan Value)(v.nonPrimitive)
}

func (v *Value) GetNumber() float64 {
	switch v.VType {
	case ValueTypeNumber:
		return *interpret[float64](&v.primitive)
	default:
		return 0
	}
}

func (v *Value) SetObject(o map[string]Value) {
	v.VType = ValueTypeObject
	v.nonPrimitive = unsafe.Pointer(&o)
}

func (v *Value) GetObject() map[string]Value {
	return *(*map[string]Value)(v.nonPrimitive)
}

func (v *Value) SetModule(m map[string]Value) {
	v.VType = ValueTypeModule
	v.nonPrimitive = unsafe.Pointer(&m)
}

func (v *Value) GetModule() map[string]Value {
	return *(*map[string]Value)(v.nonPrimitive)
}

func (v *Value) GetBool() bool {
	return *interpret[bool](&v.primitive)
}

func (v *Value) GetArray() *[]Value {
	return (*[]Value)(v.nonPrimitive)
}

func (v *Value) SetArray(a []Value) {
	v.VType = ValueTypeArray
	v.nonPrimitive = unsafe.Pointer(&a)
}

func (v *Value) SetBool(b bool) {
	v.VType = ValueTypeBool
	*interpret[bool](&v.primitive) = b
}

func (v *Value) SetNil() {
	v.VType = ValueTypeNil
	v.nonPrimitive = nil
}

type Error struct {
	msg  string
	data Value
}

func (e *Error) Error() string {
	return e.msg
}

func (v *Value) SetError(msg string, data Value) {
	e := Error{msg: msg, data: data}
	v.VType = ValueTypeError
	v.nonPrimitive = unsafe.Pointer(&e)
}

func (v *Value) GetError() *Error {
	return (*Error)(v.nonPrimitive)
}

type FunctionValue struct {
	NumVars      int
	Instructions []opcode.OpCode
	FreeVars     []Value
	Constants    []Value
	Path         string
}

func (v *Value) SetFunction(f FunctionValue) {
	v.VType = ValueTypeFunction
	v.nonPrimitive = unsafe.Pointer(&f)
}

func (v *Value) GetFunction() *FunctionValue {
	return (*FunctionValue)(v.nonPrimitive)
}

func (v *Value) SetString(s string) {
	v.VType = ValueTypeString
	v.nonPrimitive = unsafe.Pointer(&s)
}

func (v *Value) GetString() string {
	return *(*string)(v.nonPrimitive)
}

func (v *Value) SetNativeObject(obj interface{}) {
	v.VType = ValueTypeNativeObject
	v.nonPrimitive = unsafe.Pointer(&obj)
}

func (v *Value) GetNativeObject() interface{} {
	return *(*interface{})(v.nonPrimitive)
}

type NativeFunctionArgs []Value

func (a NativeFunctionArgs) Get(i int, types ...ValueType) (Value, bool) {
	if i >= len(a) {
		return NewError("invalid number of arguments", Value{}), false
	}

	return CheckValueType(a[i], types...)
}

func (a NativeFunctionArgs) Len() int {
	return len(a)
}

func CheckValueType(val Value, types ...ValueType) (Value, bool) {
	if val.VType.Is(types...) {
		return val, true
	}

	return NewError(fmt.Sprintf("invalid argument type, expected %v", types), Value{}), false
}

type NativeFunction func(v *VM, args NativeFunctionArgs) Value

func (v *Value) GetNativeFunction() NativeFunction {
	return *(*NativeFunction)(v.nonPrimitive)
}

func (v *Value) SetNativeFunction(f NativeFunction) {
	v.VType = ValueTypeNativeFunction
	v.nonPrimitive = unsafe.Pointer(&f)
}

func NewString(s string) Value {
	val := Value{}
	val.SetString(s)
	return val
}

func NewArray(a []Value) Value {
	val := Value{}
	val.SetArray(a)
	return val
}

func NewNativeFunction(f NativeFunction) Value {
	val := Value{}
	val.SetNativeFunction(f)
	return val
}

func NewNativeObject(o interface{}) Value {
	val := Value{}
	val.SetNativeObject(o)
	return val
}

func NewNumber(f float64) Value {
	val := Value{}
	val.SetNumber(f)
	return val
}

func NewObject(m map[string]Value) Value {
	val := Value{}
	val.SetObject(m)
	return val
}

func NewBoolean(b bool) Value {
	val := Value{}
	val.SetBool(b)
	return val
}

func NewBool(b bool) Value {
	val := Value{}
	val.SetBool(b)
	return val
}

func NewError(msg string, data Value) Value {
	val := Value{}
	val.SetError(msg, data)
	return val
}

func NewTime(t time.Time) Value {
	val := Value{}
	val.SetTime(t)
	return val
}

func NewIter(iter iter.Seq[Value]) Value {
	val := Value{}
	val.SetIter(iter)
	return val
}

func NewErrFromErr(err error) Value {
	return NewError(err.Error(), Value{})
}

func (v *Value) String() string { return v.string(0) }

func (v *Value) string(i int) string {
	switch v.VType {
	case ValueTypeModule:
		return "module"

	case ValueTypeString:
		str := v.GetString()
		return str

	case ValueTypeNumber:
		num := v.GetNumber()
		return strconv.FormatFloat(num, 'f', -1, 64)

	case ValueTypeNil:
		return "nil"

	case ValueTypeObject:
		builder := strings.Builder{}
		builder.WriteString("{\n")
		for k, v := range v.GetObject() {
			builder.WriteString(fmt.Sprintf("%s%s: %s \n", strings.Repeat("  ", i+1), k, v.string(i+1)))
		}
		builder.WriteString(fmt.Sprintf("%s}", strings.Repeat("  ", i)))
		return builder.String()

	case ValueTypeBool:
		return strconv.FormatBool(v.GetBool())

	case ValueTypeFunction:
		return "function"

	case ValueTypeArray:
		builder := strings.Builder{}
		builder.WriteString("[\n")
		for _, v := range *v.GetArray() {
			builder.WriteString(fmt.Sprintf("%s%s \n", strings.Repeat("  ", i+1), v.string(i+1)))
		}
		builder.WriteString(fmt.Sprintf("%s]", strings.Repeat("  ", i)))
		return builder.String()

	case ValueTypeNativeFunction:
		return "native function"

	case ValueTypeError:
		err := v.GetError()
		msg := err.msg
		if err.data.VType == ValueTypeNil {
			return fmt.Sprintf("error(%s)", msg)
		}

		return fmt.Sprintf("error(%s, %s)", msg, err.data.String())

	case ValueTypeChannel:
		return "channel"

	case ValueTypeTask:
		return "task"

	case ValueTypeTime:
		return v.GetTime().String()

	case ValueTypeIterator:
		return "iterator"

	default:
		panic(fmt.Sprintf("Value.String(): unimplemented %T", v.VType))
	}
}

func (v *Value) IsTruthy() bool {
	switch v.VType {
	case ValueTypeBool:
		return v.GetBool()
	case ValueTypeNil:
		return false
	case ValueTypeError:
		return false
	default:
		return true
	}
}

func (v *Value) Add(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetNumber(v.GetNumber() + other.GetNumber())
		return
	}

	if v.VType == ValueTypeString && other.VType == ValueTypeString {
		res.SetString(v.GetString() + other.GetString())
		return
	}

	panic(fmt.Sprintf("illegal operation %s + %s", v.VType, other.VType))
}

func (v *Value) Sub(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetNumber(v.GetNumber() - other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s - %s", v.VType, other.VType))
}

func (v *Value) Mul(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetNumber(v.GetNumber() * other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s * %s", v.VType, other.VType))
}

func (v *Value) Div(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetNumber(v.GetNumber() / other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s / %s", v.VType, other.VType))
}

func (v *Value) LessThan(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetBool(v.GetNumber() < other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s < %s", v.VType, other.VType))
}

func (v *Value) LessThanEqual(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetBool(v.GetNumber() <= other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s <= %s", v.VType, other.VType))
}

func (v *Value) GreaterThan(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetBool(v.GetNumber() > other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s > %s", v.VType, other.VType))
}

func (v *Value) GreaterThanEqual(other *Value, res *Value) {
	if v.VType == ValueTypeNumber && other.VType == ValueTypeNumber {
		res.SetBool(v.GetNumber() >= other.GetNumber())
		return
	}

	panic(fmt.Sprintf("illegal operation %s >= %s", v.VType, other.VType))
}

func (v *Value) Equal(other *Value, res *Value) {
	if v.VType != other.VType {
		res.SetBool(false)
		return
	}

	switch v.VType {
	case ValueTypeNil:
		res.SetBool(true)
	case ValueTypeNumber:
		res.SetBool(v.GetNumber() == other.GetNumber())
	case ValueTypeString:
		res.SetBool(v.GetString() == other.GetString())
	case ValueTypeArray:
		res.SetBool(v.GetArray() == other.GetArray())
	case ValueTypeBool:
		res.SetBool(v.GetBool() == other.GetBool())
	case ValueTypeFunction:
		res.SetBool(v.GetFunction() == other.GetFunction())
	case ValueTypeTime:
		res.SetBool(v.GetTime().Equal(other.GetTime()))
	default:
		res.SetBool(false)
	}
}

func (v *Value) NotEqual(other *Value, res *Value) {
	if v.VType != other.VType {
		res.SetBool(true)
		return
	}

	switch v.VType {
	case ValueTypeNil:
		res.SetBool(false)
	case ValueTypeNumber:
		res.SetBool(v.GetNumber() != other.GetNumber())
	case ValueTypeString:
		res.SetBool(v.GetString() != other.GetString())
	case ValueTypeArray:
		res.SetBool(v.GetArray() != other.GetArray())
	case ValueTypeBool:
		res.SetBool(v.GetBool() != other.GetBool())
	case ValueTypeFunction:
		res.SetBool(v.GetFunction() != other.GetFunction())
	case ValueTypeTime:
		res.SetBool(!v.GetTime().Equal(other.GetTime()))
	default:
		res.SetBool(true)
	}
}

func (v *Value) Negate(res *Value) {
	switch v.VType {
	case ValueTypeNumber:
		res.SetNumber(-v.GetNumber())
	default:
		panic(fmt.Sprintf("illegal operation -%s", v))
	}
}

func (v *Value) Mod(other *Value, res *Value) {
	if (v.VType != ValueTypeNumber) || (other.VType != ValueTypeNumber) {
		panic(fmt.Sprintf("illegal operation %s %% %s", v, other))
	}

	res.SetNumber(float64(int(v.GetNumber()) % int(other.GetNumber())))
}

var (
	ErrInvalidArrayIndexType  = errors.New("invalid array index type")
	ErrInvalidObjectIndexType = errors.New("invalid object index type")
	ErrInvalidErrorIndexType  = errors.New("invalid error index type")
)

func (v *Value) Index(idx *Value, res *Value) {
	switch v.VType {
	case ValueTypeArray:
		switch idx.VType {
		case ValueTypeNumber:
			res.Set((*v.GetArray())[int(idx.GetNumber())])
			return
		}
	case ValueTypeObject:
		switch idx.VType {
		case ValueTypeString:
			idx := idx.GetString()
			obj := v.GetObject()
			res.Set(obj[idx])
			return
		}

	case ValueTypeError:
		switch idx.VType {
		case ValueTypeString:
			err := v.GetError()
			switch idx.GetString() {
			case "msg":
				res.SetString(err.msg)
				return
			case "data":
				res.Set(err.data)
				return
			}
		}

	case ValueTypeLock:
		lock := v.GetLock()

		switch idx.VType {
		case ValueTypeString:
			switch idx.GetString() {
			case "lock":
				res.SetNativeFunction(lock.lock)
				return
			case "unlock":
				res.SetNativeFunction(lock.unlock)
				return
			}
		}
	}

	res.SetNil()
}

func (v *Value) SetIndex(idx *Value, val Value) {
	switch v.VType {
	case ValueTypeArray:
		switch idx.VType {
		case ValueTypeNumber:
			(*v.GetArray())[int(idx.GetNumber())] = val
		default:
			panic(ErrInvalidArrayIndexType)
		}
	case ValueTypeObject:
		switch idx.VType {
		case ValueTypeString:
			v.GetObject()[idx.GetString()] = val
		default:
			panic(ErrInvalidObjectIndexType)
		}
	case ValueTypeError:
		switch idx.VType {
		case ValueTypeString:
			err := v.GetError()
			key := idx.GetString()

			switch key {
			case "msg":
				err.msg = val.GetString()
			case "data":
				err.data = val
			default:
				panic(ErrInvalidObjectIndexType)
			}
		default:
			panic(ErrInvalidErrorIndexType)
		}
	}
}

func (v *Value) IsError() bool {
	return v.VType.Is(ValueTypeError)
}
