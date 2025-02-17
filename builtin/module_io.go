package builtin

import (
	"fmt"

	"github.com/joetifa2003/weaver/vm"
)

func registerIOModule(builder *RegistryBuilder) {
	builder.RegisterModule("io", map[string]vm.Value{
		"println": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) (res vm.Value) {
			if len(args) != 1 {
				panic("println() takes exactly 1 argument")
			}

			val := args[0]
			fmt.Println(val.String())
			return
		}),
		"print": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) (res vm.Value) {
			if len(args) != 1 {
				panic("println() takes exactly 1 argument")
			}

			val := args[0]
			fmt.Print(val.String())
			return
		}),
	})
}
