package lexer

import (
	"fmt"
	"testing"
)

const (
	TT_IDENT = iota
	TT_WHITESPACE
)

func TestSimpleLexer(t *testing.T) {
	l := New([]Pattern{
		{TT_IDENT, "[a-zA-Z]+"},
		{TT_WHITESPACE, "\\s+"},
	}, WithElide(TT_WHITESPACE))

	tokens, err := l.Lex("axxx   xx")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(tokens)
}
