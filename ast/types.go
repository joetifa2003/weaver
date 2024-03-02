package ast

import "fmt"

type Type interface {
	typ()
}

type StringType struct{}

func (t StringType) typ() {}

type IntType struct{}

func (t IntType) typ() {}

type FloatType struct{}

func (t FloatType) typ() {}

type ObjectType struct{}

func (t ObjectType) typ() {}

type CustomType struct{}

func (t CustomType) typ() {}

type AnyType struct{}

func (t AnyType) typ() {}

type TypeNode struct {
	Type Type
}

func (t *TypeNode) Capture(values []string) error {
	ident := values[0]

	switch ident {
	case "string":
		t.Type = StringType{}
	case "int":
		t.Type = IntType{}
	case "any":
		t.Type = AnyType{}
	default:
		return fmt.Errorf("unknown type %s", ident)
	}

	return nil
}
