package vm

import (
	"sync"
	"sync/atomic"
)

type Executor struct {
	Vms []*VMState
	l   sync.RWMutex
	Reg *Registry
}

type VMState struct {
	*VM
	busy atomic.Bool
}

func NewExecutor(reg *Registry) *Executor {
	return &Executor{
		Reg: reg,
	}
}

func (e *Executor) Run(frame Frame, args int) *ExecutorTask {

	e.l.RLock()
	for _, vm := range e.Vms {
		if vm.busy.CompareAndSwap(false, true) {
			vm.Resurrect()
			task := newExecutorTask(vm.VM)

			go func() {
				defer vm.busy.Store(false)
				task.Complete(vm.Run(frame, args))
			}()

			e.l.RUnlock()

			return task
		}
	}
	e.l.RUnlock()

	newVm := &VMState{
		VM:   New(e),
		busy: atomic.Bool{},
	}
	newVm.busy.Store(true)
	e.l.Lock()
	e.Vms = append(e.Vms, newVm)
	e.l.Unlock()

	task := newExecutorTask(newVm.VM)
	go func() {
		defer newVm.busy.Store(false)
		task.Complete(newVm.Run(frame, args))
	}()

	return task
}

type ExecutorTask struct {
	done chan struct{}
	vm   *VM
	once *sync.Once
	val  Value
}

func newExecutorTask(vm *VM) *ExecutorTask {
	return &ExecutorTask{
		done: make(chan struct{}),
		vm:   vm,
		once: &sync.Once{},
	}
}

func (t *ExecutorTask) Wait() Value {
	<-t.done
	return t.val
}

func (t *ExecutorTask) Complete(val Value) {
	t.once.Do(func() {
		t.val = val
		t.vm.running.Store(false)
		close(t.done)
	})
}

func (t *ExecutorTask) Cancel() {
	t.once.Do(func() {
		t.vm.running.Store(false)
		t.vm.ctxCancel()
		t.val = Value{} // return nil value on cancel, TODO: maybe cancellation error?
		close(t.done)
	})
}
