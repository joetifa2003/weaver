package builtin

import (
	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncsArr(builder *vm.RegistryBuilder) {
	builder.RegisterFunc("makeArr", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		length := 0
		capacity := 0

		lengthArg, ok := args.Get(0)
		if !ok {
			return lengthArg
		}

		length = int(lengthArg.GetNumber())

		if args.Len() > 1 {
			capacityArg, ok := args.Get(1)
			if !ok {
				return capacityArg
			}
			capacity = int(capacityArg.GetNumber())
			return vm.NewArray(make([]vm.Value, length, capacity))
		}

		return vm.NewArray(make([]vm.Value, length))
	})

	builder.RegisterFunc("push", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg
		}

		val, ok := args.Get(1)
		if !ok {
			return val
		}

		arr := arrArg.GetArray()
		*arr = append(*arr, val)
		return arrArg
	})

	builder.RegisterFunc("map", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg
		}

		fnArg, ok := args.Get(1, vm.ValueTypeFunction)
		if !ok {
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
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg
		}

		fnArg, ok := args.Get(1, vm.ValueTypeFunction)
		if !ok {
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
		arrArg, ok := args.Get(0, vm.ValueTypeArray)
		if !ok {
			return arrArg
		}

		f, ok := args.Get(1)
		if !ok {
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
