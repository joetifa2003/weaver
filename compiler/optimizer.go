package compiler

import (
	"github.com/joetifa2003/weaver/opcode"
)

type Rule func([]opcode.DecodedOpCode) (bool, int)

func eq(op opcode.OpCode) Rule {
	return func(instructions []opcode.DecodedOpCode) (bool, int) {
		if len(instructions) == 0 {
			return false, 0
		}

		if instructions[0].Op != op {
			return false, 0
		}

		return true, 1
	}
}

func seq(rules ...Rule) Rule {
	return func(instructions []opcode.DecodedOpCode) (bool, int) {
		if len(instructions) < len(rules) {
			return false, 0
		}

		prev := 0
		for _, rule := range rules {
			matched, eaten := rule(instructions[prev:])
			if !matched {
				return false, 0
			}
			prev += eaten
		}

		return true, prev
	}
}

func or(rules ...Rule) Rule {
	return func(instructions []opcode.DecodedOpCode) (bool, int) {
		for _, rule := range rules {
			matched, eaten := rule(instructions)
			if matched {
				return true, eaten
			}
		}

		return false, 0
	}
}

type Optimizer struct {
	Seq Rule
	Fn  func([]opcode.DecodedOpCode) []opcode.OpCode
}

const MAX_SEQ_LEN = 4

var optimizers = []Optimizer{
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_CONST), eq(opcode.OP_ADD), eq(opcode.OP_STORE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_CONST_ADD_STORE,
				doc[1].Args[0],
				doc[0].Args[0],
				doc[3].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_CONST), eq(opcode.OP_ADD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_CONST_ADD,
				doc[1].Args[0],
				doc[0].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_CONST), eq(opcode.OP_ADD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_CONST_ADD,
				doc[0].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_CONST), eq(opcode.OP_STORE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_CONST_STORE,
				doc[0].Args[0],
				doc[1].Args[0],
			}
		},
	},
}

func (c *Compiler) optimize(instructions []opcode.OpCode) []opcode.OpCode {
	var optimized []opcode.OpCode

	window := []opcode.DecodedOpCode{}
	for _, instr := range opcode.OpCodeIterator(instructions) {
		window = append(window, instr)

		for _, opt := range optimizers {
			matched, eaten := opt.Seq(window)
			if matched {
				optimized = append(optimized, opt.Fn(window[:eaten])...)
				window = window[eaten:]
			}
		}

		if len(window) >= MAX_SEQ_LEN {
			first := window[0]
			optimized = append(optimized, first.Op)
			optimized = append(optimized, first.Args...)
			window = window[1:]
		}
	}

	for _, instr := range window {
		optimized = append(optimized, instr.Op)
		optimized = append(optimized, instr.Args...)
	}

	return optimized
}
