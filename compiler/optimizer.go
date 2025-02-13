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

type Optimizer struct {
	Seq Rule
	Fn  func([]opcode.DecodedOpCode) []opcode.OpCode
}

var optimizers = []Optimizer{
	{
		Seq: seq(eq(opcode.OP_STORE), eq(opcode.OP_POP)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LET,
				doc[0].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_NOT), eq(opcode.OP_PJUMP_F)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_PJUMP_T,
				doc[1].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_NOT), eq(opcode.OP_PJUMP_T)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_PJUMP_F,
				doc[1].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_NOT), eq(opcode.OP_JUMP_F)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_JUMP_T,
				doc[1].Args[0],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_NOT), eq(opcode.OP_JUMP_T)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_JUMP_F,
				doc[1].Args[0],
			}
		},
	},

	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_ADD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_ADD,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_ADD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_ADD,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_SUB)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_SUB,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_SUB)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_SUB,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_MUL)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_MUL,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_MUL)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_MUL,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_DIV)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_DIV,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_DIV)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_DIV,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_MOD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_MOD,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_MOD)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_MOD,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_LT)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_LT,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LT)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LT,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_LTE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_LTE,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LTE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LTE,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_GT)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_GT,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_GT)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_GT,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_GTE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_GTE,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_GTE)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_GTE,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_EQ)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_EQ,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_EQ)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_EQ,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_LOAD), eq(opcode.OP_NEQ)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_NEQ,
				doc[0].Args[0],
				doc[0].Args[1],
				doc[1].Args[0],
				doc[1].Args[1],
			}
		},
	},
	{
		Seq: seq(eq(opcode.OP_LOAD), eq(opcode.OP_NEQ)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_NEQ,
				doc[0].Args[0],
				doc[0].Args[1],
			}
		},
	},
}

func (c *Compiler) optimizePass(instructions []opcode.OpCode) ([]opcode.OpCode, bool) {
	dirty := false

	for _, opt := range optimizers {
		decodedInstructions := opcode.DecodeInstructions(instructions)
		optimized := make([]opcode.OpCode, 0, len(instructions))

		for len(decodedInstructions) > 0 {
			matched, eaten := opt.Seq(decodedInstructions)
			if matched {
				optimized = append(optimized, opt.Fn(decodedInstructions[:eaten])...)
				decodedInstructions = decodedInstructions[eaten:]
				dirty = true
			} else {
				first := decodedInstructions[0]
				optimized = append(optimized, first.Op)
				optimized = append(optimized, first.Args...)
				decodedInstructions = decodedInstructions[1:]
			}
		}

		instructions = optimized
	}

	return instructions, dirty
}

func (c *Compiler) optimize(instructions []opcode.OpCode) []opcode.OpCode {
	// return instructions
	for {
		optimized, dirty := c.optimizePass(instructions)
		if !dirty {
			return optimized
		}
		instructions = optimized
	}
}
