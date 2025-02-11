package vm

import (
	"fmt"
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
	ValueTypRef
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

// Value poor mans union/enum
type Value struct {
	VType        ValueType
	primitive    [8]byte
	nonPrimitive unsafe.Pointer
}

func interpret[T any](b *[8]byte) *T {
	return (*T)(unsafe.Pointer(b))
}

func (v *Value) Deref() *Value {
	if v.VType == ValueTypRef {
		return (*Value)(v.nonPrimitive)
	}

	return v
}

func (v *Value) SetRef(r *Value) {
	v.VType = ValueTypRef
	v.nonPrimitive = unsafe.Pointer(r)
}

func (v *Value) GetInt() int {
	vr := v.Deref()

	if vr.VType != ValueTypeInt {
		panic("Value.GetInt(): not an int")
	}

	return *interpret[int](&vr.primitive)
}

func (v *Value) GetFloat() float64 {
	vr := v.Deref()

	if vr.VType != ValueTypeFloat {
		panic("Value.GetFloat(): not a float")
	}

	return *interpret[float64](&vr.primitive)
}

func (v *Value) SetObject(o map[string]Value) {
	vr := v.Deref()
	vr.VType = ValueTypeObject
	vr.nonPrimitive = unsafe.Pointer(&o)
}

func (v *Value) GetObject() map[string]Value {
	vr := v.Deref()

	if vr.VType != ValueTypeObject {
		panic("Value.GetObject(): not an object")
	}

	return *(*map[string]Value)(vr.nonPrimitive)
}

func (v *Value) GetBool() bool {
	vr := v.Deref()

	if vr.VType != ValueTypeBool {
		panic("Value.GetBool(): not a bool")
	}

	return *interpret[bool](&vr.primitive)
}

func (v *Value) GetArray() *[]Value {
	vr := v.Deref()

	if vr.VType != ValueTypeArray {
		panic("Value.GetArray(): not an array")
	}

	return (*[]Value)(v.nonPrimitive)
}

func (v *Value) SetArray(a []Value) {
	vr := v.Deref()
	vr.VType = ValueTypeArray
	vr.nonPrimitive = unsafe.Pointer(&a)
}

func (v *Value) SetBool(b bool) {
	vr := v.Deref()
	vr.VType = ValueTypeBool
	*interpret[bool](&vr.primitive) = b
}

func (v *Value) SetInt(i int) {
	vr := v.Deref()
	vr.VType = ValueTypeInt
	*interpret[int](&vr.primitive) = i
}

func (v *Value) SetFloat(f float64) {
	vr := v.Deref()
	vr.VType = ValueTypeFloat
	*interpret[float64](&vr.primitive) = f
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
	vr := v.Deref()
	vr.VType = ValueTypeFunction
	vr.nonPrimitive = unsafe.Pointer(&f)
}

func (v *Value) GetFunction() *FunctionValue {
	vr := v.Deref()
	if vr.VType != ValueTypeFunction {
		panic("Value.GetFunction(): not a function")
	}

	return (*FunctionValue)(vr.nonPrimitive)
}

func (v *Value) SetString(s string) {
	vr := v.Deref()
	vr.VType = ValueTypeString
	vr.nonPrimitive = unsafe.Pointer(&s)
}

func (v *Value) GetString() string {
	vr := v.Deref()
	if vr.VType != ValueTypeString {
		panic("Value.GetString(): not a string")
	}

	return *(*string)(vr.nonPrimitive)
}

type NativeFunction func(v *VM, args ...Value) Value

func (v *Value) GetNativeFunction() NativeFunction {
	vr := v.Deref()
	if vr.VType != ValueTypeNativeFunction {
		panic("Value.GetNativeFunction(): not a native function")
	}

	return *(*NativeFunction)(vr.nonPrimitive)
}

func (v *Value) SetNativeFunction(f NativeFunction) {
	vr := v.Deref()
	vr.VType = ValueTypeNativeFunction
	vr.nonPrimitive = unsafe.Pointer(&f)
}

func (v Value) String() string {
	vr := v.Deref()
	switch vr.VType {
	case ValueTypeString:
		str := vr.GetString()
		return str

	case ValueTypeInt:
		integer := vr.GetInt()
		return fmt.Sprint(integer)

	case ValueTypeFloat:
		float := vr.GetFloat()
		return fmt.Sprint(float)

	case ValueTypeNil:
		return "nil"

	case ValueTypeObject:
		return fmt.Sprint(vr.GetObject())

	case ValueTypeBool:
		return fmt.Sprint(vr.GetBool())

	case ValueTypeFunction:
		return "function"

	case ValueTypeArray:
		arr := *vr.GetArray()
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

func (v Value) Add(other Value, res *Value) {
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

func (v Value) Sub(other Value, res *Value) {
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

func (v Value) Mul(other Value, res *Value) {
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

func (v Value) Div(other Value, res *Value) {
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

func (v Value) Mod(other Value, res *Value) {
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

func (v Value) LessThan(other Value, res *Value) {
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

func (v Value) LessThanEqual(other Value, res *Value) {
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

func (v Value) GreaterThan(other Value, res *Value) {
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

func (v Value) GreaterThanEqual(other Value, res *Value) {
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

// TODO: implement object equality
func (v Value) Equal(other Value, res *Value) {
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

func (v Value) NotEqual(other Value, res *Value) {
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
