package vm

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/joetifa2003/weaver/opcode"
)

type ValueType int

const (
	ValueTypeNil ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeString
	ValueTypeObject
	ValueTypeBool
	ValueTypeFunction
	ValueTypeArray
	ValueTypeNativeFunction
)

func (t ValueType) String() string {
	switch t {
	case ValueTypeInt:
		return "int"
	case ValueTypeFloat:
		return "float"
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
	default:
		panic(fmt.Sprintf("unimplemented %T", t))
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

func (v *Value) GetInt() int {
	if v.VType != ValueTypeInt {
		panic("Value.GetInt(): not an int")
	}

	return *interpret[int](&v.primitive)
}

func (v *Value) GetFloat() float64 {
	if v.VType != ValueTypeFloat {
		panic("Value.GetFloat(): not a float")
	}

	return *interpret[float64](&v.primitive)
}

func (v *Value) SetObject(o map[string]Value) {
	v.VType = ValueTypeObject
	v.nonPrimitive = unsafe.Pointer(&o)
}

func (v *Value) GetObject() map[string]Value {
	if v.VType != ValueTypeObject {
		panic("Value.GetObject(): not an object")
	}

	return *(*map[string]Value)(v.nonPrimitive)
}

func (v *Value) GetBool() bool {
	if v.VType != ValueTypeBool {
		panic("Value.GetBool(): not a bool")
	}

	return *interpret[bool](&v.primitive)
}

func (v *Value) GetArray() *[]Value {
	if v.VType != ValueTypeArray {
		panic("Value.GetArray(): not an array")
	}

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

func (v *Value) SetInt(i int) {
	v.VType = ValueTypeInt
	*interpret[int](&v.primitive) = i
}

func (v *Value) SetFloat(f float64) {
	v.VType = ValueTypeFloat
	*interpret[float64](&v.primitive) = f
}

func (v *Value) SetNil() {
	v.VType = ValueTypeNil
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
	if v.VType != ValueTypeFunction {
		panic("Value.GetFunction(): not a function")
	}

	return (*FunctionValue)(v.nonPrimitive)
}

func (v *Value) SetString(s string) {
	v.VType = ValueTypeString
	v.nonPrimitive = unsafe.Pointer(&s)
}

func (v *Value) GetString() string {
	if v.VType != ValueTypeString {
		panic("Value.GetString(): not a string")
	}

	return *(*string)(v.nonPrimitive)
}

type NativeFunction func(v *VM, args ...Value) Value

func (v *Value) GetNativeFunction() NativeFunction {
	if v.VType != ValueTypeNativeFunction {
		panic("Value.GetNativeFunction(): not a native function")
	}

	return *(*NativeFunction)(v.nonPrimitive)
}

func (v *Value) SetNativeFunction(f NativeFunction) {
	v.VType = ValueTypeNativeFunction
	v.nonPrimitive = unsafe.Pointer(&f)
}

func (v Value) String() string {
	switch v.VType {
	case ValueTypeString:
		str := v.GetString()
		return str

	case ValueTypeInt:
		integer := v.GetInt()
		return strconv.Itoa(integer)

	case ValueTypeFloat:
		float := v.GetFloat()
		return fmt.Sprint(float)

	case ValueTypeNil:
		return "nil"

	case ValueTypeObject:
		return fmt.Sprint(v.GetObject())

	case ValueTypeBool:
		return strconv.FormatBool(v.GetBool())

	case ValueTypeFunction:
		return "function"

	case ValueTypeArray:
		arr := *v.GetArray()
		return fmt.Sprint(arr)

	case ValueTypeNativeFunction:
		return "native function"

	default:
		panic(fmt.Sprintf("Value.String(): unimplemented %T", v.VType))
	}
}

func (v Value) IsTruthy() bool {
	switch v.VType {
	case ValueTypeBool:
		return v.GetBool()
	case ValueTypeInt:
		return v.GetInt() != 0
	case ValueTypeFloat:
		return v.GetFloat() != 0
	case ValueTypeString:
		return len(v.GetString()) != 0
	case ValueTypeObject:
		return len(v.GetObject()) != 0
	case ValueTypeArray:
		return len(*v.GetArray()) != 0
	case ValueTypeNil:
		return false
	default:
		return false
	}
}

func (v Value) Add(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeString:
		switch other.VType {
		case ValueTypeString:
			res.SetString(v.GetString() + other.GetString())
		default:
			panic(fmt.Sprintf("illegal operation %s + %s", v, other))
		}
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetInt(v.GetInt() + other.GetInt())
		case ValueTypeFloat:
			res.SetFloat(float64(v.GetInt()) + other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s + %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetFloat(v.GetFloat() + float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetFloat(v.GetFloat() + other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s + %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s + %s", v, other))
	}
}

func (v Value) Sub(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetInt(v.GetInt() - other.GetInt())
		case ValueTypeFloat:
			res.SetFloat(float64(v.GetInt()) - other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s - %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetFloat(v.GetFloat() - float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetFloat(v.GetFloat() - other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s - %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s - %s", v, other))
	}
}

func (v Value) Mul(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetInt(v.GetInt() * other.GetInt())
		case ValueTypeFloat:
			res.SetFloat(float64(v.GetInt()) * other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s * %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetFloat(v.GetFloat() * float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetFloat(v.GetFloat() * other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s * %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s * %s", v, other))
	}
}

func (v Value) Div(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetFloat(float64(v.GetInt()) / float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetFloat(float64(v.GetInt()) / other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s / %s", v, other))
		}

	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetFloat(v.GetFloat() / float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetFloat(v.GetFloat() / other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s / %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s / %s", v, other))
	}
}

func (v Value) Mod(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetInt(v.GetInt() % other.GetInt())
		case ValueTypeFloat:
			res.SetInt(v.GetInt() % int(other.GetFloat()))
		default:
			panic(fmt.Sprintf("illegal operation %s %% %s", v, other))
		}

	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetInt(int(v.GetFloat()) % other.GetInt())
		case ValueTypeFloat:
			res.SetInt(int(v.GetFloat()) % int(other.GetFloat()))
		default:
			panic(fmt.Sprintf("illegal operation %s %% %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s %% %s", v, other))
	}
}

func (v Value) LessThan(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetInt() < other.GetInt())
		case ValueTypeFloat:
			res.SetBool(float64(v.GetInt()) < other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s < %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetFloat() < float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetBool(v.GetFloat() < other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s < %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s < %s", v, other))
	}
}

func (v Value) LessThanEqual(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetInt() <= other.GetInt())
		case ValueTypeFloat:
			res.SetBool(float64(v.GetInt()) <= other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s <= %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetFloat() <= float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetBool(v.GetFloat() <= other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s <= %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s <= %s", v, other))
	}
}

func (v Value) GreaterThan(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetInt() > other.GetInt())
		case ValueTypeFloat:
			res.SetBool(float64(v.GetInt()) > other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s > %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetFloat() > float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetBool(v.GetFloat() > other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s > %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s > %s", v, other))
	}
}

func (v Value) GreaterThanEqual(other *Value, res *Value) {
	switch v.VType {
	case ValueTypeInt:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetInt() >= other.GetInt())
		case ValueTypeFloat:
			res.SetBool(float64(v.GetInt()) >= other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s >= %s", v, other))
		}
	case ValueTypeFloat:
		switch other.VType {
		case ValueTypeInt:
			res.SetBool(v.GetFloat() >= float64(other.GetInt()))
		case ValueTypeFloat:
			res.SetBool(v.GetFloat() >= other.GetFloat())
		default:
			panic(fmt.Sprintf("illegal operation %s >= %s", v, other))
		}
	default:
		panic(fmt.Sprintf("illegal operation %s >= %s", v, other))
	}
}

// TODO: implement object equality.
func (v Value) Equal(other *Value, res *Value) {
	if v.VType != other.VType {
		res.SetBool(false)
		return
	}

	switch v.VType {
	case ValueTypeNil:
		res.SetBool(true)
	case ValueTypeInt:
		res.SetBool(v.GetInt() == other.GetInt())
	case ValueTypeFloat:
		res.SetBool(v.GetFloat() == other.GetFloat())
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

func (v Value) NotEqual(other *Value, res *Value) {
	if v.VType != other.VType {
		res.SetBool(true)
		return
	}

	switch v.VType {
	case ValueTypeNil:
		res.SetBool(false)
	case ValueTypeInt:
		res.SetBool(v.GetInt() != other.GetInt())
	case ValueTypeFloat:
		res.SetBool(v.GetFloat() != other.GetFloat())
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
