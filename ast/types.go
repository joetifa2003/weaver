package ast

type Type interface {
	typ()
}

type BuiltInType struct {
	Name string `@("string" | "bool" | "number" | "any")`
}

func (t *BuiltInType) typ() {}

type ObjectType struct {
	Fields []*ObjectTypeField `@@ ("," @@)*`
}

type ObjectTypeField struct {
	Name string `@Ident [":"`
	Expr Type   `@@]?`
}

func (t *ObjectType) typ() {}
