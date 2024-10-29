package compiler

import "github.com/joetifa2003/weaver/opcode"

type Frame struct {
	Register     []*Reg
	Vars         map[string]int
	Instructions []opcode.OpCode
}

func newFrame() *Frame {
	return &Frame{
		Vars: map[string]int{},
	}
}

type Reg struct {
	Index int
	Free  bool
}

func (r *Reg) free() {
	r.Free = true
}

func (f *Frame) allocReg() *Reg {
	for _, r := range f.Register {
		if r.Free {
			r.Free = false
			return r
		}
	}

	f.Register = append(f.Register, &Reg{
		Free:  false,
		Index: len(f.Register),
	})
	return f.Register[len(f.Register)-1]
}

func (f *Frame) allocVar(name string) *Reg {
	reg := f.allocReg()
	f.Vars[name] = reg.Index

	return reg
}

func (f *Frame) addInstructions(instructions []opcode.OpCode) {
	f.Instructions = append(f.Instructions, instructions...)
}
