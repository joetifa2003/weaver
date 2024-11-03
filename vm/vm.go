package vm

import (
	"fmt"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/value"
)

type Frame struct {
	*compiler.Frame
	ip int
}

type VM struct {
	constants []value.Value
	stack     *ds.Stack[value.Value]
	frames    *ds.Stack[*Frame]
}

func New(constants []value.Value, mainFrame *compiler.Frame) *VM {
	return &VM{
		constants: constants,
		stack:     ds.NewStack[value.Value](),
		frames:    ds.NewStack(&Frame{mainFrame, 0}),
	}
}

func (v *VM) Run() {
	v.initializeFrame()

	for v.currentFrame().ip < len(v.currentFrame().Instructions) {
		switch v.currentFrame().Instructions[v.currentFrame().ip] {
		case opcode.OP_CONSTANT:
			v.incrementIP()
			index := v.currentInstruction()
			v.stack.Push(v.constants[index])
			v.incrementIP()
		case opcode.OP_LET:
			v.incrementIP()
			index := v.currentInstruction()
			v.stack.Set(int(index), v.stack.Pop())
			v.incrementIP()
		case opcode.OP_ECHO:
			value := v.stack.Pop()
			fmt.Println(value.String())
			v.incrementIP()
		case opcode.OP_ADD:
			left := v.stack.Pop()
			right := v.stack.Pop()
			v.stack.Push(value.NewInt(left.GetInt() + right.GetInt()))
			v.incrementIP()
		case opcode.OP_MUL:
			left := v.stack.Pop()
			right := v.stack.Pop()
			v.stack.Push(value.NewInt(left.GetInt() * right.GetInt()))
			v.incrementIP()
		case opcode.OP_LOAD:
			v.incrementIP()
			index := v.currentInstruction()
			val := v.stack.Get(int(index))
			v.stack.Push(val)
			v.incrementIP()
		default:
			panic(fmt.Sprintf("unimplemented %s", v.currentInstruction()))
		}
	}
}

func (v *VM) currentFrame() *Frame {
	return v.frames.Peek()
}

func (v *VM) incrementIP() {
	v.currentFrame().ip++
}

func (v *VM) currentInstruction() opcode.OpCode {
	return v.currentFrame().Instructions[v.currentFrame().ip]
}

func (v *VM) initializeFrame() {
	f := v.currentFrame()
	for range f.Vars {
		v.stack.Push(value.Value{})
	}
}
