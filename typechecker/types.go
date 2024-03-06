package typechecker

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Type interface {
	typ()
	Pos() lexer.Position
	EndPos() lexer.Position
	IsAssignableTo(other Type) bool
	String() string
}

type StringType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t StringType) typ() {}

func (t StringType) Pos() lexer.Position { return t.pos }

func (t StringType) EndPos() lexer.Position { return t.endPos }

func (t StringType) IsAssignableTo(other Type) bool { return isType[StringType](other) }

func (t StringType) String() string { return "string" }

type NumberType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t NumberType) typ() {}

func (t NumberType) Pos() lexer.Position { return t.pos }

func (t NumberType) EndPos() lexer.Position { return t.endPos }

func (t NumberType) IsAssignableTo(other Type) bool { return isType[NumberType](other) }

func (t NumberType) String() string { return "number" }

type BoolType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t BoolType) typ() {}

func (t BoolType) Pos() lexer.Position { return t.pos }

func (t BoolType) EndPos() lexer.Position { return t.endPos }

func (t BoolType) IsAssignableTo(other Type) bool { return isType[BoolType](other) }

func (t BoolType) String() string { return "bool" }

type AnyType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t AnyType) typ() {}

func (t AnyType) Pos() lexer.Position { return t.pos }

func (t AnyType) EndPos() lexer.Position { return t.endPos }

func (t AnyType) IsAssignableTo(other Type) bool { return true }

func (t AnyType) String() string { return "any" }

type ObjectType struct {
	pos    lexer.Position
	Fields map[string]Type
	endPos lexer.Position
}

func (t ObjectType) typ() {}

func (t ObjectType) Pos() lexer.Position { return t.pos }

func (t ObjectType) EndPos() lexer.Position { return t.endPos }

func (t ObjectType) String() string {
	var res strings.Builder

	res.WriteString("{ ")

	idx := 0
	for name, typ := range t.Fields {
		res.WriteString(fmt.Sprintf("%s: %s", name, typ))

		if idx != len(t.Fields)-1 {
			res.WriteString(", ")
		}

		idx++
	}
	res.WriteString(" }")

	return res.String()
}

func (t ObjectType) IsAssignableTo(other Type) bool {
	if !isType[ObjectType](other) {
		return false
	}

	otherObj := other.(ObjectType)

	for k, v := range otherObj.Fields {
		typ, ok := t.Fields[k]
		if !ok {
			return false
		}

		if !typ.IsAssignableTo(v) {
			return false
		}
	}

	return true
}

func isType[T Type](t Type) bool {
	_, ok := t.(T)
	return ok
}
