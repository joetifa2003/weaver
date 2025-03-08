package builtin

import "github.com/joetifa2003/weaver/vm"

func registerBuiltinFuncsArr(builder *RegistryBuilder) {
	builder.RegisterFunc("makeArr", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		if val.IsError() {
			return val
		}

		res := vm.Value{}
		arr := make([]vm.Value, int(val.GetNumber()))
		res.SetArray(arr)
		return res
	})

	builder.RegisterFunc("push", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg := args.Get(0, vm.ValueTypeArray)
		if arrArg.IsError() {
			return arrArg
		}

		val := args.Get(1)
		if val.IsError() {
			return val
		}

		arr := arrArg.GetArray()
		*arr = append(*arr, val)
		return arrArg
	})

	builder.RegisterFunc("map", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg := args.Get(0, vm.ValueTypeArray)
		if arrArg.IsError() {
			return arrArg
		}

		fnArg := args.Get(1, vm.ValueTypeFunction)
		if fnArg.IsError() {
			return fnArg
		}

		arr := *arrArg.GetArray()
		newArr := make([]vm.Value, len(arr))
		for i, val := range arr {
			newArr[i] = v.RunFunction(fnArg, val)
		}

		var result vm.Value
		result.SetArray(newArr)
		return result
	})

	builder.RegisterFunc("filter", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg := args.Get(0, vm.ValueTypeArray)
		if arrArg.IsError() {
			return arrArg
		}

		fnArg := args.Get(1, vm.ValueTypeFunction)
		if fnArg.IsError() {
			return fnArg
		}

		arr := *arrArg.GetArray()
		newArr := make([]vm.Value, 0)
		for _, val := range arr {
			r := v.RunFunction(fnArg, val)
			if r.IsTruthy() {
				newArr = append(newArr, val)
			}
		}

		var result vm.Value
		result.SetArray(newArr)
		return result
	})

	builder.RegisterFunc("contains", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg := args.Get(0, vm.ValueTypeArray)
		if arrArg.IsError() {
			return arrArg
		}

		f := args.Get(1)
		if f.IsError() {
			return f
		}

		arr := *arrArg.GetArray()
		if f.VType == vm.ValueTypeFunction {
			for _, val := range arr {
				r := v.RunFunction(f, val)
				if r.IsTruthy() {
					return vm.NewBool(true)
				}
			}
		} else {
			isEqual := vm.Value{}
			for _, val := range arr {
				val.Equal(&f, &isEqual)
				if isEqual.IsTruthy() {
					return vm.NewBool(true)
				}
			}
		}

		return vm.NewBool(false)
	})
}
