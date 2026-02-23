package builtin

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"

	"github.com/joetifa2003/weaver/vm"
)

type pluginCallContextKey struct{}

type pluginCallContext struct {
	args    vm.NativeFunctionArgs
	ret     vm.Value
	handles []vm.Value
}

func (p *pluginCallContext) addHandle(v vm.Value) int32 {
	p.handles = append(p.handles, v)
	return int32(len(p.handles) - 1)
}

func (p *pluginCallContext) getHandle(id int32) *vm.Value {
	if id < 0 || int(id) >= len(p.handles) {
		return nil
	}
	return &p.handles[id]
}

func registerPluginModule(r *vm.RegistryBuilder) {
	r.RegisterModule("plugin", func() vm.Value {
		return vm.NewObject(map[string]vm.Value{
			"load": vm.NewNativeFunction("plugin.load", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
				pathArg, ok := args.Get(0, vm.ValueTypeString)
				if !ok {
					return pathArg, false
				}
				path := pathArg.GetString()

				wasmBytes, err := os.ReadFile(path)
				if err != nil {
					return vm.NewError("failed to read wasm plugin", vm.NewString(err.Error())), false
				}

				ctx := context.Background()
				runtime := wazero.NewRuntime(ctx)

				wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

				_, err = runtime.NewHostModuleBuilder("weaver").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						return int32(len(pCtx.args.Args))
					}).Export("get_arg_len").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, idx int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if idx < 0 || int(idx) >= len(pCtx.args.Args) {
							return -1 // Invalid
						}
						// Map vm.ValueType to an int32 for the SDK
						return int32(pCtx.args.Args[idx].VType)
					}).Export("get_arg_type").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, idx int32) float64 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if idx < 0 || int(idx) >= len(pCtx.args.Args) {
							return 0
						}
						return pCtx.args.Args[idx].GetNumber()
					}).Export("get_arg_num").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, idx int32, ptr int32, maxLen int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if idx < 0 || int(idx) >= len(pCtx.args.Args) {
							return -1
						}
						str := pCtx.args.Args[idx].GetString()
						bytes := []byte(str)

						copyLen := int32(len(bytes))
						if copyLen > maxLen {
							copyLen = maxLen
						}

						m.Memory().Write(uint32(ptr), bytes[:copyLen])
						return int32(len(bytes)) // Return the full length so caller knows if it was truncated
					}).Export("get_arg_str").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, val float64) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						pCtx.ret = vm.NewNumber(val)
					}).Export("return_num").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, ptr int32, len int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if bytes, ok := m.Memory().Read(uint32(ptr), uint32(len)); ok {
							pCtx.ret = vm.NewString(string(bytes))
						} else {
							pCtx.ret = vm.NewError("plugin memory read out of bounds", vm.Value{})
						}
					}).Export("return_str").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, val int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						pCtx.ret = vm.NewBool(val != 0)
					}).Export("return_bool").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						pCtx.ret = vm.Value{}
						pCtx.ret.SetNil()
					}).Export("return_nil").
					// === Handle-Based Builder API ===
					// Value creation
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := vm.Value{}
						v.SetNil()
						return pCtx.addHandle(v)
					}).Export("value_new_nil").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, val int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						return pCtx.addHandle(vm.NewBool(val != 0))
					}).Export("value_new_bool").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, val float64) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						return pCtx.addHandle(vm.NewNumber(val))
					}).Export("value_new_number").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, ptr int32, length int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if bytes, ok := m.Memory().Read(uint32(ptr), uint32(length)); ok {
							return pCtx.addHandle(vm.NewString(string(bytes)))
						}
						return -1
					}).Export("value_new_string").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						return pCtx.addHandle(vm.NewObject(map[string]vm.Value{}))
					}).Export("value_new_object").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						return pCtx.addHandle(vm.NewArray([]vm.Value{}))
					}).Export("value_new_array").
					// Value mutation
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, objHandle int32, keyPtr int32, keyLen int32, valHandle int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						obj := pCtx.getHandle(objHandle)
						val := pCtx.getHandle(valHandle)
						if obj == nil || val == nil {
							return
						}
						keyBytes, ok := m.Memory().Read(uint32(keyPtr), uint32(keyLen))
						if !ok {
							return
						}
						obj.GetObject()[string(keyBytes)] = *val
					}).Export("value_object_set").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, arrHandle int32, valHandle int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						arr := pCtx.getHandle(arrHandle)
						val := pCtx.getHandle(valHandle)
						if arr == nil || val == nil {
							return
						}
						arrSlice := arr.GetArray()
						*arrSlice = append(*arrSlice, *val)
					}).Export("value_array_push").
					// Value reading (for args)
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, idx int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if idx < 0 || int(idx) >= len(pCtx.args.Args) {
							return -1
						}
						return pCtx.addHandle(pCtx.args.Args[idx])
					}).Export("get_arg_value").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, handle int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := pCtx.getHandle(handle)
						if v == nil {
							return -1
						}
						return int32(v.VType)
					}).Export("value_type").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, handle int32) float64 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := pCtx.getHandle(handle)
						if v == nil {
							return 0
						}
						return v.GetNumber()
					}).Export("value_get_number").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, handle int32, ptr int32, maxLen int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := pCtx.getHandle(handle)
						if v == nil {
							return -1
						}
						str := v.GetString()
						bytes := []byte(str)
						if ptr == 0 && maxLen == 0 {
							return int32(len(bytes))
						}
						copyLen := int32(len(bytes))
						if copyLen > maxLen {
							copyLen = maxLen
						}
						m.Memory().Write(uint32(ptr), bytes[:copyLen])
						return int32(len(bytes))
					}).Export("value_get_string").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, handle int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := pCtx.getHandle(handle)
						if v == nil {
							return 0
						}
						if v.GetBool() {
							return 1
						}
						return 0
					}).Export("value_get_bool").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, objHandle int32, keyPtr int32, keyLen int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						obj := pCtx.getHandle(objHandle)
						if obj == nil {
							return -1
						}
						keyBytes, ok := m.Memory().Read(uint32(keyPtr), uint32(keyLen))
						if !ok {
							return -1
						}
						val, exists := obj.GetObject()[string(keyBytes)]
						if !exists {
							return -1
						}
						return pCtx.addHandle(val)
					}).Export("value_object_get").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, arrHandle int32, idx int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						arr := pCtx.getHandle(arrHandle)
						if arr == nil {
							return -1
						}
						arrSlice := arr.GetArray()
						if idx < 0 || int(idx) >= len(*arrSlice) {
							return -1
						}
						return pCtx.addHandle((*arrSlice)[idx])
					}).Export("value_array_get").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, arrHandle int32) int32 {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						arr := pCtx.getHandle(arrHandle)
						if arr == nil {
							return -1
						}
						return int32(len(*arr.GetArray()))
					}).Export("value_array_len").
					// Return via handle
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, handle int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						v := pCtx.getHandle(handle)
						if v != nil {
							pCtx.ret = *v
						}
					}).Export("return_value").
					NewFunctionBuilder().
					WithFunc(func(ctx context.Context, m api.Module, ptr int32, length int32) {
						pCtx := ctx.Value(pluginCallContextKey{}).(*pluginCallContext)
						if bytes, ok := m.Memory().Read(uint32(ptr), uint32(length)); ok {
							pCtx.ret = vm.NewError(string(bytes), vm.Value{})
						}
					}).Export("return_error").
					Instantiate(ctx)

				if err != nil {
					return vm.NewError("failed to instantiate host module", vm.NewString(err.Error())), false
				}

				compiledModule, err := runtime.CompileModule(ctx, wasmBytes)
				if err != nil {
					return vm.NewError("failed to compile wasm module", vm.NewString(err.Error())), false
				}

				module, err := runtime.InstantiateModule(ctx, compiledModule, wazero.NewModuleConfig())
				if err != nil {
					return vm.NewError("failed to instantiate wasm module", vm.NewString(err.Error())), false
				}

				if initFn := module.ExportedFunction("_initialize"); initFn != nil {
					if _, err := initFn.Call(ctx); err != nil {
						return vm.NewError("failed to initialize wasm module", vm.NewString(err.Error())), false
					}
				}

				exportedFuncs := map[string]vm.Value{}

				// Expose a close function to free resources
				exportedFuncs["close"] = vm.NewNativeFunction("close", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					runtime.Close(context.Background())
					return vm.Value{}, true
				})

				for name, def := range module.ExportedFunctionDefinitions() {
					// Skip internal functions (often starting with _, memory, or WASI internals)
					if name == "_start" || name == "memory" || name == "malloc" || name == "free" || def.Name() == "" {
						continue
					}

					funcName := name
					fn := module.ExportedFunction(funcName)

					// Only expose functions that take no arguments and return no arguments via Wasm ABI.
					// We pass arguments via the host context.
					if len(def.ParamTypes()) == 0 && len(def.ResultTypes()) == 0 {
						exportedFuncs[funcName] = vm.NewNativeFunction(funcName, func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
							callCtx := &pluginCallContext{
								args: args,
								ret:  vm.Value{},
							}
							callCtx.ret.SetNil()

							ctx := context.WithValue(context.Background(), pluginCallContextKey{}, callCtx)
							_, err := fn.Call(ctx)
							if err != nil {
								return vm.NewError(fmt.Sprintf("plugin function %s panic", funcName), vm.NewString(err.Error())), false
							}
							return callCtx.ret, true
						})
					}
				}

				return vm.NewObject(exportedFuncs), true
			}),
		})
	})
}
