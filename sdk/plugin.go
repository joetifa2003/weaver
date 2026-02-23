//go:build wasip1

package sdk

import (
	"runtime"
	"unsafe"
)

//go:wasmimport weaver get_arg_len
func getArgLen() int32

// Handle-based builder API imports

//go:wasmimport weaver value_new_nil
func valueNewNil() int32

//go:wasmimport weaver value_new_bool
func valueNewBool(val int32) int32

//go:wasmimport weaver value_new_number
func valueNewNumber(val float64) int32

//go:wasmimport weaver value_new_string
func valueNewString(ptr int32, length int32) int32

//go:wasmimport weaver value_new_object
func valueNewObject() int32

//go:wasmimport weaver value_new_array
func valueNewArray() int32

//go:wasmimport weaver value_object_set
func valueObjectSet(objHandle int32, keyPtr int32, keyLen int32, valHandle int32)

//go:wasmimport weaver value_array_push
func valueArrayPush(arrHandle int32, valHandle int32)

//go:wasmimport weaver get_arg_value
func getArgValue(idx int32) int32

//go:wasmimport weaver value_type
func valueType(handle int32) int32

//go:wasmimport weaver value_get_number
func valueGetNumber(handle int32) float64

//go:wasmimport weaver value_get_string
func valueGetString(handle int32, ptr int32, maxLen int32) int32

//go:wasmimport weaver value_get_bool
func valueGetBool(handle int32) int32

//go:wasmimport weaver value_object_get
func valueObjectGet(objHandle int32, keyPtr int32, keyLen int32) int32

//go:wasmimport weaver value_array_get
func valueArrayGet(arrHandle int32, idx int32) int32

//go:wasmimport weaver value_array_len
func valueArrayLen(arrHandle int32) int32

//go:wasmimport weaver return_value
func returnValue(handle int32)

//go:wasmimport weaver return_error
func returnError(ptr int32, length int32)

// === Legacy convenience functions (delegate to handle API) ===

func ArgLen() int {
	return int(getArgLen())
}

func ArgNum(idx int) float64 {
	return Arg(idx).AsNumber()
}

func ArgStr(idx int) string {
	return Arg(idx).AsString()
}

func ReturnNum(val float64) {
	Return(Number(val))
}

func ReturnStr(val string) {
	Return(String(val))
}

func ReturnBool(val bool) {
	Return(Bool(val))
}

func ReturnNil() {
	Return(Nil())
}

// === Handle-based Value API ===

// Value is an opaque handle to a Weaver value on the host.
type Value struct{ handle int32 }

// Value type constants (matching vm.ValueType iota order)
const (
	TypeNil    = 0
	TypeNumber = 1
	TypeString = 2
	TypeObject = 3
	TypeBool   = 4
	TypeArray  = 6
)

// Constructors

// Nil creates a nil value handle.
func Nil() Value {
	return Value{valueNewNil()}
}

// Bool creates a bool value handle.
func Bool(b bool) Value {
	v := int32(0)
	if b {
		v = 1
	}
	return Value{valueNewBool(v)}
}

// Number creates a number value handle.
func Number(n float64) Value {
	return Value{valueNewNumber(n)}
}

// String creates a string value handle.
func String(s string) Value {
	if len(s) == 0 {
		return Value{valueNewString(0, 0)}
	}
	bytes := []byte(s)
	ptr := int32(uintptr(unsafe.Pointer(&bytes[0])))
	h := valueNewString(ptr, int32(len(bytes)))
	runtime.KeepAlive(bytes)
	return Value{h}
}

// Object creates an empty object value handle.
func Object() Value {
	return Value{valueNewObject()}
}

// Array creates an empty array value handle.
func Array() Value {
	return Value{valueNewArray()}
}

// Object methods

// Set sets a field on an object value.
func (v Value) Set(key string, val Value) {
	if len(key) == 0 {
		return
	}
	bytes := []byte(key)
	ptr := int32(uintptr(unsafe.Pointer(&bytes[0])))
	valueObjectSet(v.handle, ptr, int32(len(bytes)), val.handle)
	runtime.KeepAlive(bytes)
}

// Get gets a field from an object value. Returns a Value with handle -1 if not found.
func (v Value) Get(key string) Value {
	if len(key) == 0 {
		return Value{-1}
	}
	bytes := []byte(key)
	ptr := int32(uintptr(unsafe.Pointer(&bytes[0])))
	h := valueObjectGet(v.handle, ptr, int32(len(bytes)))
	runtime.KeepAlive(bytes)
	return Value{h}
}

// Array methods

// Push appends a value to an array.
func (v Value) Push(val Value) {
	valueArrayPush(v.handle, val.handle)
}

// Index gets an element from an array by index.
func (v Value) Index(i int) Value {
	return Value{valueArrayGet(v.handle, int32(i))}
}

// Len returns the length of an array.
func (v Value) Len() int {
	return int(valueArrayLen(v.handle))
}

// Reading values

// Type returns the ValueType of this handle.
func (v Value) Type() int {
	return int(valueType(v.handle))
}

// AsNumber reads the number value from this handle.
func (v Value) AsNumber() float64 {
	return valueGetNumber(v.handle)
}

// AsString reads the string value from this handle.
func (v Value) AsString() string {
	l := valueGetString(v.handle, 0, 0)
	if l <= 0 {
		return ""
	}
	buf := make([]byte, l)
	ptr := int32(uintptr(unsafe.Pointer(&buf[0])))
	valueGetString(v.handle, ptr, l)
	runtime.KeepAlive(buf)
	return string(buf)
}

// AsBool reads the bool value from this handle.
func (v Value) AsBool() bool {
	return valueGetBool(v.handle) != 0
}

// Args & Return

// Arg gets a function argument as a Value handle.
func Arg(idx int) Value {
	return Value{getArgValue(int32(idx))}
}

// Return sets the return value from a Value handle.
func Return(v Value) {
	returnValue(v.handle)
}

// ReturnErrorMsg returns an error with the given message.
func ReturnErrorMsg(msg string) {
	if len(msg) == 0 {
		return
	}
	bytes := []byte(msg)
	ptr := int32(uintptr(unsafe.Pointer(&bytes[0])))
	returnError(ptr, int32(len(bytes)))
	runtime.KeepAlive(bytes)
}
