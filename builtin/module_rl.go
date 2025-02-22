package builtin

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/joetifa2003/weaver/vm"
)

func registerModuleRL(builder *RegistryBuilder) {
	builder.RegisterModule("rl", map[string]vm.Value{
		"initWindow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			widthArg, err := args.Get(0, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			heightArg, err := args.Get(1, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			titleArg, err := args.Get(2, vm.ValueTypeString)
			if err != nil {
				return vm.Value{}, err
			}

			width := widthArg.GetInt()
			height := heightArg.GetInt()
			title := titleArg.GetString()

			rl.InitWindow(int32(width), int32(height), title)
			return vm.Value{}, nil
		}),
		"windowShouldClose": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			return vm.NewBool(rl.WindowShouldClose()), nil
		}),
		"closeWindow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			rl.CloseWindow()
			return vm.Value{}, nil
		}),
		"beginDrawing": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			rl.BeginDrawing()
			return vm.Value{}, nil
		}),
		"endDrawing": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			rl.EndDrawing()
			return vm.Value{}, nil
		}),
		"clearBackground": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			colorArg, err := args.Get(0, vm.ValueTypeNativeObject)
			if err != nil {
				return vm.Value{}, err
			}
			color := colorArg.GetNativeObject().(rl.Color)
			rl.ClearBackground(color)
			return vm.Value{}, nil
		}),
		"colorRayWhite": vm.NewNativeObject(rl.RayWhite),
	})
}
