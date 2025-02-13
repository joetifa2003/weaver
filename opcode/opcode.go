package opcode

import (
	"fmt"
	"iter"
	"slices"
)

type OpCode int

const (
	OP_POP OpCode = iota

	OP_LABEL

	OP_CALL
	OP_RET
	OP_HALT
	OP_INC_LOCAL

	OP_LET          // arg1: variable index
	OP_STORE        // arg1: variable index
	OP_STORE_FREE   // arg1: variable index
	OP_STORE_GLOBAL // arg1: variable index
	OP_LOAD         // arg1: scope; arg2: variable index
	OP_STORE_IDX

	OP_JUMP    // arg1: jump offset
	OP_PJUMP_F // arg1: jump offset
	OP_PJUMP_T // arg1: jump offset
	OP_JUMP_F  // arg1: jump offset
	OP_JUMP_T  // arg1: jump offset

	OP_ARRAY // initialize array
	OP_APUSH // push value to array
	OP_INDEX // push array/object value at index to stack

	OP_OBJ   // initialize object
	OP_OPUSH // push value to object

	OP_ADD // +
	OP_SUB // -
	OP_MUL // *
	OP_DIV // /
	OP_MOD // %
	OP_LT  // <
	OP_LTE // <=
	OP_GT  // >
	OP_GTE // >=
	OP_EQ  // ==
	OP_NEQ // !=
	OP_NOT // !

	OP_ECHO
	OP_FUNC // arg1: constant index; arg2: free variables count

	// Super instructions.
	OP_LOAD_LOAD_ADD // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_ADD      // arg1: scope; arg2: index
	OP_LOAD_LOAD_SUB // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_SUB      // arg1: scope; arg2: index
	OP_LOAD_LOAD_MUL // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_MUL      // arg1: scope; arg2: index
	OP_LOAD_LOAD_DIV // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_DIV      // arg1: scope; arg2: index
	OP_LOAD_LOAD_MOD // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_MOD      // arg1: scope; arg2: index
	OP_LOAD_LOAD_LT  // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_LT       // arg1: scope; arg2: index
	OP_LOAD_LOAD_LTE // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 OP_INDEX
	OP_LOAD_LTE      // arg1: scope; arg2: index
	OP_LOAD_LOAD_GT  // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_GT       // arg1: scope; arg2: index
	OP_LOAD_LOAD_GTE // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_GTE      // arg1: scope; arg2: index
	OP_LOAD_LOAD_EQ  // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_EQ       // arg1: scope; arg2: index
	OP_LOAD_LOAD_NEQ // arg1: v1 scope; arg2: v1 index; arg3: v2 scope; arg4: v2 index
	OP_LOAD_NEQ      // arg1: scope; arg2: index
)

type ScopeType = OpCode

const (
	ScopeTypeLocal ScopeType = iota
	ScopeTypeFree
	ScopeTypeGlobal
	ScopeTypeConst
)

type OpCodeDef struct {
	Code      OpCode
	Name      string
	ArgsCount int
}

