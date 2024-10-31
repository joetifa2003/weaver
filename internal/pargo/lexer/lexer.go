package lexer

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

type SimpleToken struct {
	location Location
	ttype    SimpleTokenType
	lit      string
}

func (t SimpleToken) String() string { return t.lit }

func (t SimpleToken) Type() int { return int(t.ttype) }

func (t SimpleToken) Location() Location { return t.location }

type SimpleTokenType int

const (
	TT_CHARACTER SimpleTokenType = iota
)

type SimpleLexer struct {
}

func New() *SimpleLexer {
	return &SimpleLexer{}
}

func (l *SimpleLexer) Lex(input string) ([]Token, error) {
	runes := []rune(input)

	res := []Token{}

	i := 0
	column := 0
	line := 1

	for i < len(runes) {
		if runes[i] == ' ' || runes[i] == '\r' || runes[i] == '\n' || runes[i] == '\t' {
			if runes[i] == '\n' {
				line++
				column = 0
			}
			i++
			continue
		}

		res = append(res, SimpleToken{
			ttype: TT_CHARACTER,
			lit:   string(runes[i]),
			location: Location{
				Line:   line,
				Column: column + 1,
			},
		})

		i++
		column++
	}

	return res, nil
}

func isCharacter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\r'
}

func isSymbol(r rune) bool {
	return r == '='
}
