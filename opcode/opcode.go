package opcode

import (
	"fmt"
	"iter"
	"slices"
)

type OpCode int

const (
	OP_CONST OpCode = iota // arg1: constant index
	OP_POP

	OP_LABEL

	OP_CALL
	OP_RET
	OP_HALT

	OP_STORE // arg1: variable index
	OP_LOAD  // arg1: variable index

	OP_JUMP  // arg1: jump offset
	OP_JUMPF // arg1: jump offset

	OP_ADD
	OP_MUL
	OP_DIV
	OP_MOD
	OP_SUB
	OP_LT
	OP_LTE
	OP_EQ
	OP_GT
	OP_GTE

	OP_ECHO

	OP_CONST_STORE          // arg1: constant index; arg2: variable index
	OP_LOAD_CONST_ADD       // arg1: constant index; arg2: variable index
	OP_LOAD_CONST_ADD_STORE // arg1: constant index; arg2: variable index; arg3: variable index
	OP_CONST_ADD            // arg1: constant index;
)

type OpCodeDef struct {
	Code      OpCode
	Name      string
	ArgsCount int
}

var opCodeDefs = map[OpCode]OpCodeDef{
	OP_CONST:                {OP_CONST, "const", 1},
	OP_POP:                  {OP_POP, "pop", 0},
	OP_CALL:                 {OP_CALL, "call", 1},
	OP_RET:                  {OP_RET, "ret", 0},
	OP_HALT:                 {OP_HALT, "halt", 0},
	OP_STORE:                {OP_STORE, "store", 1},
	OP_LOAD:                 {OP_LOAD, "load", 1},
	OP_JUMP:                 {OP_JUMP, "jmp", 1},
	OP_JUMPF:                {OP_JUMPF, "jmpf", 1},
	OP_ADD:                  {OP_ADD, "add", 0},
	OP_MUL:                  {OP_MUL, "mul", 0},
	OP_DIV:                  {OP_DIV, "div", 0},
	OP_MOD:                  {OP_MOD, "mod", 0},
	OP_SUB:                  {OP_SUB, "sub", 0},
	OP_LT:                   {OP_LT, "lt", 0},
	OP_EQ:                   {OP_EQ, "eq", 0},
	OP_ECHO:                 {OP_ECHO, "echo", 0},
	OP_LABEL:                {OP_ECHO, "label", 1},
	OP_CONST_STORE:          {OP_CONST_STORE, "sconst", 2},
	OP_LOAD_CONST_ADD:       {OP_LOAD_CONST_ADD, "lcadd", 2},
	OP_LOAD_CONST_ADD_STORE: {OP_LOAD_CONST_ADD_STORE, "lcastore", 3},
	OP_CONST_ADD:            {OP_CONST_ADD, "cadd", 1},
}

type DecodedOpCode struct {
	Op   OpCode
	Addr int
	Name string
	Args []OpCode
}

func OpCodeIterator(instruction []OpCode, skip ...OpCode) iter.Seq2[int, DecodedOpCode] {
	return func(yield func(int, DecodedOpCode) bool) {

		for i := 0; i < len(instruction); i++ {
			res := DecodedOpCode{}
			def, ok := opCodeDefs[instruction[i]]
			if !ok {
				panic(fmt.Sprintf("unknown instruction %d", instruction))
			}
			res.Name = def.Name
			res.Op = instruction[i]
			res.Addr = i

			for range def.ArgsCount {
				i++
				res.Args = append(res.Args, instruction[i])
			}

			if slices.Contains(skip, res.Op) {
				continue
			}

			if !yield(i, res) {
				return
			}
		}
	}
}

func PrintOpcodes(instructions []OpCode) string {
	var out string

	maxNameLength := 0
	for _, instr := range OpCodeIterator(instructions) {
		if len(instr.Name) > maxNameLength {
			maxNameLength = len(instr.Name)
		}
	}

	for _, instr := range OpCodeIterator(instructions) {
		argsStr := ""
		for _, arg := range instr.Args {
			argsStr = fmt.Sprintf("%s%06d ", argsStr, arg)
		}
		out += fmt.Sprintf("%06d: %-*s %s\n", instr.Addr, maxNameLength, instr.Name, argsStr)
	}

	return out
}
