package vm

import "fmt"

type opDef struct {
	op1 ValueType
	op2 ValueType
	op  opFunc
}

type opFunc = func(v *Value, other *Value, res *Value)

type opTable [valueTypeEnd][valueTypeEnd]opFunc

func (t *opTable) Call(v *Value, other *Value, res *Value) {
	t[v.VType][other.VType](v, other, res)
}

func initOpTable(op string, defs ...opDef) opTable {
	var table [valueTypeEnd][valueTypeEnd]opFunc

	illegalOp := func(v *Value, other *Value, res *Value) {
		panic(fmt.Sprintf("illegal operation %s %s %s", v.VType, op, other.VType))
	}

	for i := range table {
		for j := range table[i] {
			table[i][j] = illegalOp
		}
	}

	for _, def := range defs {
		table[def.op1][def.op2] = def.op
	}

	return table
}
