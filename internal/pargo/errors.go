package pargo

import "parser-comb/lexer"

type ParseError struct {
	Expected string
	Found    lexer.Token
	Source   string
}

func (e ParseError) Error() string {
	return "Expected " + e.Expected + " but found " + e.Found.String()
}

func NewParseError(src string, expected string, found lexer.Token) ParseError {
	return ParseError{
		Expected: expected,
		Found:    found,
		Source:   src,
	}
}
