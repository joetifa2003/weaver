package vm

import (
	"sync"
	"sync/atomic"
)

type Executor struct {
	Vms       []*VMState
	Constants []Value
	l         sync.RWMutex
}

type VMState struct {
	*VM
	busy atomic.Bool
}

func NewExecutor(constants []Value) *Executor {
	return &Executor{
		Constants: constants,
	}
}

func (e *Executor) Run(frame Frame, args int) Value {
	e.l.RLock()
	for _, vm := range e.Vms {
		if vm.busy.CompareAndSwap(false, true) {
			defer vm.busy.Store(false)
			val := vm.Run(frame, args)
			return val
		}
	}
	e.l.RUnlock()

	newVm := &VMState{
		VM:   New(e.Constants, e),
		busy: atomic.Bool{},
	}
	newVm.busy.Store(true)
	e.l.Lock()
	e.Vms = append(e.Vms, newVm)
	e.l.Unlock()

	val := newVm.Run(frame, args)

	newVm.busy.Store(false)

	return val
}
