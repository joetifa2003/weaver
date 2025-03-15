package builtin

import (
	"sync"

	"github.com/joetifa2003/weaver/vm"
)

func registerFiberModule(builder *RegistryBuilder) {
	builder.RegisterModule("fiber", map[string]vm.Value{

		"run": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			fnArg, ok := args.Get(0, vm.ValueTypeFunction)
			if !ok {
				return fnArg
			}

			fn := fnArg.GetFunction()

			res := make(chan vm.Value)
			go func() {
				v := v.Executor.Run(&vm.Frame{
					Instructions: fn.Instructions,
					NumVars:      fn.NumVars,
					FreeVars:     fn.FreeVars,
					HaltAfter:    true,
				}, 0)
				res <- v
				close(res)
			}()

			resVal := vm.Value{}
			resVal.SetTask(res)

			return resVal
		}),

		"wait": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			taskArg, ok := args.Get(0, vm.ValueTypeTask)
			if !ok {
				return taskArg
			}

			task := taskArg.GetTask()
			task.L.Lock()
			defer task.L.Unlock()

			if task.Done {
				return task.Value
			}

			val := <-task.C
			task.Done = true
			task.Value = val
			return val
		}),

		"newLock": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			l := vm.Value{}
			l.SetLock(&sync.Mutex{})
			return l
		}),

		"newChannel": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			var buffer int
			if len(args) > 0 {
				bufferArg, ok := args.Get(0, vm.ValueTypeNumber)
				if !ok {
					return bufferArg
				}
				buffer = int(bufferArg.GetNumber())
			}
			val := vm.Value{}
			val.SetChannel(make(chan vm.Value, buffer))
			return val
		}),

		"send": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			chArg, ok := args.Get(0, vm.ValueTypeChannel)
			if !ok {
				return chArg
			}

			valArg, ok := args.Get(1)
			if !ok {
				return valArg
			}

			ch := chArg.GetChannel()
			ch <- valArg

			return valArg
		}),

		"close": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			chArg, ok := args.Get(0, vm.ValueTypeChannel)
			if !ok {
				return chArg
			}

			ch := chArg.GetChannel()
			close(ch)

			return vm.Value{}
		}),

		"recv": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			chArg, ok := args.Get(0, vm.ValueTypeChannel)
			if !ok {
				return chArg
			}

			ch := chArg.GetChannel()
			val := <-ch
			return val
		}),

		"onRecv": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			chArg, ok := args.Get(0, vm.ValueTypeChannel)
			if !ok {
				return chArg
			}

			ch := chArg.GetChannel()
			fnArg, ok := args.Get(1, vm.ValueTypeFunction)
			if !ok {
				return fnArg
			}

			for val := range ch {
				v.RunFunction(fnArg, val)
			}

			return vm.Value{}
		}),
	})
}
