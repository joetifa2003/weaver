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

type RegexToken struct {
	location Location
	ttype    int
	lit      string
}

func (t RegexToken) String() string { return t.lit }

func (t RegexToken) Type() int { return t.ttype }

func (t RegexToken) Location() Location { return t.location }

type RegexLexer struct {
	patterns   []Pattern
	ellide     map[int]struct{}
	transforms map[int]func(string) string
}

type Pattern struct {
	TokenType int
	Regex     string
}

type RegexLexerOption func(*RegexLexer)

func WithEllide(ellide ...int) RegexLexerOption {
	return func(l *RegexLexer) {
		for _, ttype := range ellide {
			l.ellide[ttype] = struct{}{}
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
		ellide:     map[int]struct{}{},
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
				if _, ok := l.ellide[pattern.TokenType]; !ok {
					lit := match
					if transform, ok := l.transforms[pattern.TokenType]; ok {
						lit = transform(match)
					}
					tokens = append(tokens, RegexToken{
						location: location,
						ttype:    pattern.TokenType,
						lit:      lit,
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
