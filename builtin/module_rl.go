//go:build cgo

package builtin

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/joetifa2003/weaver/vm"
)

func registerModuleRL(builder *vm.RegistryBuilder) {

	builder.RegisterModule("rl", func() vm.Value {
		m := map[string]vm.Value{
			"initWindow": vm.NewNativeFunction("initWindow", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				widthArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return widthArg, false
				}
				heightArg, ok := args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return heightArg, false
				}
				titleArg, ok := args.Get(2, vm.ValueTypeString)
				if !ok {
					return titleArg, false
				}

				width := widthArg.GetNumber()
				height := heightArg.GetNumber()
				title := titleArg.GetString()

				rl.InitWindow(int32(width), int32(height), title)
				return vm.Value{}, true
			}),
			"windowShouldClose": vm.NewNativeFunction("windowShouldClose", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				return vm.NewBool(rl.WindowShouldClose()), false
			}),
			"setTargetFPS": vm.NewNativeFunction("setTargetFPS", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				fpsArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return fpsArg, false
				}
				rl.SetTargetFPS(int32(fpsArg.GetNumber()))
				return vm.Value{}, true
			}),
			"closeWindow": vm.NewNativeFunction("closeWindow", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				rl.CloseWindow()
				return vm.Value{}, true
			}),
			"beginDrawing": vm.NewNativeFunction("beginDrawing", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				rl.BeginDrawing()
				return vm.Value{}, true
			}),
			"endDrawing": vm.NewNativeFunction("endDrawing", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				rl.EndDrawing()
				return vm.Value{}, true
			}),
			"clearBackground": vm.NewNativeFunction("clearBackground", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				colorArg, ok := args.Get(0, vm.ValueTypeNativeObject)
				if !ok {
					return colorArg, false
				}
				color := colorArg.GetNativeObject().Obj.(rl.Color)
				rl.ClearBackground(color)
				return vm.Value{}, true
			}),
			"drawRectangle": vm.NewNativeFunction("drawRectangle", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				posXArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return posXArg, false
				}
				posYArg, ok := args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return posYArg, false
				}
				widthArg, ok := args.Get(2, vm.ValueTypeNumber)
				if !ok {
					return widthArg, false
				}
				heightArg, ok := args.Get(3, vm.ValueTypeNumber)
				if !ok {
					return heightArg, false
				}
				colorArg, ok := args.Get(4, vm.ValueTypeNativeObject)
				if !ok {
					return colorArg, false
				}

				posX := posXArg.GetNumber()
				posY := posYArg.GetNumber()
				width := widthArg.GetNumber()
				height := heightArg.GetNumber()
				color := colorArg.GetNativeObject().Obj.(rl.Color)

				rl.DrawRectangle(int32(posX), int32(posY), int32(width), int32(height), color)
				return vm.Value{}, true
			}),
			"drawFps": vm.NewNativeFunction("drawFps", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				posXArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return posXArg, false
				}
				posYArg, ok := args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return posYArg, false
				}
				rl.DrawFPS(int32(posXArg.GetNumber()), int32(posYArg.GetNumber()))
				return vm.Value{}, true
			}),
			"isKeyPressed": vm.NewNativeFunction("isKeyPressed", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				keyArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return keyArg, false
				}
				return vm.NewBool(rl.IsKeyPressed(int32(keyArg.GetNumber()))), true
			}),
			"isKeyDown": vm.NewNativeFunction("isKeyDown", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				keyArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return keyArg, false
				}
				return vm.NewBool(rl.IsKeyDown(int32(keyArg.GetNumber()))), true
			}),
			"getFrameTime": vm.NewNativeFunction("getFrameTime", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				return vm.NewNumber(float64(rl.GetFrameTime())), true
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
		return vm.NewObject(m)
	})
}
