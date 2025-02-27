package pargo

import (
	"fmt"

	"github.com/joetifa2003/weaver/internal/pargo/lexer"
)

type ParseError struct {
	Expected string
	Found    lexer.Token
	Source   string
}

func (e ParseError) Error() string {
	location := e.Found.Location()

	return fmt.Sprintf("expected %s but found %s at %d:%d", e.Expected, e.Found.String(), location.Line, location.Column)
}

func (e ParseError) Location() lexer.Location {
	return e.Found.Location()
}

func NewParseError(src string, expected string, found lexer.Token) ParseError {
	return ParseError{
		Expected: expected,
		Found:    found,
		Source:   src,
	}
}
