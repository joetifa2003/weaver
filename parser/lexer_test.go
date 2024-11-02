package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	assert := require.New(t)
	l := newLexer()

	tokens, err := l.Lex("axxx xx")
	assert.NoError(err)

	assert.Equal(2, len(tokens))
	assert.Equal(int(TT_IDENT), tokens[0].Type())
	assert.Equal("axxx", tokens[0].String())
	assert.Equal(int(TT_IDENT), tokens[1].Type())
	assert.Equal("xx", tokens[1].String())
}
