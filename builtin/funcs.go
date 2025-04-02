package builtin

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *RegistryBuilder) {
	builder.RegisterFunc("error", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		msgArg, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return msgArg
		}
		if len(args) == 1 {
			return vm.NewError(msgArg.GetString(), vm.Value{})
		}

		dataArg, ok := args.Get(1)
		if !ok {
			return dataArg
		}

		return vm.NewError(msgArg.GetString(), dataArg)
	})

	builder.RegisterFunc("isError", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0)
		if !ok {
			return val
		}
		return vm.NewBool(val.IsError())
	})

	builder.RegisterFunc("sleep", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		timeArg, ok := args.Get(0, vm.ValueTypeNumber)
		if !ok {
			return timeArg
		}

		time.Sleep(time.Duration(timeArg.GetNumber()) * time.Millisecond)
		return vm.Value{}
	})

	builder.RegisterFunc("assert", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0)
		if !ok {
			return val
		}

		if !val.IsTruthy() {
			return vm.NewError("assertion failed", vm.Value{})
		}

		return vm.Value{}
	})

	builder.RegisterFunc("echo", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0)
		if !ok {
			return val
		}

		fmt.Println(val.String())
		return vm.Value{}
	})

	builder.RegisterFunc("rand", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		return vm.NewNumber(rand.Float64())
	})

	builder.RegisterFunc("len", func(x *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0, vm.ValueTypeArray, vm.ValueTypeString, vm.ValueTypeObject)
		if !ok {
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
		val, ok := args.Get(0)
		if !ok {
			return val
		}
		return vm.NewString(val.VType.String())
	})

	builder.RegisterFunc("string", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0)
		if !ok {
			return val
		}

		return vm.NewString(val.String())
	})

	builder.RegisterFunc("int", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		val, ok := args.Get(0, vm.ValueTypeNumber)
		if !ok {
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
