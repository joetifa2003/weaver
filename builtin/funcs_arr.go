package builtin

import (
	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncsArr(builder *vm.RegistryBuilder) {
	builder.RegisterFunc("makeArr", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		length := 0
		capacity := 0

		lengthArg, ok := args.Get(0)
		if !ok {
			return lengthArg, false
		}

		length = int(lengthArg.GetNumber())

		if args.Len() > 1 {
			capacityArg, ok := args.Get(1)
			if !ok {
				return capacityArg, false
			}
			capacity = int(capacityArg.GetNumber())
			return vm.NewArray(make([]vm.Value, length, capacity)), true
		}

		return vm.NewArray(make([]vm.Value, length)), true
	})

	builder.RegisterFunc("push", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		val, ok := args.Get(1)
		if !ok {
			return val, false
		}

		arr := arrArg.GetArray()
		*arr = append(*arr, val)
		return arrArg, true
	})

	builder.RegisterFunc("map", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		fnArg, ok := args.Get(1, vm.ValueTypeFunction)
		if !ok {
			return fnArg, false
		}

		arr := *arrArg.GetArray()
		newArr := make([]vm.Value, len(arr))
		for i, val := range arr {
			mapped, ok := v.RunFunction(fnArg, val)
			if !ok {
				return mapped, false
			}

			newArr[i] = mapped
		}

		var result vm.Value
		result.SetArray(newArr)
		return result, true
	})

	builder.RegisterFunc("filter", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		fnArg, ok := args.Get(1, vm.ValueTypeFunction)
		if !ok {
			return fnArg, false
		}

		arr := *arrArg.GetArray()
		newArr := make([]vm.Value, 0)
		for _, val := range arr {
			r, ok := v.RunFunction(fnArg, val)
			if !ok {
				return r, false
			}

			if r.IsTruthy() {
				newArr = append(newArr, val)
			}
		}

		var result vm.Value
		result.SetArray(newArr)
		return result, true
	})

	builder.RegisterFunc("contains", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		f, ok := args.Get(1)
		if !ok {
			return f, false
		}

		arr := *arrArg.GetArray()
		if f.VType == vm.ValueTypeFunction {
			for _, val := range arr {
				r, ok := v.RunFunction(f, val)
				if !ok {
					return r, false
				}
				if r.IsTruthy() {
					return vm.NewBool(true), true
				}
			}
		} else {
			isEqual := vm.Value{}
			for _, val := range arr {
				val.Equal(&f, &isEqual)
				if isEqual.IsTruthy() {
					return vm.NewBool(true), true
				}
			}
		}

		return vm.NewBool(false), true
	})

	builder.RegisterFunc("find", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		f, ok := args.Get(1)
		if !ok {
			return f, false
		}

		arr := *arrArg.GetArray()
		if f.VType == vm.ValueTypeFunction {
			for _, val := range arr {
				r, ok := v.RunFunction(f, val)
				if !ok {
					return r, false
				}

				if r.IsTruthy() {
					return val, true
				}
			}
		} else {
			isEqual := vm.Value{}
			for _, val := range arr {
				val.Equal(&f, &isEqual)
				if isEqual.IsTruthy() {
					return val, true
				}
			}
		}

		return vm.Value{}, true
	})

	builder.RegisterFunc("each", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg, false
		}

		fnArg, ok := args.Get(1, vm.ValueTypeFunction)
		if !ok {
			return fnArg, false
		}

		arr := *arrArg.GetArray()
		for _, val := range arr {
			r, ok := v.RunFunction(fnArg, val)
			if !ok {
				return r, false
			}
		}
		return vm.Value{}, true
	})
}
