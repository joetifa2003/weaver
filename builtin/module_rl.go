//go:build cgo

package builtin

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/joetifa2003/weaver/vm"
)

func registerModuleRL(builder *vm.RegistryBuilder) {
	m := map[string]vm.Value{
		"initWindow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			widthArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return widthArg
			}
			heightArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return heightArg
			}
			titleArg, ok := args.Get(2, vm.ValueTypeString)
			if !ok {
				return titleArg
			}

			width := widthArg.GetNumber()
			height := heightArg.GetNumber()
			title := titleArg.GetString()

			rl.InitWindow(int32(width), int32(height), title)
			return vm.Value{}
		}),
		"windowShouldClose": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			return vm.NewBool(rl.WindowShouldClose())
		}),
		"setTargetFPS": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			fpsArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return fpsArg
			}
			rl.SetTargetFPS(int32(fpsArg.GetNumber()))
			return vm.Value{}
		}),
		"closeWindow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			rl.CloseWindow()
			return vm.Value{}
		}),
		"beginDrawing": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			rl.BeginDrawing()
			return vm.Value{}
		}),
		"endDrawing": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			rl.EndDrawing()
			return vm.Value{}
		}),
		"clearBackground": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			colorArg, ok := args.Get(0, vm.ValueTypeNativeObject)
			if !ok {
				return colorArg
			}
			color := colorArg.GetNativeObject().Obj.(rl.Color)
			rl.ClearBackground(color)
			return vm.Value{}
		}),
		"drawRectangle": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			posXArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return posXArg
			}
			posYArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return posYArg
			}
			widthArg, ok := args.Get(2, vm.ValueTypeNumber)
			if !ok {
				return widthArg
			}
			heightArg, ok := args.Get(3, vm.ValueTypeNumber)
			if !ok {
				return heightArg
			}
			colorArg, ok := args.Get(4, vm.ValueTypeNativeObject)
			if !ok {
				return colorArg
			}

			posX := posXArg.GetNumber()
			posY := posYArg.GetNumber()
			width := widthArg.GetNumber()
			height := heightArg.GetNumber()
			color := colorArg.GetNativeObject().Obj.(rl.Color)

			rl.DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), color)
			return vm.Value{}
		}),
		"drawFps": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			posXArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return posXArg
			}
			posYArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return posYArg
			}
			rl.DrawFPS(int32(posXArg.GetNumber()), int32(posYArg.GetNumber()))
			return vm.Value{}
		}),
		"isKeyPressed": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return keyArg
			}
			return vm.NewBool(rl.IsKeyPressed(int32(keyArg.GetNumber())))
		}),
		"isKeyDown": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			keyArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return keyArg
			}
			return vm.NewBool(rl.IsKeyDown(int32(keyArg.GetNumber())))
		}),
		"getFrameTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			return vm.NewNumber(float64(rl.GetFrameTime()))
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
		m["color"+c.name] = vm.NewNativeObject(c.color, nil)
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
		m["key"+k.name] = vm.NewNumber(float64(k.key))
	}

	builder.RegisterModule("rl", vm.NewObject(m))
}
