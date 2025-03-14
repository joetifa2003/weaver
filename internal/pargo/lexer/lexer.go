package lexer

import (
	"fmt"
	"regexp"
)

type Lexer interface {
	Lex(input string) ([]Token, error)
}

type Token interface {
	String() string
	Type() int
	Location() Location
}

type Location struct {
	Line   int
	Column int
}

func (l Location) IsAfter(other Location) bool {
	return l.Line > other.Line || (l.Line == other.Line && l.Column > other.Column)
}

type RegexToken struct {
	Loc   Location
	Ttype int
	Lit   string
}

func (t RegexToken) String() string { return t.Lit }

func (t RegexToken) Type() int { return t.Ttype }

func (t RegexToken) Location() Location { return t.Loc }

type RegexLexer struct {
	patterns   []Pattern
	elide      map[int]struct{}
	transforms map[int]func(string) string
}

type Pattern struct {
	TokenType int
	Regex     string
}

type RegexLexerOption func(*RegexLexer)

func WithElide(elide ...int) RegexLexerOption {
	return func(l *RegexLexer) {
		for _, ttype := range elide {
			l.elide[ttype] = struct{}{}
		}
	}
}

func WithTransform(ttype int, transform func(string) string) RegexLexerOption {
	return func(l *RegexLexer) {
		l.transforms[ttype] = transform
	}
}

func New(patterns []Pattern, options ...RegexLexerOption) *RegexLexer {
	l := &RegexLexer{
		patterns:   patterns,
		elide:      map[int]struct{}{},
		transforms: map[int]func(string) string{},
	}

	for _, option := range options {
		option(l)
	}

	return l
}

func (l *RegexLexer) Lex(input string) ([]Token, error) {
	var tokens []Token
	location := Location{Line: 1, Column: 1}

	for len(input) > 0 {
		matched := false
		for _, pattern := range l.patterns {
			re := regexp.MustCompile("^" + pattern.Regex)
			if match := re.FindString(input); match != "" {
				if _, ok := l.elide[pattern.TokenType]; !ok {
					lit := match
					if transform, ok := l.transforms[pattern.TokenType]; ok {
						lit = transform(match)
					}
					tokens = append(tokens, RegexToken{
						Loc:   location,
						Ttype: pattern.TokenType,
						Lit:   lit,
					})
				}

				advance := len(match)
				input = input[advance:]
				for _, char := range match {
					if char == '\n' {
						location.Line++
						location.Column = 1
					} else {
						location.Column++
					}
				}
				matched = true
				break
			}
		}
		if !matched {
			return nil, fmt.Errorf("unexpected character '%s' at line %d, column %d", input[:1], location.Line, location.Column)
		}
	}

	return tokens, nil
}
