package builtin

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncs(builder *vm.RegistryBuilder) {
	builder.RegisterFunc("error", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		msgArg, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return msgArg, false
		}
		if len(args) == 1 {
			return vm.NewError(msgArg.GetString(), vm.Value{}), true
		}

		dataArg, ok := args.Get(1)
		if !ok {
			return dataArg, false
		}

		return vm.NewError(msgArg.GetString(), dataArg), true
	})

	builder.RegisterFunc("isError", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0)
		if !ok {
			return val, false
		}

		return vm.NewBool(val.IsError()), true
	})

	builder.RegisterFunc("sleep", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		timeArg, ok := args.Get(0, vm.ValueTypeNumber)
		if !ok {
			return timeArg, false
		}

		time.Sleep(time.Duration(timeArg.GetNumber()) * time.Millisecond)
		return vm.Value{}, true
	})

	builder.RegisterFunc("assert", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0)
		if !ok {
			return val, false
		}

		if !val.IsTruthy() {
			return vm.NewError("assertion failed", vm.Value{}), false
		}

		return vm.Value{}, true
	})

	builder.RegisterFunc("echo", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0)
		if !ok {
			return val, false
		}

		fmt.Println(val.String())
		return vm.Value{}, true
	})

	builder.RegisterFunc("rand", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		return vm.NewNumber(rand.Float64()), true
	})

	builder.RegisterFunc("len", func(x *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0, vm.ValueTypeArray, vm.ValueTypeString, vm.ValueTypeObject)
		if !ok {
			return val, false
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
			return vm.NewError("invalid type for len()", vm.Value{}), false
		}

		return res, true
	})

	builder.RegisterFunc("type", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0)
		if !ok {
			return val, false
		}
		return vm.NewString(val.VType.String()), true
	})

	builder.RegisterFunc("string", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0)
		if !ok {
			return val, false
		}

		return vm.NewString(val.String()), true
	})

	builder.RegisterFunc("number", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return val, false
		}

		floatVal, err := strconv.ParseFloat(val.GetString(), 64)
		if err != nil {
			return vm.NewErrFromErr(err), false
		}

		return vm.NewNumber(floatVal), true
	})

	builder.RegisterFunc("int", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		val, ok := args.Get(0, vm.ValueTypeNumber)
		if !ok {
			return val, false
		}

		switch val.VType {
		case vm.ValueTypeNumber:
			val.SetNumber(float64(int(val.GetNumber())))
			return val, true
		default:
			return vm.NewError("invalid type for int()", vm.Value{}), false
		}
	})
}
