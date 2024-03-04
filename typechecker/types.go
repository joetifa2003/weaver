package typechecker

type Type interface {
	typ()
}

type StringType struct{}

func (t StringType) typ() {}

type NumberType struct{}

func (t NumberType) typ() {}

type BoolType struct{}

func (t BoolType) typ() {}

type AnyType struct{}

func (t AnyType) typ() {}
