package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/vm"
)

var builtInFunctions = map[string]vm.NativeFunction{
	"echo": func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("echo() takes exactly 1 argument")
		}

		val := args[0]

		fmt.Println(val.String())

		return
	},
	"len": func(x *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("len() takes exactly 1 argument")
		}

		val := args[0]

		switch val.VType {
		case vm.ValueTypeArray:
			res.SetInt(len(val.GetArray()))
		case vm.ValueTypeString:
			res.SetInt(len(val.GetString()))
		case vm.ValueTypeObject:
			res.SetInt(len(val.GetObject()))
		default:
			panic("len() argument must be an array, string or object")
		}

		return
	},
	"map": func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 2 {
			panic("map() takes exactly 2 arguments")
		}

		arr := args[0].GetArray()
		fn := args[1]

		newArr := make([]vm.Value, len(arr))
		for i, val := range arr {
			newArr[i] = v.RunFunction(fn, val)
		}

		var result vm.Value
		result.SetArray(newArr)

		return result
	},
	"filter": func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 2 {
			panic("filter() takes exactly 2 arguments")
		}

		arr := args[0].GetArray()
		fn := args[1]

		newArr := make([]vm.Value, 0)
		for _, val := range arr {
			if v.RunFunction(fn, val).IsTruthy() {
				newArr = append(newArr, val)
			}
		}

		var result vm.Value
		result.SetArray(newArr)

		return result
	},
	"assert": func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("assert() takes exactly 1 argument")
		}

		val := args[0]
		if !val.IsTruthy() {
			panic("assertion failed")
		}

		return
	},
}
