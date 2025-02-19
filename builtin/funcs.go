package builtin

import (
	"fmt"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *RegistryBuilder) {
	builder.RegisterFunc("echo", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("echo() takes exactly 1 argument")
		}

		val := args[0]

		fmt.Println(val.String())

		return
	})

	builder.RegisterFunc("len", func(x *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("len() takes exactly 1 argument")
		}

		val := args[0]

		switch val.VType {
		case vm.ValueTypeArray:
			res.SetInt(len(*val.GetArray()))
		case vm.ValueTypeString:
			res.SetInt(len(val.GetString()))
		case vm.ValueTypeObject:
			res.SetInt(len(val.GetObject()))
		default:
			panic("len() argument must be an array, string or object")
		}

		return
	})

	builder.RegisterFunc("push", func(v *vm.VM, args ...vm.Value) vm.Value {
		if len(args) != 2 {
			panic("expected 1 arg")
		}

		arr := args[0].GetArray()
		val := args[1]

		*arr = append(*arr, val)

		return args[0]
	})

	builder.RegisterFunc("map", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 2 {
			panic("map() takes exactly 2 arguments")
		}

		arr := *args[0].GetArray()
		fn := args[1]

		newArr := make([]vm.Value, len(arr))
		for i, val := range arr {
			newArr[i] = v.RunFunction(fn, val)
		}

		var result vm.Value
		result.SetArray(newArr)

		return result
	})

	builder.RegisterFunc("filter", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 2 {
			panic("filter() takes exactly 2 arguments")
		}

		arr := *args[0].GetArray()
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
	})

	builder.RegisterFunc("contains", func(v *vm.VM, args ...vm.Value) vm.Value {
		if len(args) != 2 {
			panic("filter() takes exactly 2 arguments")
		}

		arr := args[0]
		if arr.VType != vm.ValueTypeArray {
			panic("filter() argument must be an array")
		}

		f := args[1]
		if f.VType == vm.ValueTypeFunction {
			for _, val := range *arr.GetArray() {
				if v.RunFunction(f, val).IsTruthy() {
					return vm.NewBool(true)
				}
			}
		} else {
			isEqual := vm.Value{}
			for _, val := range *arr.GetArray() {
				val.Equal(&f, &isEqual)
				if isEqual.IsTruthy() {
					return vm.NewBool(true)
				}
			}
		}

		return vm.NewBool(false)
	})

	builder.RegisterFunc("assert", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("assert() takes exactly 1 argument")
		}

		val := args[0]
		if !val.IsTruthy() {
			panic("assertion failed")
		}

		return
	})

	builder.RegisterFunc("type", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("type() takes exactly 1 argument")
		}

		val := args[0]
		res.SetString(val.VType.String())

		return
	})

	builder.RegisterFunc("string", func(v *vm.VM, args ...vm.Value) (res vm.Value) {
		if len(args) != 1 {
			panic("string() takes exactly 1 argument")
		}

		val := args[0]
		res.SetString(val.String())

		return
	})
}
