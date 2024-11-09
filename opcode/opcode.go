package opcode

import "fmt"

type OpCode int

const (
	OP_CONSTANT OpCode = iota // arg1: constant index
	OP_POP

	OP_CALL
	OP_RET

	OP_LET    // arg1: variable index
	OP_ASSIGN // arg1: variable index
	OP_LOAD   // arg1: variable index

	OP_JUMP  // arg1: jump offset
	OP_JUMPF // arg1: jump offset

	OP_ADD
	OP_MUL
	OP_DIV
	OP_MOD
	OP_SUB
	OP_LT
	OP_EQ

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

		case OP_JUMP:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("jump %d\n", op1)

		case OP_JUMPF:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("jumpf %d\n", op1)

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

		case OP_ASSIGN:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("assign %d\n", op1)

		case OP_LT:
			out += "lt\n"

		case OP_EQ:
			out += "eq\n"

		case OP_POP:
			out += "pop\n"

		case OP_CALL:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("call %d\n", op1)

		case OP_RET:
			out += "ret\n"

		default:
			out += fmt.Sprintf("unknown opcode %d\n", instr)
		}
	}

	return out
}
