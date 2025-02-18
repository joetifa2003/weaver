package builtin

import (
	"fmt"
	"os"

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
		"readFile": vm.NewNativeFunction(func(v *vm.VM, args ...vm.Value) vm.Value {
			if len(args) != 1 {
				panic("readFile() takes exactly 1 argument")
			}

			arg1 := args[0]
			if arg1.VType != vm.ValueTypeString {
				panic("readFile() argument must be a string")
			}

			filename := arg1.String()

			file, err := os.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			return vm.NewString(string(file))
		}),
	})
}
