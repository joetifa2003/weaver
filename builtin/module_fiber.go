package builtin

import (
	"sync"

	"github.com/joetifa2003/weaver/vm"
)

func registerFiberModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("fiber", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"run": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					fnArg, ok := args.Get(0, vm.ValueTypeFunction)
					if !ok {
						return fnArg, false
					}

					return vm.NewTask(runFunc(v, fnArg)), true
				}),

				"wait": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					taskArg, ok := args.Get(0, vm.ValueTypeTask, vm.ValueTypeArray)
					if !ok {
						return taskArg, false
					}

					if taskArg.VType == vm.ValueTypeTask {
						return waitForTask(taskArg)
					} else {
						vals := make([]vm.Value, 0, len(*taskArg.GetArray()))
						for _, task := range *taskArg.GetArray() {
							if err, ok := vm.CheckValueType(task, vm.ValueTypeTask); !ok {
								return err, false
							}
							val, ok := waitForTask(task)
							if !ok {
								return val, false
							}
							vals = append(vals, val)
						}

						return vm.NewArray(vals), true
					}
				}),

				"cancel": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					taskArg, ok := args.Get(0, vm.ValueTypeTask)
					if !ok {
						return taskArg, false
					}

					task := taskArg.GetTask()
					task.Cancel()

					return vm.Value{}, true
				}),

				"newLock": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					l := vm.Value{}
					l.SetLock(&sync.Mutex{})
					return l, true
				}),

				"newChannel": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					var buffer int
					if len(args) > 0 {
						bufferArg, ok := args.Get(0, vm.ValueTypeNumber)
						if !ok {
							return bufferArg, false
						}
						buffer = int(bufferArg.GetNumber())
					}
					val := vm.Value{}
					val.SetChannel(make(chan vm.Value, buffer))
					return val, true
				}),

				"send": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					chArg, ok := args.Get(0, vm.ValueTypeChannel)
					if !ok {
						return chArg, false
					}

					valArg, ok := args.Get(1)
					if !ok {
						return valArg, false
					}

					ch := chArg.GetChannel()
					ch <- valArg

					return valArg, true
				}),

				"close": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					chArg, ok := args.Get(0, vm.ValueTypeChannel)
					if !ok {
						return chArg, false
					}

					ch := chArg.GetChannel()
					close(ch)

					return vm.Value{}, true
				}),

				"recv": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					chArg, ok := args.Get(0, vm.ValueTypeChannel)
					if !ok {
						return chArg, false
					}

					ch := chArg.GetChannel()
					val := <-ch
					return val, true
				}),

				"onRecv": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					chArg, ok := args.Get(0, vm.ValueTypeChannel)
					if !ok {
						return chArg, false
					}

					ch := chArg.GetChannel()
					fnArg, ok := args.Get(1, vm.ValueTypeFunction)
					if !ok {
						return fnArg, false
					}

					for val := range ch {
						v.RunFunction(fnArg, val)
					}

					return vm.Value{}, true
				}),
			},
		)
	})
}

func waitForTask(taskArg vm.Value) (vm.Value, bool) {
	val, ok := taskArg.GetTask().Wait()
	if !ok {
		return val, false
	}

	return val, true
}

func runFunc(v *vm.VM, fnArg vm.Value, args ...vm.Value) *vm.ExecutorTask {
	task := v.Executor.Run(fnArg, args...)
	return task
}
