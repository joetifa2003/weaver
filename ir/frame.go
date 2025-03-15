package ir

import (
	"fmt"

	"github.com/joetifa2003/weaver/internal/pkg/ds"
)

type frame struct {
	Parent     *frame
	Vars       []*basicVar
	FreeVars   []*basicVar
	Blocks     *ds.Stack[*basicBlock]
	Statements []Statement
}

func NewFrame(parent *frame) *frame {
	return &frame{
		Parent: parent,
		Vars:   []*basicVar{},
		Blocks: ds.NewStack[*basicBlock](&basicBlock{}),
	}
}

func (c *frame) define(name string) *basicVar {
	if name == "" {
		name = fmt.Sprintf("__$b%dv%d", c.Blocks.Len()-1, len(c.Vars))
	}

	block := c.currentBlock()

	for _, v := range c.Vars {
		if v.Free {
			v.Name = name
			block.vars = append(block.vars, v)
			v.Free = false
			return v
		}
	}

	v := &basicVar{Name: name, Scope: VarScopeLocal, Index: len(c.Vars)}
	c.Vars = append(c.Vars, v)
	block.vars = append(block.vars, v)
	return v
}

func (c *frame) defineFreeVar(name string, parent *basicVar) *basicVar {
	parent.Ref = true
	v := &basicVar{Name: name, Scope: VarScopeFree, Index: len(c.FreeVars), Parent: parent}
	c.FreeVars = append(c.FreeVars, v)
	return v
}

func (c *frame) resolve(name string) (*basicVar, error) {
	for _, v := range c.Vars {
		if v.Name == name {
			return v, nil
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

	return nil, fmt.Errorf("cannot find variable %s", name)
}

func (c *frame) currentBlock() *basicBlock {
	return c.Blocks.Peek()
}

func (c *frame) pushBlock() *basicBlock {
	b := &basicBlock{}
	c.Blocks.Push(b)
	return b
}

func (c *frame) popBlock() *basicBlock {
	b := c.Blocks.Pop()
	for _, v := range b.vars {
		v.Free = true
	}
	return b
}

func (c *frame) pushStmt(s Statement) {
	c.Statements = append(c.Statements, s)
}

func (c *frame) export() FrameExpr {
	f := FrameExpr{
		VarCount: len(c.Vars),
		Body:     c.Statements,
	}

	for _, v := range c.FreeVars {
		f.FreeVars = append(f.FreeVars, v.Parent.export())
	}

	return f
}

type VarScope int

const (
	VarScopeLocal VarScope = iota
	VarScopeFree
)

func (s VarScope) String() string {
	switch s {
	case VarScopeLocal:
		return "local"
	case VarScopeFree:
		return "free"
	default:
		panic(fmt.Sprintf("unknown scope %d", s))
	}
}

type basicVar struct {
	Scope VarScope
	Name  string
	Index int
	Free  bool
	Ref   bool

	Parent *basicVar
}

func (b *basicVar) assign(expr Expr) VarAssignExpr {
	return VarAssignExpr{
		Var: Var{
			Idx:   b.Index,
			Scope: b.Scope,
		},
		Value: expr,
	}
}

func (b *basicVar) assignStmt(expr Expr) Statement {
	return ExpressionStmt{
		Expr: b.assign(expr),
	}
}

func (b *basicVar) load() LoadExpr {
	return LoadExpr{
		Var: Var{
			Idx:   b.Index,
			Scope: b.Scope,
			Ref:   b.Ref,
		},
	}
}

func (c *basicVar) export() Var {
	return Var{
		Idx:   c.Index,
		Scope: c.Scope,
	}
}

type basicBlock struct {
	vars       []*basicVar
	statements []Statement
}

func (b *basicBlock) pushStmt(s Statement) {
	b.statements = append(b.statements, s)
}

func (b *basicBlock) export() BlockStmt {
	return BlockStmt{
		Statements: b.statements,
	}
}

func (b *basicBlock) freeAll() {
	for _, v := range b.vars {
		v.Free = true
	}
}
