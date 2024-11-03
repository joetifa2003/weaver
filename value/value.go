package value

import (
	"fmt"
	"unsafe"
)

type ValueType int

const (
	ValueTypeInt ValueType = iota
	ValueTypeFloat
	ValueTypeString
	ValueTypeObject
)

// Value poor mans union/enum
type Value struct {
	VType        ValueType
	primitive    [8]byte
	nonPrimitive unsafe.Pointer
}

func interpret[T any](b []byte) *T {
	return (*T)(unsafe.Pointer(&b[0]))
}

func (v *Value) GetInt() int {
	if v.VType != ValueTypeInt {
		panic("Value.GetInt(): not an int")
	}

	return *interpret[int](v.primitive[:])
}

func (v *Value) GetFloat() float64 {
	if v.VType != ValueTypeFloat {
		panic("Value.GetFloat(): not a float")
	}

	return *interpret[float64](v.primitive[:])
}

func (v *Value) GetString() string {
	if v.VType != ValueTypeString {
		panic("Value.GetString(): not a string")
	}

	return *(*string)(v.nonPrimitive)
}

func (v *Value) GetObject() map[string]Value {
	if v.VType != ValueTypeString {
		panic("Value.GetObject(): not an object")
	}

	return *(*map[string]Value)(v.nonPrimitive)
}

func (v *Value) String() string {
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
	default:
		panic(fmt.Sprintf("Value.String(): unimplemented %T", v.VType))
	}
}

func NewString(s string) Value {
	return Value{
		VType:        ValueTypeString,
		nonPrimitive: unsafe.Pointer(&s),
	}
}

func NewInt(i int) Value {
	v := Value{
		VType: ValueTypeInt,
	}
	p := interpret[int](v.primitive[:])
	*p = i

	return v
}

func NewObject(o map[string]Value) Value {
	return Value{
		VType:        ValueTypeObject,
		nonPrimitive: unsafe.Pointer(&o),
	}
}