var opCodeDefs = map[OpCode]OpCodeDef{
	OP_POP:          {OP_POP, "pop", 0},
	OP_CALL:         {OP_CALL, "call", 1},
	OP_RET:          {OP_RET, "ret", 0},
	OP_HALT:         {OP_HALT, "halt", 0},
	OP_STORE:        {OP_STORE, "store", 1},
	OP_STORE_FREE:   {OP_STORE_FREE, "storef", 1},
	OP_STORE_GLOBAL: {OP_STORE_GLOBAL, "storeg", 1},
	OP_LET:          {OP_LET, "let", 1},
	OP_LOAD:         {OP_LOAD, "load", 2},
	OP_JUMP:         {OP_JUMP, "jmp", 1},
	OP_PJUMP_F:      {OP_PJUMP_F, "pjmpf", 1},
	OP_PJUMP_T:      {OP_PJUMP_T, "pjmpt", 1},
	OP_JUMP_F:       {OP_JUMP_F, "jmpf", 1},
	OP_JUMP_T:       {OP_JUMP_T, "jmpt", 1},
	OP_ADD:          {OP_ADD, "add", 0},
	OP_MUL:          {OP_MUL, "mul", 0},
	OP_DIV:          {OP_DIV, "div", 0},
	OP_MOD:          {OP_MOD, "mod", 0},
	OP_SUB:          {OP_SUB, "sub", 0},
	OP_LT:           {OP_LT, "lt", 0},
	OP_LTE:          {OP_LTE, "lte", 0},
	OP_GT:           {OP_GT, "gt", 0},
	OP_GTE:          {OP_GTE, "gte", 0},
	OP_EQ:           {OP_EQ, "eq", 0},
	OP_NEQ:          {OP_NEQ, "neq", 0},
	OP_NOT:          {OP_NOT, "not", 0},
	OP_APUSH:        {OP_APUSH, "apsh", 0},
	OP_ARRAY:        {OP_ARRAY, "arr", 0},
	OP_INDEX:        {OP_INDEX, "idx", 0},
	OP_FUNC:         {OP_FUNC, "func", 2},
	OP_STORE_IDX:    {OP_STORE_IDX, "storeidx", 0},
	OP_OBJ:          {OP_OBJ, "obj", 0},
	OP_OPUSH:        {OP_OPUSH, "opsh", 0},
	OP_LABEL:        {OP_LABEL, "label", 1},
	OP_INC_LOCAL:    {OP_INC_LOCAL, "incl", 1},

	OP_ECHO: {OP_ECHO, "echo", 0},

	OP_LOAD_LOAD_ADD: {OP_LOAD_LOAD_ADD, "lladd", 4},
	OP_LOAD_ADD:      {OP_LOAD_ADD, "ladd", 2},
	OP_LOAD_LOAD_SUB: {OP_LOAD_LOAD_SUB, "llsub", 4},
	OP_LOAD_SUB:      {OP_LOAD_SUB, "lsub", 2},
	OP_LOAD_LOAD_MUL: {OP_LOAD_LOAD_MUL, "llmul", 4},
	OP_LOAD_MUL:      {OP_LOAD_MUL, "lmul", 2},
	OP_LOAD_LOAD_DIV: {OP_LOAD_LOAD_DIV, "lldiv", 4},
	OP_LOAD_DIV:      {OP_LOAD_DIV, "ldiv", 2},
	OP_LOAD_LOAD_MOD: {OP_LOAD_LOAD_MOD, "llmod", 4},
	OP_LOAD_MOD:      {OP_LOAD_MOD, "lmod", 2},
	OP_LOAD_LOAD_LT:  {OP_LOAD_LOAD_LT, "lllt", 4},
	OP_LOAD_LT:       {OP_LOAD_LT, "llt", 2},
	OP_LOAD_LOAD_LTE: {OP_LOAD_LOAD_LTE, "lllte", 4},
	OP_LOAD_LTE:      {OP_LOAD_LTE, "lte", 2},
	OP_LOAD_LOAD_GT:  {OP_LOAD_LOAD_GT, "llgt", 4},
	OP_LOAD_GT:       {OP_LOAD_GT, "lgt", 2},
	OP_LOAD_LOAD_GTE: {OP_LOAD_LOAD_GTE, "llgte", 4},
	OP_LOAD_GTE:      {OP_LOAD_GTE, "lgte", 2},
	OP_LOAD_LOAD_EQ:  {OP_LOAD_LOAD_EQ, "lleq", 4},
	OP_LOAD_EQ:       {OP_LOAD_EQ, "leq", 2},
	OP_LOAD_LOAD_NEQ: {OP_LOAD_LOAD_NEQ, "llneq", 4},
	OP_LOAD_NEQ:      {OP_LOAD_NEQ, "lneq", 2},
}

type DecodedOpCode struct {
	Op   OpCode
	Addr int
	Args []OpCode
	Name string
}

func OpCodeIterator(instruction []OpCode, skip ...OpCode) iter.Seq2[int, DecodedOpCode] {
	return func(yield func(int, DecodedOpCode) bool) {
		for i := 0; i < len(instruction); i++ {
			res := DecodedOpCode{}
			def, ok := opCodeDefs[instruction[i]]
			if !ok {
				panic(fmt.Sprintf("unknown instruction %d", instruction[i]))
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

func DecodeInstructions(instructions []OpCode) []DecodedOpCode {
	//nolint:prealloc
	var decoded []DecodedOpCode
	for _, instr := range OpCodeIterator(instructions) {
		decoded = append(decoded, instr)
	}
	return decoded
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
		out += fmt.Sprintf("%06d: %-*s %s\n", instr.Addr, maxNameLength+10, instr.Name, argsStr)
	}

	return out
}
