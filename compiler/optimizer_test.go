package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joetifa2003/weaver/opcode"
)

func TestOptimizer(t *testing.T) {
	assert := require.New(t)

	c := New()

	testCases := []struct {
		name string
		in   []opcode.OpCode
		out  []opcode.OpCode
	}{
		{
			name: "load_const_add_let",
			in: []opcode.OpCode{
				opcode.OP_LOAD,
				opcode.OpCode(0),
				opcode.OP_CONST,
				opcode.OpCode(0),
				opcode.OP_ADD,
				opcode.OP_STORE,
				opcode.OpCode(0),
				opcode.OP_POP,
			},
			out: []opcode.OpCode{
				opcode.OP_LOAD_CONST_ADD_LET,
				opcode.OpCode(0),
				opcode.OpCode(0),
				opcode.OpCode(0),
			},
		},
		{
			name: "load_const_add_let2",
			in: []opcode.OpCode{
				opcode.OP_LOAD,
				opcode.OpCode(0),
				opcode.OP_CONST,
				opcode.OpCode(0),
				opcode.OP_ADD,
				opcode.OP_LET,
				opcode.OpCode(0),
			},
			out: []opcode.OpCode{
				opcode.OP_LOAD_CONST_ADD_LET,
				opcode.OpCode(0),
				opcode.OpCode(0),
				opcode.OpCode(0),
			},
		},
		{
			name: "let_const",
			in: []opcode.OpCode{
				opcode.OP_CONST,
				opcode.OpCode(0),
				opcode.OP_LET,
				opcode.OpCode(1),
			},
			out: []opcode.OpCode{
				opcode.OP_CONST_LET,
				opcode.OpCode(0),
				opcode.OpCode(1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(tc.out, c.optimize(tc.in))
		})
	}
}
