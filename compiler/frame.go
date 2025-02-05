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
	Scope     VarScope
	Name      string
	Index     int
	Parent    *Var
	Available bool
}

func (v *Var) load() []opcode.OpCode {
	switch v.Scope {
	case VarScopeLocal:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeLocal,
			opcode.OpCode(v.Index),
		}
	case VarScopeFree:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeFree,
			opcode.OpCode(v.Index),
		}
	case VarScopeGlobal:
		return []opcode.OpCode{
			opcode.OP_LOAD,
			opcode.ScopeTypeGlobal,
			opcode.OpCode(v.Index),
		}
	default:
		panic(fmt.Sprintf("unknown scope %d", v.Scope))
	}
}

func (v *Var) store() []opcode.OpCode {
	switch v.Scope {
	case VarScopeLocal:
		return []opcode.OpCode{
			opcode.OP_STORE,
			opcode.OpCode(v.Index),
		}
	case VarScopeFree:
		return []opcode.OpCode{
			opcode.OP_STORE_FREE,
			opcode.OpCode(v.Index),
		}
	case VarScopeGlobal:
		return []opcode.OpCode{
			opcode.OP_STORE_GLOBAL,
			opcode.OpCode(v.Index),
		}
	default:
		panic(fmt.Sprintf("unknown scope %d", v.Scope))
	}
}

type Block struct {
	Vars []*Var
}

func (c *Frame) addInstructions(instructions []opcode.OpCode) {
	c.Instructions = append(c.Instructions, instructions...)
}

func (c *Frame) defineVar(name string) int {
	for _, v := range c.Vars {
		if v.Available {
			v.Available = false
			v.Scope = VarScopeLocal
			v.Parent = nil
			v.Name = name
			c.Blocks.Peek().Vars = append(c.Blocks.Peek().Vars, v)
			return v.Index
		}
	}

	v := &Var{Name: name, Index: len(c.Vars), Scope: VarScopeLocal}
	c.Vars = append(c.Vars, v)
	c.Blocks.Peek().Vars = append(c.Blocks.Peek().Vars, v)
	return len(c.Vars) - 1
}

func (c *Frame) defineFreeVar(name string, parent *Var) *Var {
	v := &Var{Name: name, Index: len(c.FreeVars), Scope: VarScopeFree, Parent: parent}
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
		v, err := c.Parent.resolve(name)
		if err != nil {
			return nil, err
		}

		return c.defineFreeVar(name, v), nil
	}

	return nil, fmt.Errorf("%w: %s", ErrVarNotFound, name)
}

func (c *Frame) beginBlock() {
	b := &Block{}
	c.Blocks.Push(b)
}

func (c *Frame) endBlock() {
	block := c.Blocks.Pop()
	for _, v := range block.Vars {
		v.Available = true
	}
}
