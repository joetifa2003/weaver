package vm

import (
	"sync"

	"github.com/joetifa2003/weaver/internal/pkg/pool"
)

type Executor struct {
	l    sync.RWMutex
	Reg  *Registry
	Pool *pool.Pool[*VM]
}

func NewExecutor(reg *Registry) *Executor {
	e := &Executor{
		Reg: reg,
	}
	e.Pool = pool.New(func() *VM {
		return New(e)
	})

	return e
}

func (e *Executor) Run(function Value, args ...Value) *ExecutorTask {
	v := e.Pool.Get()
	v.Resurrect()

	task := newExecutorTask(v)
	go func() {
		defer e.Pool.Put(v)
		task.Complete(v.RunFunction(function, args...))
	}()

	return task
}

type ExecutorTask struct {
	done chan struct{}
	vm   *VM
	once *sync.Once
	val  Value
	ok   bool
}

func newExecutorTask(vm *VM) *ExecutorTask {
	return &ExecutorTask{
		done: make(chan struct{}),
		vm:   vm,
		once: &sync.Once{},
	}
}

func (t *ExecutorTask) Wait() (Value, bool) {
	<-t.done
	return t.val, t.ok
}

func (t *ExecutorTask) Complete(val Value, ok bool) {
	t.once.Do(func() {
		t.val = val
		t.ok = ok
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
