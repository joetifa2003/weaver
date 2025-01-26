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

func repeat(rule Rule, n int) Rule {
	rules := make([]Rule, n)
	for i := range rules {
		rules[i] = rule
	}
	return seq(rules...)
}

func some(rule Rule, atLeast int) Rule {
	return func(instructions []opcode.DecodedOpCode) (bool, int) {
		someEaten := 0

		for {
			matched, eaten := rule(instructions[someEaten:])
			if !matched {
				if someEaten < atLeast {
					return false, 0
				}

				return true, someEaten
			}

			someEaten += eaten
		}
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
		Seq: seq(eq(opcode.OP_NOT), eq(opcode.OP_JUMP_F)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_JUMP_T,
				doc[1].Args[0],
			}
		},
	},
	{
		Seq: seq(
			eq(opcode.OP_LOAD),
			eq(opcode.OP_CONST),
			eq(opcode.OP_ADD),
			eq(opcode.OP_LET),
		),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_CONST_ADD_LET,
				doc[1].Args[0],
				doc[0].Args[0],
				doc[3].Args[0],
			}
		},
	},
	{
		Seq: seq(repeat(eq(opcode.OP_LOAD), 2), eq(opcode.OP_LT)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_LOAD_LOAD_LT,
				doc[1].Args[0],
				doc[0].Args[0],
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
		Seq: seq(eq(opcode.OP_CONST), eq(opcode.OP_LET)),
		Fn: func(doc []opcode.DecodedOpCode) []opcode.OpCode {
			return []opcode.OpCode{
				opcode.OP_CONST_LET,
				doc[0].Args[0],
				doc[1].Args[0],
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
