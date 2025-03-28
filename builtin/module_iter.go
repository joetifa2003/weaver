package builtin

import "github.com/joetifa2003/weaver/vm"

func registerIterModule(builder *RegistryBuilder) {
	builder.RegisterModule("iter", map[string]vm.Value{
		"each": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			iterArg, ok := args.Get(0, vm.ValueTypeIterator)
			if !ok {
				return iterArg
			}

			fnArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return fnArg
			}

			iter := iterArg.GetIter()

			for val := range iter {
				if val.IsError() {
					return val
				}

				retVal := v.RunFunction(fnArg, val)

				if retVal.IsError() {
					return retVal
				}
			}

			return vm.Value{}
		}),
	})
}
