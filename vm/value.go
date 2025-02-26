package vm

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/joetifa2003/weaver/opcode"
)

type ValueType uint8

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
	ValueTypeModule
	ValueTypeNativeObject
	ValueTypeRef
	valueTypeEnd
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
	case ValueTypeNativeObject:
		return "native object"
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

func (v *Value) GetInt() int {
	switch v.VType {
	case ValueTypeInt:
		return *interpret[int](&v.primitive)
	case ValueTypeFloat:
		return int(*interpret[float64](&v.primitive))
	default:
		return 0
	}
}

func (v *Value) GetFloat() float64 {
	switch v.VType {
	case ValueTypeInt:
		return float64(*interpret[int](&v.primitive))
	case ValueTypeFloat:
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

func NewInt(i int) Value {
	val := Value{}
	val.SetInt(i)
	return val
}

func NewFloat(f float64) Value {
	val := Value{}
	val.SetFloat(f)
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

func (v *Value) String() string {
	switch v.VType {
	case ValueTypeModule:
		return "module"

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

var addTable = initOpTable("+",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetInt(v.GetInt() + other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() + other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() + other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() + other.GetFloat())
	}},
	opDef{ValueTypeString, ValueTypeString, func(v *Value, other *Value, res *Value) {
		res.SetString(v.GetString() + other.GetString())
	}},
)

func (v *Value) Add(other *Value, res *Value) {
	addTable.Call(v, other, res)
}

var subTable = initOpTable("-",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetInt(v.GetInt() - other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() - other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() - other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() - other.GetFloat())
	}},
)

func (v *Value) Sub(other *Value, res *Value) {
	subTable.Call(v, other, res)
}

var mulTable = initOpTable("*",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetInt(v.GetInt() * other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() * other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() * other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() * other.GetFloat())
	}},
)

func (v *Value) Mul(other *Value, res *Value) {
	mulTable.Call(v, other, res)
}

var divTable = initOpTable("/",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetInt(v.GetInt() / other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() / other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() / other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetFloat(v.GetFloat() / other.GetFloat())
	}},
)

func (v *Value) Div(other *Value, res *Value) {
	divTable.Call(v, other, res)
}

var lessThanTable = initOpTable("<",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetInt() < other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() < other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() < other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() < other.GetFloat())
	}},
)

func (v *Value) LessThan(other *Value, res *Value) {
	lessThanTable.Call(v, other, res)
}

var lessThanEqualTable = initOpTable("<=",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetInt() <= other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() <= other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() <= other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() <= other.GetFloat())
	}},
)

func (v *Value) LessThanEqual(other *Value, res *Value) {
	lessThanEqualTable.Call(v, other, res)
}

var greaterThanTable = initOpTable(">",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetInt() > other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() > other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() > other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() > other.GetFloat())
	}},
)

func (v *Value) GreaterThan(other *Value, res *Value) {
	greaterThanTable.Call(v, other, res)
}

var greaterThanEqual = initOpTable(">=",
	opDef{ValueTypeInt, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetInt() >= other.GetInt())
	}},
	opDef{ValueTypeInt, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() >= other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeInt, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() >= other.GetFloat())
	}},
	opDef{ValueTypeFloat, ValueTypeFloat, func(v *Value, other *Value, res *Value) {
		res.SetBool(v.GetFloat() >= other.GetFloat())
	}},
)

func (v *Value) GreaterThanEqual(other *Value, res *Value) {
	greaterThanEqual.Call(v, other, res)
}

func (v *Value) Equal(other *Value, res *Value) {
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

func (v *Value) NotEqual(other *Value, res *Value) {
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

func (v *Value) Negate(res *Value) {
	switch v.VType {
	case ValueTypeInt:
		res.SetInt(-v.GetInt())
	case ValueTypeFloat:
		res.SetFloat(-v.GetFloat())
	default:
		panic(fmt.Sprintf("illegal operation -%s", v))
	}
}

func (v *Value) Mod(other *Value, res *Value) {
	if (v.VType != ValueTypeInt && v.VType != ValueTypeFloat) || (other.VType != ValueTypeInt && other.VType != ValueTypeFloat) {
		panic(fmt.Sprintf("illegal operation %s %% %s", v, other))
	}

	res.SetInt(v.GetInt() % other.GetInt())
}
