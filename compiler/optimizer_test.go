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
			name: "let_const",
			in: []opcode.OpCode{
				opcode.OP_CONST,
				opcode.OpCode(0),
				opcode.OP_STORE,
				opcode.OpCode(1),
			},
			out: []opcode.OpCode{
				opcode.OP_CONST_STORE,
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
