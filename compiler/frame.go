package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/opcode"
)

type Frame struct {
	Parent       *Frame
	Vars         []*Var
	FreeVars     []*Var
	Instructions []opcode.OpCode
	Blocks       *ds.Stack[*Block]
}

func NewFrame(parent *Frame) *Frame {
	return &Frame{
		Parent:       parent,
		Vars:         []*Var{},
		Instructions: []opcode.OpCode{},
		Blocks:       ds.NewStack(&Block{}),
	}
}

type VarScope int

const (
	VarScopeGlobal VarScope = iota
	VarScopeLocal
	VarScopeFree
)

type Var struct {
	Scope VarScope
	Name  string
	Index int
}

type Block struct {
	Vars []*Var
}

func (c *Frame) addInstructions(instructions []opcode.OpCode) {
	c.Instructions = append(c.Instructions, instructions...)
}

func (c *Frame) defineVar(name string) int {
	v := &Var{Name: name, Index: len(c.Vars)}
	c.Vars = append(c.Vars, v)
	c.Blocks.Peek().Vars = append(c.Blocks.Peek().Vars, v)
	return len(c.Vars) - 1
}

func (c *Frame) defineFreeVar(name string, parentIdx int) *Var {
	v := &Var{Name: name, Index: len(c.FreeVars)}
	c.FreeVars = append(c.FreeVars, v)
	return v
}

func (c *Frame) resolve(name string) (*Var, error) {
	for _, b := range c.Blocks.Iter() {
		for _, v := range b.Vars {
			if v.Name == name {
				return v, nil
			}
		}
	}

	for _, v := range c.FreeVars {
		if v.Name == name {
			return v, nil
		}
	}

	if c.Parent != nil {
		idx, err := c.Parent.resolve(name)
		if err == nil {
			return nil, err
		}

		return c.defineFreeVar(name), nil
	}

	return nil, fmt.Errorf("%w: %s", ErrVarNotFound, name)
}

func (c *Frame) beginBlock() {
	b := &Block{}
	c.Blocks.Push(b)
}

func (c *Frame) endBlock() {
	c.Blocks.Pop()
}
