package typechecker

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Type interface {
	typ()
	Pos() TypePos
	SetPos(TypePos)
	Nullable() bool

	IsAssignableTo(other Type) bool
	String() string
}

type TypePos struct {
	start lexer.Position
	end   lexer.Position
}

type BaseType struct {
	pos      TypePos
	nullable bool
}

func NewBase(start lexer.Position, end lexer.Position) BaseType {
	return BaseType{
		pos: TypePos{
			start: start,
			end:   end,
		},
	}
}

func (t BaseType) typ() {}

func (t BaseType) Nullable() bool { return t.nullable }

func (t BaseType) Pos() TypePos { return t.pos }

func (t BaseType) SetPos(pos TypePos) { t.pos = pos }

type StringType struct{ BaseType }

func (t StringType) IsAssignableTo(other Type) bool { return isType[StringType](other) }

func (t StringType) String() string { return "string" }

type NumberType struct{ BaseType }

func (t NumberType) IsAssignableTo(other Type) bool { return isType[NumberType](other) }

func (t NumberType) String() string { return "number" }

type BoolType struct{ BaseType }

func (t BoolType) IsAssignableTo(other Type) bool { return isType[BoolType](other) }

func (t BoolType) String() string { return "bool" }

type AnyType struct{ BaseType }

func (t AnyType) IsAssignableTo(other Type) bool { return true }

func (t AnyType) String() string { return "any" }

type ObjectType struct {
	BaseType
	Fields map[string]Type
}

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
	BaseType

	Args       []Type
	ReturnType Type
}

func (t FnType) String() string {
	// TODO: proper Fn type string repr
	return "<fn>"
}

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
