package builtin

import (
	"fmt"
	"math/rand/v2"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *RegistryBuilder) {
	builder.RegisterFunc("error", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		msgArg, err := args.Get(0, vm.ValueTypeString)
		if err != nil {
			return vm.Value{}, err
		}
		if len(args) == 1 {
			return vm.NewError(msgArg.GetString(), vm.Value{}), nil
		}

		dataArg, err := args.Get(1)
		if err != nil {
			return vm.Value{}, err
		}

		return vm.NewError(msgArg.GetString(), dataArg), nil
	})

	builder.RegisterFunc("isError", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		val, err := args.Get(0)
		if err != nil {
			return vm.Value{}, err
		}

		return vm.NewBool(val.VType == vm.ValueTypeError), nil
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
		return vm.NewNumber(rand.Float64()), nil
	})

	builder.RegisterFunc("len", func(x *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
		res := vm.Value{}

		val, err := args.Get(0, vm.ValueTypeArray, vm.ValueTypeString, vm.ValueTypeObject)
		if err != nil {
			return res, err
		}

		switch val.VType {
		case vm.ValueTypeArray:
			res.SetNumber(float64(len(*val.GetArray())))
		case vm.ValueTypeString:
			res.SetNumber(float64(len(val.GetString())))
		case vm.ValueTypeObject:
			res.SetNumber(float64(len(val.GetObject())))
		default:
			panic("unreachable")
		}

		return res, nil
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
		val, err := args.Get(0, vm.ValueTypeNumber)
		if err != nil {
			return vm.Value{}, err
		}

		switch val.VType {
		case vm.ValueTypeNumber:
			val.SetNumber(float64(int(val.GetNumber())))
			return val, nil
		default:
			panic("unreachable")
		}
	})
}
