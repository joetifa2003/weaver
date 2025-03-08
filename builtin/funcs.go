package builtin

import (
	"fmt"
	"math/rand/v2"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *RegistryBuilder) {
	builder.RegisterFunc("error", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		msgArg := args.Get(0, vm.ValueTypeString)
		if msgArg.IsError() {
			return msgArg
		}
		if len(args) == 1 {
			return vm.NewError(msgArg.GetString(), vm.Value{})
		}

		dataArg := args.Get(1)
		if dataArg.IsError() {
			return dataArg
		}

		return vm.NewError(msgArg.GetString(), dataArg)
	})

	builder.RegisterFunc("isError", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		return vm.NewBool(val.IsError())
	})

	builder.RegisterFunc("assert", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		if val.IsError() {
			return val
		}

		if !val.IsTruthy() {
			return vm.NewError("assertion failed", vm.Value{})
		}

		return vm.Value{}
	})

	builder.RegisterFunc("echo", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		if val.IsError() {
			return val
		}

		fmt.Println(val.String())
		return vm.Value{}
	})

	builder.RegisterFunc("rand", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		return vm.NewNumber(rand.Float64())
	})

	builder.RegisterFunc("len", func(x *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0, vm.ValueTypeArray, vm.ValueTypeString, vm.ValueTypeObject)
		if val.IsError() {
			return val
		}

		res := vm.Value{}
		switch val.VType {
		case vm.ValueTypeArray:
			res.SetNumber(float64(len(*val.GetArray())))
		case vm.ValueTypeString:
			res.SetNumber(float64(len(val.GetString())))
		case vm.ValueTypeObject:
			res.SetNumber(float64(len(val.GetObject())))
		default:
			return vm.NewError("invalid type for len()", vm.Value{})
		}

		return res
	})

	builder.RegisterFunc("type", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		return vm.NewString(val.VType.String())
	})

	builder.RegisterFunc("string", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0)
		if val.IsError() {
			return val
		}

		return vm.NewString(val.String())
	})

	builder.RegisterFunc("int", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val := args.Get(0, vm.ValueTypeNumber)
		if val.IsError() {
			return val
		}

		switch val.VType {
		case vm.ValueTypeNumber:
			val.SetNumber(float64(int(val.GetNumber())))
			return val
		default:
			return vm.NewError("invalid type for int()", vm.Value{})
		}
	})
}
