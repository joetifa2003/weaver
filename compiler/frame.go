package compiler

import (
	"fmt"

	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/opcode"
)

type Frame struct {
	Vars         []*Var
	Instructions []opcode.OpCode
	Blocks       *ds.Stack[*Block]
}

func NewFrame() *Frame {
	return &Frame{
		Vars:         []*Var{},
		Instructions: []opcode.OpCode{},
		Blocks:       ds.NewStack(&Block{}),
	}
}

type Var struct {
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

func (c *Frame) resolve(name string) (int, error) {
	for _, b := range c.Blocks.Iter() {
		for _, v := range b.Vars {
			if v.Name == name {
				return v.Index, nil
			}
		}
	}

	return -1, fmt.Errorf("%w: %s", ErrVarNotFound, name)
}

func (c *Frame) beginBlock() {
	b := &Block{}
	c.Blocks.Push(b)
}

func (c *Frame) endBlock() {
	c.Blocks.Pop()
}
