package ast

type Type interface {
	typ()
}

type BuiltInType struct {
	Name string `@("string" | "bool" | "number" | "any")`
}

func (t *BuiltInType) typ() {}

type ObjectType struct {
	Fields []*ObjectTypeField `"{" @@ ("," @@)* "}"`
}

type ObjectTypeField struct {
	Name string `@Ident ":"`
	Type Type   `@@`
}

func (t *ObjectType) typ() {}

type CustomType struct {
	Name string `@Ident`
}

func (t *CustomType) typ() {}
