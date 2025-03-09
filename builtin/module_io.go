package builtin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joetifa2003/weaver/vm"
)

func registerIOModule(builder *RegistryBuilder) {
	builder.RegisterModule("io", map[string]vm.Value{
		"println": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			val, ok := args.Get(0)
			if !ok {
				return val
			}

			fmt.Println(val.String())
			return vm.Value{}
		}),
		"print": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			val, ok := args.Get(0)
			if !ok {
				return val
			}

			fmt.Print(val.String())
			return vm.Value{}
		}),
		"readFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			arg1, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return arg1
			}

			filename := arg1.String()
			file, err := os.ReadFile(filename)
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}

			return vm.NewString(string(file))
		}),
		"writeFile": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}
			contentArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return contentArg
			}

			err := os.WriteFile(pathArg.GetString(), []byte(contentArg.GetString()), 0644)
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.Value{}
		}),
		"exists": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			_, err := os.Stat(pathArg.GetString())
			return vm.NewBool(!os.IsNotExist(err))
		}),
		"mkdir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			err := os.MkdirAll(pathArg.GetString(), 0755)
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.Value{}
		}),
		"remove": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			err := os.RemoveAll(pathArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.Value{}
		}),
		"rename": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			oldArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return oldArg
			}
			newArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return newArg
			}

			err := os.Rename(oldArg.GetString(), newArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.Value{}
		}),
		"join": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			paths := make([]string, len(args))
			for i, arg := range args {
				if arg.IsError() {
					return arg
				}
				paths[i] = arg.GetString()
			}
			return vm.NewString(filepath.Join(paths...))
		}),
		"dirname": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}
			return vm.NewString(filepath.Dir(pathArg.GetString()))
		}),
		"basename": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}
			return vm.NewString(filepath.Base(pathArg.GetString()))
		}),
		"extname": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}
			return vm.NewString(filepath.Ext(pathArg.GetString()))
		}),
		"readDir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			entries, err := os.ReadDir(pathArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}

			results := make([]vm.Value, len(entries))
			for i, entry := range entries {
				results[i] = vm.NewString(entry.Name())
			}
			return vm.NewArray(results)
		}),
		"size": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			info, err := os.Stat(pathArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.NewNumber(float64(info.Size()))
		}),
		"isDir": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			info, err := os.Stat(pathArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.NewBool(info.IsDir())
		}),
		"modTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			pathArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return pathArg
			}

			info, err := os.Stat(pathArg.GetString())
			if err != nil {
				return vm.NewError(err.Error(), vm.Value{})
			}
			return vm.NewNumber(float64(info.ModTime().Unix()))
		}),
	})
}
