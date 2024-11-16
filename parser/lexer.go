package parser

import "github.com/joetifa2003/weaver/internal/pargo/lexer"

const (
	TT_IDENT int = iota
	TT_INT
	TT_FLOAT
	TT_STRING
	TT_PLUS
	TT_MINUS
	TT_SYMBOL
	TT_MULTIPLY
	TT_DIVIDE
	TT_EQUAL
	TT_LESS_THAN
	TT_GREATER_THAN
	TT_LESS_THAN_EQUAL
	TT_GREATER_THAN_EQUAL
	TT_NOT_EQUAL
	TT_AND
	TT_OR
	TT_LPAREN
	TT_RPAREN
	TT_ASSIGN
	TT_VARDECL
	TT_WHITESPACE
)

func newLexer() *lexer.RegexLexer {
	return lexer.New(
		[]lexer.Pattern{
			{TokenType: TT_IDENT, Regex: "[a-zA-Z_]+"},
			{TokenType: TT_FLOAT, Regex: "[0-9]+\\.[0-9]+"},
			{TokenType: TT_INT, Regex: "[0-9]+"},
			{TokenType: TT_STRING, Regex: `"(?:[^"\\]|\\.)*"`},
			// ========== operators ==========
			{TokenType: TT_SYMBOL, Regex: "{"},
			{TokenType: TT_SYMBOL, Regex: "}"},
			{TokenType: TT_SYMBOL, Regex: "\\|"},
			{TokenType: TT_SYMBOL, Regex: ","},

			{TokenType: TT_SYMBOL, Regex: "%"},
			{TokenType: TT_PLUS, Regex: "\\+"},
			{TokenType: TT_MINUS, Regex: "-"},
			{TokenType: TT_MULTIPLY, Regex: "\\*"},
			{TokenType: TT_DIVIDE, Regex: "/"},
			{TokenType: TT_EQUAL, Regex: "=="},
			{TokenType: TT_VARDECL, Regex: ":="},
			{TokenType: TT_ASSIGN, Regex: "="},
			{TokenType: TT_LESS_THAN_EQUAL, Regex: "<="},
			{TokenType: TT_GREATER_THAN_EQUAL, Regex: ">="},
			{TokenType: TT_LESS_THAN, Regex: "<"},
			{TokenType: TT_GREATER_THAN, Regex: ">"},
			{TokenType: TT_NOT_EQUAL, Regex: "!="},
			{TokenType: TT_AND, Regex: "&&"},
			{TokenType: TT_OR, Regex: "\\|\\|"},
			{TokenType: TT_LPAREN, Regex: "\\("},
			{TokenType: TT_RPAREN, Regex: "\\)"},
			// ===============================
			{TokenType: TT_WHITESPACE, Regex: "\\s+"},
		},
		lexer.WithEllide(TT_WHITESPACE),
		lexer.WithTransform(TT_STRING, func(s string) string {
			return s[1 : len(s)-1]
		}),
	)
}
