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

	assert.Len(tokens, 2)
	assert.Equal(TT_IDENT, tokens[0].Type())
	assert.Equal("axxx", tokens[0].String())
	assert.Equal(TT_IDENT, tokens[1].Type())
	assert.Equal("xx", tokens[1].String())
}
