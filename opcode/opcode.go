package opcode

import "fmt"

type OpCode int

const (
	OP_CONSTANT OpCode = iota // (const idx, reg)
	OP_ADD                    // (reg1, reg2) => reg3
	OP_LET                    // (reg1, reg2)
	OP_ECHO                   // (reg)
)

func PrintOpcodes(instructions []OpCode) string {
	var out string

	for i := 0; i < len(instructions); i++ {
		instr := instructions[i]
		switch instr {
		case OP_ADD:
			i++
			op1 := i
			i++
			op2 := i
			out += fmt.Sprintf("add %d %d\n", op1, op2)
		case OP_CONSTANT:
			i++
			op1 := instructions[i]
			i++
			op2 := instructions[i]
			out += fmt.Sprintf("const %d %d\n", op1, op2)
		case OP_LET:
			i++
			op1 := instructions[i]
			i++
			op2 := instructions[i]
			out += fmt.Sprintf("let %d %d\n", op1, op2)
		case OP_ECHO:
			i++
			op1 := instructions[i]
			out += fmt.Sprintf("echo %d\n", op1)
		}
	}

	return out
}
