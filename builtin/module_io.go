package builtin

import (
	"fmt"
	"os"

	"github.com/joetifa2003/weaver/vm"
)

func registerIOModule(builder *RegistryBuilder) {
	builder.RegisterModule("io", map[string]vm.Value{
		"println": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			val, err := args.Get(0, vm.ValueTypeAny)
			if err != nil {
				return vm.Value{}, err
			}

			fmt.Println(val.String())
			return vm.Value{}, nil
		}),
		"print": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			val, err := args.Get(0, vm.ValueTypeAny)
			if err != nil {
				return vm.Value{}, err
			}

			fmt.Print(val.String())
			return vm.Value{}, nil
		}),
		"readFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			arg1, err := args.Get(0, vm.ValueTypeString)
			if err != nil {
				return vm.Value{}, err
			}

			filename := arg1.String()

			file, err := os.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			return vm.NewString(string(file)), nil
		}),
	})
}
