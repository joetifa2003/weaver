package ir

import (
	"errors"
	"fmt"
)

type basicBlock struct {
	BlockStmt
	vars   []*basicVar
	parent *basicBlock
	idx    int
}

func (b *basicBlock) allocate(name string) (*basicVar, error) {
	if name != "" {
		for _, v := range b.vars {
			if v.name == name && !v.free {
				return nil, fmt.Errorf("variable %s already defined", name)
			}
		}
	}

	for _, v := range b.vars {
		if v.free {
			v.free = false
			v.name = name
			v.reused = true
			return v, nil
		}
	}

	newVar := &basicVar{
		name:     name,
		idx:      len(b.vars),
		free:     false,
		blockIdx: b.idx,
	}

	b.vars = append(b.vars, newVar)

	return newVar, nil
}

func (b *basicBlock) resolve(name string) (*basicVar, error) {
	for _, v := range b.vars {
		if v.name == name && !v.free {
			return v, nil
		}
	}

	if b.parent != nil {
		return b.parent.resolve(name)
	}

	return nil, errors.New(fmt.Sprintf("cannot find variable %s", name))
}

func (b *basicBlock) deallocateAll() {
	for _, v := range b.vars {
		v.deallocate()
	}
}

type basicVar struct {
	name     string
	idx      int
	blockIdx int
	free     bool
	noInit   bool
	reused   bool
}

func (b *basicVar) id() string {
	return fmt.Sprintf("__$b%dv%d", b.blockIdx, b.idx)
}

func (b *basicVar) deallocate() {
	b.free = true
}
