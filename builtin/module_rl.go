package builtin

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/joetifa2003/weaver/vm"
)

func registerModuleRL(builder *RegistryBuilder) {
	m := map[string]vm.Value{
		"initWindow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			widthArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			heightArg, err := args.Get(1, vm.ValueTypeFloat, vm.ValueTypeInt)
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
		"setTargetFPS": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			fpsArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			rl.SetTargetFPS(int32(fpsArg.GetInt()))
			return vm.Value{}, nil
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
		"drawRectangle": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			posXArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			posYArg, err := args.Get(1, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			widthArg, err := args.Get(2, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			heightArg, err := args.Get(3, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			colorArg, err := args.Get(4, vm.ValueTypeNativeObject)
			if err != nil {
				return vm.Value{}, err
			}

			posX := posXArg.GetInt()
			posY := posYArg.GetInt()
			width := widthArg.GetInt()
			height := heightArg.GetInt()
			color := colorArg.GetNativeObject().(rl.Color)

			rl.DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), color)
			return vm.Value{}, nil
		}),
		"drawFps": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			posXArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			posYArg, err := args.Get(1, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			rl.DrawFPS(int32(posXArg.GetInt()), int32(posYArg.GetInt()))
			return vm.Value{}, nil
		}),
		"isKeyPressed": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			keyArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			return vm.NewBool(rl.IsKeyPressed(int32(keyArg.GetInt()))), nil
		}),
		"isKeyDown": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			keyArg, err := args.Get(0, vm.ValueTypeFloat, vm.ValueTypeInt)
			if err != nil {
				return vm.Value{}, err
			}
			return vm.NewBool(rl.IsKeyDown(int32(keyArg.GetInt()))), nil
		}),
		"getFrameTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, error) {
			return vm.NewFloat(float64(rl.GetFrameTime())), nil
		}),
	}

	colors := [...]struct {
		name  string
		color rl.Color
	}{
		{"LightGray", rl.LightGray},
		{"Gray", rl.Gray},
		{"DarkGray", rl.DarkGray},
		{"Yellow", rl.Yellow},
		{"Gold", rl.Gold},
		{"Orange", rl.Orange},
		{"Pink", rl.Pink},
		{"Red", rl.Red},
		{"Maroon", rl.Maroon},
		{"Green", rl.Green},
		{"Lime", rl.Lime},
		{"DarkGreen", rl.DarkGreen},
		{"SkyBlue", rl.SkyBlue},
		{"Blue", rl.Blue},
		{"DarkBlue", rl.DarkBlue},
		{"Purple", rl.Purple},
		{"Violet", rl.Violet},
		{"DarkPurple", rl.DarkPurple},
		{"Beige", rl.Beige},
		{"Brown", rl.Brown},
		{"DarkBrown", rl.DarkBrown},
		{"White", rl.White},
		{"Black", rl.Black},
		{"Blank", rl.Blank},
		{"Magenta", rl.Magenta},
		{"RayWhite", rl.RayWhite},
	}
	for _, c := range colors {
		m["color"+c.name] = vm.NewNativeObject(c.color)
	}

	keys := [...]struct {
		name string
		key  int
	}{
		{"A", rl.KeyA},
		{"B", rl.KeyB},
		{"C", rl.KeyC},
		{"D", rl.KeyD},
		{"E", rl.KeyE},
		{"F", rl.KeyF},
		{"G", rl.KeyG},
		{"H", rl.KeyH},
		{"I", rl.KeyI},
		{"J", rl.KeyJ},
		{"K", rl.KeyK},
		{"L", rl.KeyL},
		{"M", rl.KeyM},
		{"N", rl.KeyN},
		{"O", rl.KeyO},
		{"P", rl.KeyP},
		{"Q", rl.KeyQ},
		{"R", rl.KeyR},
		{"S", rl.KeyS},
		{"T", rl.KeyT},
		{"U", rl.KeyU},
		{"V", rl.KeyV},
		{"W", rl.KeyW},
		{"X", rl.KeyX},
		{"Y", rl.KeyY},
		{"Z", rl.KeyZ},
		{"Enter", rl.KeyEnter},
		{"Left", rl.KeyLeft},
		{"Right", rl.KeyRight},
		{"Up", rl.KeyUp},
		{"Down", rl.KeyDown},
	}

	for _, k := range keys {
		m["key"+k.name] = vm.NewInt(k.key)
	}

	builder.RegisterModule("rl", m)
}
