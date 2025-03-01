package vm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func (v *Value) SetNumber(f float64) {
	v.VType = ValueTypeNumber
	*interpret[float64](&v.primitive) = f
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
}

type Error struct {
	Data Value
}

func (v *Value) SetError(data Value) {
	e := Error{Data: data}
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

func (a NativeFunctionArgs) Get(i int, types ...ValueType) (Value, error) {
	if i >= len(a) {
		return Value{}, ErrInvalidNumberOfArguments
	}

	v := a[i]
	if !v.VType.Is(types...) {
		return Value{}, ErrInvalidArgType
	}

	return a[i], nil
}

type NativeFunction func(v *VM, args NativeFunctionArgs) (Value, error)

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

func NewBool(b bool) Value {
	val := Value{}
	val.SetBool(b)
	return val
}

func NewError(data Value) Value {
	val := Value{}
	val.SetError(data)
	return val
}

func (v *Value) String() string {
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
		builder.WriteString("{")
		for k, v := range v.GetObject() {
			builder.WriteString(fmt.Sprintf("%s: %s, ", k, v.String()))
		}
		builder.WriteString("}")
		return builder.String()

	case ValueTypeBool:
		return strconv.FormatBool(v.GetBool())

	case ValueTypeFunction:
		return "function"

	case ValueTypeArray:
		arr := *v.GetArray()
		return fmt.Sprint(arr)

	case ValueTypeNativeFunction:
		return "native function"

	case ValueTypeError:
		data := v.GetError().Data
		return fmt.Sprintf("error(%s)", data.String())

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

	panic(fmt.Sprintf("illegal operation %s - %s", v.VType, other.VType))
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
		default:
			panic(ErrInvalidArrayIndexType)
		}
	case ValueTypeObject:
		switch idx.VType {
		case ValueTypeString:
			res.Set(v.GetObject()[idx.GetString()])
		default:
			panic(ErrInvalidObjectIndexType)
		}
	case ValueTypeError:
		switch idx.VType {
		case ValueTypeString:
			if idx.GetString() == "data" {
				res.Set(v.GetError().Data)
			} else {
				res.Set(Value{})
			}
		default:
			panic(ErrInvalidErrorIndexType)
		}
	}
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
			if idx.GetString() == "data" {
				v.GetError().Data = val
			} else {
				panic(ErrInvalidErrorIndexType)
			}
		default:
			panic(ErrInvalidErrorIndexType)
		}
	}
}
