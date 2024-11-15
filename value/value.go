package value

import (
	"fmt"
	"unsafe"

	"github.com/joetifa2003/weaver/opcode"
)

type ValueType int

const (
	ValueTypeInt ValueType = iota
	ValueTypeFloat
	ValueTypeString
	ValueTypeObject
	ValueTypeBool
	ValueTypeFunction
	ValueTypeNil
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
	default:
		panic(fmt.Sprintf("unimplemented %T", t))
	}
}

// Value poor mans union/enum
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

func (v *Value) GetObject() map[string]Value {
	if v.VType != ValueTypeString {
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
}

func (v *Value) SetFunction(f FunctionValue) {
	v.VType = ValueTypeFunction
	v.nonPrimitive = unsafe.Pointer(&f)
}

func (v *Value) GetFunction() FunctionValue {
	if v.VType != ValueTypeFunction {
		panic("Value.GetFunction(): not a function")
	}

	return *(*FunctionValue)(v.nonPrimitive)
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

func (v Value) String() string {
	switch v.VType {
	case ValueTypeString:
		str := v.GetString()
		return str
	case ValueTypeInt:
		integer := v.GetInt()
		return fmt.Sprint(integer)
	case ValueTypeFloat:
		float := v.GetFloat()
		return fmt.Sprint(float)
	case ValueTypeNil:
		return "nil"
	case ValueTypeObject:
		return "object"
	case ValueTypeBool:
		return fmt.Sprint(v.GetBool())
	case ValueTypeFunction:
		return "function"
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
	default:
		panic(fmt.Sprintf("Value.IsTruthy(): unimplemented %T", v.VType))
	}
}
