package builtin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/vm"
)

func registerIOModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("io", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"println": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					val, ok := args.Get(0)
					if !ok {
						return val, false
					}

					fmt.Println(val.String())
					return vm.Value{}, true
				}),
				"print": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					val, ok := args.Get(0)
					if !ok {
						return val, false
					}

					fmt.Print(val.String())
					return vm.Value{}, true
				}),
				"readFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					arg1, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return arg1, false
					}

					filename := arg1.String()
					file, err := os.ReadFile(filename)
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}

					return vm.NewString(string(file)), true
				}),
				"writeFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}
					contentArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return contentArg, false
					}

					err := os.WriteFile(pathArg.GetString(), []byte(contentArg.GetString()), 0644)
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.Value{}, true
				}),
				"exists": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					_, err := os.Stat(pathArg.GetString())
					return vm.NewBool(!os.IsNotExist(err)), true
				}),
				"mkdir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					err := os.MkdirAll(pathArg.GetString(), 0755)
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.Value{}, true
				}),
				"remove": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					err := os.RemoveAll(pathArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.Value{}, true
				}),
				"rename": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					oldArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return oldArg, false
					}
					newArg, ok := args.Get(1, vm.ValueTypeString)
					if !ok {
						return newArg, false
					}

					err := os.Rename(oldArg.GetString(), newArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.Value{}, true
				}),
				"join": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					paths := make([]string, len(args))
					for i, arg := range args {
						if arg.IsError() {
							return arg, false
						}
						paths[i] = arg.GetString()
					}
					return vm.NewString(filepath.Join(paths...)), true
				}),
				"dirname": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}
					return vm.NewString(filepath.Dir(pathArg.GetString())), true
				}),
				"basename": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}
					return vm.NewString(filepath.Base(pathArg.GetString())), true
				}),
				"extname": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}
					return vm.NewString(filepath.Ext(pathArg.GetString())), true
				}),
				"readDir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					entries, err := os.ReadDir(pathArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), true
					}

					results := make([]vm.Value, len(entries))
					for i, entry := range entries {
						results[i] = vm.NewString(entry.Name())
					}
					return vm.NewArray(results), true
				}),
				"size": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					info, err := os.Stat(pathArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.NewNumber(float64(info.Size())), true
				}),
				"isDir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					info, err := os.Stat(pathArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.NewBool(info.IsDir()), true
				}),
				"modTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					pathArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return pathArg, false
					}

					info, err := os.Stat(pathArg.GetString())
					if err != nil {
						return vm.NewError(err.Error(), vm.Value{}), false
					}
					return vm.NewNumber(float64(info.ModTime().Unix())), true
				}),
				"exec": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					cmdArg, ok := args.Get(0, vm.ValueTypeString)
					if !ok {
						return cmdArg, false
					}

					argsArg, ok := args.Get(1, vm.ValueTypeArray)
					if !ok {
						return argsArg, false
					}

					cmdArgs := helpers.SliceMap(*argsArg.GetArray(), func(v vm.Value) string {
						return v.String()
					})

					cmd := exec.Command(cmdArg.GetString(), cmdArgs...)
					output, err := cmd.CombinedOutput()
					if err != nil {
						return vm.NewError(err.Error(), vm.NewString(string(output))), false
					}

					return vm.NewString(string(output)), true
				}),
			},
		)
	})
}
