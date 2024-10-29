package lexer

import (
	"fmt"
	"testing"
)

func TestSimpleLexer(t *testing.T) {
	l := NewSimple()

	tokens, err := l.Lex("xyz yyy")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(tokens)
}
