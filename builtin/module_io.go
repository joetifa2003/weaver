package builtin

import (
	"fmt"
	"os"

	"github.com/joetifa2003/weaver/vm"
)

func registerIOModule(builder *RegistryBuilder) {
	builder.RegisterModule("io", map[string]vm.Value{
		"println": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			val := args.Get(0)
			if val.IsError() {
				return val
			}

			fmt.Println(val.String())
			return vm.Value{}
		}),
		"print": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			val := args.Get(0)
			if val.IsError() {
				return val
			}

			fmt.Print(val.String())
			return vm.Value{}
		}),
		"readFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			arg1 := args.Get(0, vm.ValueTypeString)
			if arg1.IsError() {
				return arg1
			}

			filename := arg1.String()
			file, err := os.ReadFile(filename)
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}

			return vm.NewString(string(file))
		}),
	})
}
