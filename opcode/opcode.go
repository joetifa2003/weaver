package opcode

import "fmt"

type OpCode int

const (
	OP_CONSTANT OpCode = iota // arg1: constant index

	OP_LET
	OP_LOAD

	OP_ADD
	OP_MUL
	OP_DIV
	OP_MOD
	OP_SUB

	OP_ECHO
)

func PrintOpcodes(instructions []OpCode) string {
	var out string

	for i := 0; i < len(instructions); i++ {
		instr := instructions[i]
		switch instr {
		case OP_ADD:
			out += "add\n"

		case OP_MUL:
			out += "mul\n"

		case OP_DIV:
			out += "div\n"

		case OP_MOD:
			out += "mod\n"

		case OP_SUB:
			out += "sub\n"

		case OP_CONSTANT:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("const %d\n", op1)
		case OP_LET:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("let %d \n", op1)
		case OP_ECHO:
			out += fmt.Sprintf("echo\n")
		case OP_LOAD:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("load %d\n", op1)
		}
	}

	return out
}
