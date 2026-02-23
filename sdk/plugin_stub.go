//go:build !wasip1

package sdk

// Legacy convenience functions

func ArgLen() int {
	panic("sdk can only be used in a wasip1 environment")
}

func ArgNum(idx int) float64 {
	panic("sdk can only be used in a wasip1 environment")
}

func ArgStr(idx int) string {
	panic("sdk can only be used in a wasip1 environment")
}

func ReturnNum(val float64) {
	panic("sdk can only be used in a wasip1 environment")
}

func ReturnStr(val string) {
	panic("sdk can only be used in a wasip1 environment")
}

func ReturnBool(val bool) {
	panic("sdk can only be used in a wasip1 environment")
}

func ReturnNil() {
	panic("sdk can only be used in a wasip1 environment")
}

// Handle-based Value API

type Value struct{ handle int32 }

const (
	TypeNil    = 0
	TypeNumber = 1
	TypeString = 2
	TypeObject = 3
	TypeBool   = 4
	TypeArray  = 6
)

func Nil() Value {
	panic("sdk can only be used in a wasip1 environment")
}

func Bool(b bool) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func Number(n float64) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func String(s string) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func Object() Value {
	panic("sdk can only be used in a wasip1 environment")
}

func Array() Value {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Set(key string, val Value) {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Get(key string) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Push(val Value) {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Index(i int) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Len() int {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) Type() int {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) AsNumber() float64 {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) AsString() string {
	panic("sdk can only be used in a wasip1 environment")
}

func (v Value) AsBool() bool {
	panic("sdk can only be used in a wasip1 environment")
}

func Arg(idx int) Value {
	panic("sdk can only be used in a wasip1 environment")
}

func Return(v Value) {
	panic("sdk can only be used in a wasip1 environment")
}

func ReturnErrorMsg(msg string) {
	panic("sdk can only be used in a wasip1 environment")
}
