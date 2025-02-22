package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/vm"
)

func registerStringModule(builder *RegistryBuilder) {
	builder.RegisterModule("strings", map[string]vm.Value{
		"concat": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {

			var res string
			for _, arg := range args {
				res += arg.String()
			}

			return vm.NewString(res), nil
		}),
		"split": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			strArg, err := args.Get(0, vm.ValueTypeString)
			if err != nil {
				return vm.Value{}, nil
			}

			sepArg, err := args.Get(1, vm.ValueTypeString)
			if err != nil {
				return vm.Value{}, nil
			}

			str := strArg.GetString()
			sep := sepArg.GetString()

			parts := helpers.SliceMap(strings.Split(str, sep), func(s string) vm.Value {
				return vm.NewString(s)
			})

			return vm.NewArray(parts), nil
		}),
	})
}
