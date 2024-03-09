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
	Nullable() bool
}

type StringType struct {
	pos      lexer.Position
	endPos   lexer.Position
	nullable bool
}

func (t StringType) typ() {}

func (t StringType) Pos() lexer.Position { return t.pos }

func (t StringType) EndPos() lexer.Position { return t.endPos }

func (t StringType) IsAssignableTo(other Type) bool { return isType[StringType](other) }

func (t StringType) String() string { return "string" }

func (t StringType) Nullable() bool { return t.nullable }

type NumberType struct {
	pos      lexer.Position
	endPos   lexer.Position
	nullable bool
}

func (t NumberType) typ() {}

func (t NumberType) Pos() lexer.Position { return t.pos }

func (t NumberType) EndPos() lexer.Position { return t.endPos }

func (t NumberType) IsAssignableTo(other Type) bool { return isType[NumberType](other) }

func (t NumberType) String() string { return "number" }

func (t NumberType) Nullable() bool { return t.nullable }

type BoolType struct {
	pos      lexer.Position
	endPos   lexer.Position
	nullable bool
}

func (t BoolType) typ() {}

func (t BoolType) Pos() lexer.Position { return t.pos }

func (t BoolType) EndPos() lexer.Position { return t.endPos }

func (t BoolType) IsAssignableTo(other Type) bool { return isType[BoolType](other) }

func (t BoolType) String() string { return "bool" }

func (t BoolType) Nullable() bool { return t.nullable }

type AnyType struct {
	pos      lexer.Position
	endPos   lexer.Position
	nullable bool
}

func (t AnyType) typ() {}

func (t AnyType) Pos() lexer.Position { return t.pos }

func (t AnyType) EndPos() lexer.Position { return t.endPos }

func (t AnyType) IsAssignableTo(other Type) bool { return true }

func (t AnyType) String() string { return "any" }

func (t AnyType) Nullable() bool { return t.nullable }

type ObjectType struct {
	pos      lexer.Position
	Fields   map[string]Type
	endPos   lexer.Position
	nullable bool
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

func (t ObjectType) Nullable() bool { return t.nullable }

func (t ObjectType) IsAssignableTo(other Type) bool {
	if !isType[ObjectType](other) {
		return false
	}

	otherObj := other.(ObjectType)

	for k, v := range otherObj.Fields {
		typ, ok := t.Fields[k]
		if !ok {
			if v.Nullable() {
				return true
			}
			return false
		}

		if !typ.IsAssignableTo(v) {
			return false
		}
	}

	return true
}

type FnType struct {
	Args       []Type
	ReturnType Type

	// TODO: Add base type that has all of that plus the methods
	pos      lexer.Position
	endPos   lexer.Position
	nullable bool
}

func (t FnType) typ() {}

func (t FnType) Pos() lexer.Position { return t.pos }

func (t FnType) EndPos() lexer.Position { return t.endPos }

func (t FnType) String() string {
	// TODO: proper Fn type string repr
	return "<fn>"
}

func (t FnType) Nullable() bool { return t.nullable }

func (t FnType) IsAssignableTo(other Type) bool {
	if !isType[FnType](other) {
		return false
	}

	otherObj := other.(FnType)

	if len(t.Args) != len(otherObj.Args) {
		return false
	}

	if !t.ReturnType.IsAssignableTo(otherObj.ReturnType) {
		return false
	}

	for i, arg := range t.Args {
		if !arg.IsAssignableTo(otherObj.Args[i]) {
			return false
		}
	}

	return true
}

func isType[T Type](t Type) bool {
	_, ok := t.(T)
	return ok
}
