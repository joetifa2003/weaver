package builtin

import (
	"fmt"
	"math/rand/v2"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *RegistryBuilder) {
	builder.RegisterFunc("echo", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		res := vm.Value{}

		val, err := args.Get(0)
		if err != nil {
			return res, err
		}

		fmt.Println(val.String())

		return res, nil
	})

	builder.RegisterFunc("rand", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		return vm.NewFloat(rand.Float64()), nil
	})

	builder.RegisterFunc("makeArr", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		res := vm.Value{}

		val, err := args.Get(0)
		if err != nil {
			return res, err
		}

		arr := make([]vm.Value, val.GetInt())
		res.SetArray(arr)

		return res, nil
	})

	builder.RegisterFunc("len", func(x *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		res := vm.Value{}

		val, err := args.Get(0, vm.ValueTypeArray, vm.ValueTypeString, vm.ValueTypeObject)
		if err != nil {
			return res, err
		}

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

	builder.RegisterFunc("assert", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		val, err := args.Get(0)
		if err != nil {
			return vm.Value{}, err
		}

		if !val.IsTruthy() {
			panic("assertion failed")
		}

		return vm.Value{}, nil
	})

	builder.RegisterFunc("type", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		val, err := args.Get(0)
		if err != nil {
			return vm.Value{}, err
		}

		return vm.NewString(val.VType.String()), nil
	})

	builder.RegisterFunc("string", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		val, err := args.Get(0)
		if err != nil {
			return vm.Value{}, err
		}

		return vm.NewString(val.String()), nil
	})

	builder.RegisterFunc("int", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		val, err := args.Get(0, vm.ValueTypeInt, vm.ValueTypeFloat)
		if err != nil {
			return vm.Value{}, err
		}

		switch val.VType {
		case vm.ValueTypeInt:
			return val, nil
		case vm.ValueTypeFloat:
			return vm.NewInt(int(val.GetFloat())), nil
		default:
			panic("unreachable")
		}
	})
}
