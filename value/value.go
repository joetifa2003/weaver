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

func (v *Value) getInt() (*int, error) {
	if v.VType != ValueTypeInt {
		return nil, ErrInvalidType
	}

	return interpret[int](v.primitive[:]), nil
}

func (v *Value) getFloat() (*float64, error) {
	if v.VType != ValueTypeFloat {
		return nil, ErrInvalidType
	}

	return interpret[float64](v.primitive[:]), nil
}

func (v *Value) getString() (*string, error) {
	if v.VType != ValueTypeString {
		return nil, ErrInvalidType
	}

	return (*string)(v.nonPrimitive), nil
}

func (v *Value) getObject() (*map[string]Value, error) {
	if v.VType != ValueTypeString {
		return nil, ErrInvalidType
	}

	return (*map[string]Value)(v.nonPrimitive), nil
}

func (v *Value) String() string {
	switch v.VType {
	case ValueTypeString:
		str, _ := v.getString()
		return *str
	case ValueTypeInt:
		integer, _ := v.getInt()
		return fmt.Sprint(*integer)
	case ValueTypeFloat:
		float, _ := v.getFloat()
		return fmt.Sprint(*float)
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
