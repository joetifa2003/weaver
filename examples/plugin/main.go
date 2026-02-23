package main

import (
	"github.com/joetifa2003/weaver/sdk"
)

//go:wasmexport greet
func greet() {
	name := sdk.ArgStr(0)
	sdk.ReturnStr("Hello, " + name + " from Wasm!")
}

//go:wasmexport sum
func sum() {
	a := sdk.ArgNum(0)
	b := sdk.ArgNum(1)
	sdk.ReturnNum(a + b)
}

//go:wasmexport create_user
func createUser() {
	name := sdk.Arg(0).AsString()
	age := sdk.Arg(1).AsNumber()

	user := sdk.Object()
	user.Set("name", sdk.String(name))
	user.Set("age", sdk.Number(age))
	user.Set("active", sdk.Bool(true))

	tags := sdk.Array()
	tags.Push(sdk.String("admin"))
	tags.Push(sdk.String("user"))
	user.Set("tags", tags)

	sdk.Return(user)
}

//go:wasmexport create_error
func createError() {
	sdk.ReturnErrorMsg("something went wrong")
}

func main() {}
