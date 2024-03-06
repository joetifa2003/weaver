package typechecker

import "github.com/alecthomas/participle/v2/lexer"

type Type interface {
	typ()
	Pos() lexer.Position
	EndPos() lexer.Position
	Is(other Type) bool
}

type StringType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t StringType) typ() {}

func (t StringType) Pos() lexer.Position { return t.pos }

func (t StringType) EndPos() lexer.Position { return t.endPos }

func (t StringType) Is(other Type) bool { return isType[StringType](other) }

type NumberType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t NumberType) typ() {}

func (t NumberType) Pos() lexer.Position { return t.pos }

func (t NumberType) EndPos() lexer.Position { return t.endPos }

func (t NumberType) Is(other Type) bool { return isType[NumberType](other) }

type BoolType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t BoolType) typ() {}

func (t BoolType) Pos() lexer.Position { return t.pos }

func (t BoolType) EndPos() lexer.Position { return t.endPos }

func (t BoolType) Is(other Type) bool { return isType[BoolType](other) }

type AnyType struct {
	pos    lexer.Position
	endPos lexer.Position
}

func (t AnyType) typ() {}

func (t AnyType) Pos() lexer.Position { return t.pos }

func (t AnyType) EndPos() lexer.Position { return t.endPos }

func (t AnyType) Is(other Type) bool { return true }

type ObjectType struct {
	pos    lexer.Position
	Fields map[string]Type
	endPos lexer.Position
}

func (t ObjectType) typ() {}

func (t ObjectType) Pos() lexer.Position { return t.pos }

func (t ObjectType) EndPos() lexer.Position { return t.endPos }

func (t ObjectType) Is(other Type) bool {
	if !isType[ObjectType](other) {
		return false
	}

	otherObj := other.(ObjectType)

	for k, v := range otherObj.Fields {
		typ, ok := t.Fields[k]
		if !ok {
			continue
		}

		if !typ.Is(v) {
			return false
		}
	}

	return true
}

func isType[T Type](t Type) bool {
	_, ok := t.(T)
	return ok
}