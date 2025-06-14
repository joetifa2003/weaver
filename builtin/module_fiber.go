package builtin

import (
	"sync"

	"github.com/joetifa2003/weaver/vm"
)

func registerFiberModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("fiber", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"run": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					fnArg, ok := args.Get(0, vm.ValueTypeFunction)
					if !ok {
						return fnArg
					}

					return vm.NewTask(runFunc(v, fnArg))
				}),

				"wait": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					taskArg, ok := args.Get(0, vm.ValueTypeTask, vm.ValueTypeArray)
					if !ok {
						return taskArg
					}

					if taskArg.VType == vm.ValueTypeTask {
						return waitForTask(taskArg)
					} else {
						vals := make([]vm.Value, 0, len(*taskArg.GetArray()))
						for _, task := range *taskArg.GetArray() {
							if err, ok := vm.CheckValueType(task, vm.ValueTypeTask); !ok {
								return err
							}
							vals = append(vals, waitForTask(task))
						}

						return vm.NewArray(vals)
					}
				}),

				"cancel": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					taskArg, ok := args.Get(0, vm.ValueTypeTask)
					if !ok {
						return taskArg
					}

					task := taskArg.GetTask()
					task.Cancel()

					return vm.Value{}
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
			},
		)
	})
}

func waitForTask(taskArg vm.Value) vm.Value {
	return taskArg.GetTask().Wait()
}

func runFunc(v *vm.VM, fnArg vm.Value, args ...vm.Value) *vm.ExecutorTask {
	fn := fnArg.GetFunction()

	task := v.Executor.Run(vm.Frame{
		Instructions: fn.Instructions,
		NumVars:      fn.NumVars,
		FreeVars:     fn.FreeVars,
		Constants:    fn.Constants,
		Path:         fn.Path,
		HaltAfter:    true,
	}, args...)

	return task
}
