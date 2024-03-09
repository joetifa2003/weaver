package typechecker

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type TypeError struct {
	ExpectedType Type
	ActualType   Type
	Src          string
	Pos          lexer.Position
}

func (e *TypeError) Error() string {
	line := strings.Split(e.Src, "\n")[e.Pos.Line-1]
	sep := strings.Repeat("-", e.Pos.Column-1) + "^"
	return fmt.Sprintf("type error at %s: expected %s, got %s\n\n%s\n%s", e.Pos.String(), e.ExpectedType, e.ActualType, line, sep)
}

func NewTypeError(src string, expectedType Type, actualType Type, pos lexer.Position) *TypeError {
	return &TypeError{
		Src:          src,
		ExpectedType: expectedType,
		ActualType:   actualType,
		Pos:          pos,
	}
}
