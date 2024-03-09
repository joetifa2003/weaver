package typechecker

import (
	"fmt"
	"strings"
)

type TypeError struct {
	ExpectedType Type
	ActualType   Type
	Src          string
}

func (e *TypeError) Error() string {
	line := strings.Split(e.Src, "\n")[e.ActualType.Pos().start.Line-1]
	sep := strings.Repeat("-", e.ActualType.Pos().start.Column-1) + "^"
	return fmt.Sprintf("type error at %s: expected %s, got %s\n\n%s\n%s", e.ActualType.Pos().start.String(), e.ExpectedType, e.ActualType, line, sep)
}

func NewTypeError(src string, expectedType Type, actualType Type) *TypeError {
	return &TypeError{
		Src:          src,
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}
