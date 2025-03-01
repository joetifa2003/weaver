package builtin

import "github.com/joetifa2003/weaver/vm"

func registerBuiltinFuncsArr(builder *RegistryBuilder) {
	builder.RegisterFunc("makeArr", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		res := vm.Value{}

		val, err := args.Get(0)
		if err != nil {
			return res, err
		}

		arr := make([]vm.Value, int(val.GetNumber()))
		res.SetArray(arr)

		return res, nil
	})

	builder.RegisterFunc("push", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		arrArg, err := args.Get(0, vm.ValueTypeArray)
		if err != nil {
			return vm.Value{}, err
		}

		val, err := args.Get(1)
		if err != nil {
			return vm.Value{}, err
		}

		arr := arrArg.GetArray()
		*arr = append(*arr, val)

		return arrArg, nil
	})

	builder.RegisterFunc("map", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		arrArg, err := args.Get(0, vm.ValueTypeArray)
		if err != nil {
			return vm.Value{}, err
		}
		fnArg, err := args.Get(1, vm.ValueTypeFunction)
		if err != nil {
			return vm.Value{}, err
		}

		arr := *arrArg.GetArray()

		newArr := make([]vm.Value, len(arr))
		for i, val := range arr {
			newArr[i] = v.RunFunction(fnArg, val)
		}

		var result vm.Value
		result.SetArray(newArr)

		return result, nil
	})

	builder.RegisterFunc("filter", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		arrArg, err := args.Get(0, vm.ValueTypeArray)
		if err != nil {
			return vm.Value{}, err
		}
		fnArg, err := args.Get(1, vm.ValueTypeFunction)
		if err != nil {
			return vm.Value{}, err
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

		return result, nil
	})

	builder.RegisterFunc("contains", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		arrArg, err := args.Get(0, vm.ValueTypeArray)
		if err != nil {
			return vm.Value{}, err
		}
		f, err := args.Get(1)
		if err != nil {
			return vm.Value{}, err
		}

		arr := *arrArg.GetArray()

		if f.VType == vm.ValueTypeFunction {
			for _, val := range arr {
				r := v.RunFunction(f, val)
				if r.IsTruthy() {
					return vm.NewBool(true), nil
				}
			}
		} else {
			isEqual := vm.Value{}
			for _, val := range arr {
				val.Equal(&f, &isEqual)
				if isEqual.IsTruthy() {
					return vm.NewBool(true), nil
				}
			}
		}

		return vm.NewBool(false), nil
	})
}
