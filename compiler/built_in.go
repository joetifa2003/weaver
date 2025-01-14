package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/value"
)

var builtInFunctions = map[string]value.NativeFunction{
	"echo": func(args ...value.Value) (res value.Value) {
		if len(args) != 1 {
			panic("echo() takes exactly 1 argument")
		}

		v := args[0]

		fmt.Println(v.String())

		return
	},
	"len": func(args ...value.Value) (res value.Value) {
		if len(args) != 1 {
			panic("len() takes exactly 1 argument")
		}

		v := args[0]

		switch v.VType {
		case value.ValueTypeArray:
			res.SetInt(len(v.GetArray()))
		case value.ValueTypeString:
			res.SetInt(len(v.GetString()))
		case value.ValueTypeObject:
			res.SetInt(len(v.GetObject()))
		default:
			panic("len() argument must be an array, string or object")
		}

		return
	},
	// "map": func(args ...value.Value) (res value.Value) {
	// 	if len(args) != 2 {
	// 		panic("map() takes exactly 2 arguments")
	// 	}
	//
	// 	arr := args[0].GetArray()
	// 	fn := args[1].GetFunction()
	//
	// 	newArr := make([]value.Value, len(arr))
	// 	for i, v := range arr {
	// 		newArr[i] =
	// 	}
	// },
}
