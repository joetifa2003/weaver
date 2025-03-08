package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/vm"
)

func registerStringModule(builder *RegistryBuilder) {
	builder.RegisterModule("strings", map[string]vm.Value{
		"concat": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			var res string
			for _, arg := range args {
				if arg.IsError() {
					return arg
				}
				res += arg.String()
			}

			return vm.NewString(res)
		}),
		"split": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}

			sepArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return sepArg
			}

			str := strArg.GetString()
			sep := sepArg.GetString()

			parts := helpers.SliceMap(strings.Split(str, sep), func(s string) vm.Value {
				return vm.NewString(s)
			})

			return vm.NewArray(parts)
		}),
	})
}
