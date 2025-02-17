package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/vm"
)

func registerStringModule(builder *RegistryBuilder) {
	builder.RegisterModule("strings", map[string]vm.Value{
		"concat": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) vm.Value {
			if len(args) < 2 {
				panic("concat() takes at least 2 arguments")
			}

			var res string
			for _, arg := range args {
				res += arg.String()
			}

			return vm.NewString(res)
		}),
		"split": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) vm.Value {
			str := args[0]
			sep := args[1]

			if str.VType != vm.ValueTypeString {
				panic("split() first argument must be a string")
			}

			if sep.VType != vm.ValueTypeString {
				panic("split() second argument must be a string")
			}

			parts := helpers.SliceMap(strings.Split(str.GetString(), sep.GetString()), func(s string) vm.Value {
				return vm.NewString(s)
			})

			return vm.NewArray(parts)
		}),
	})
}
