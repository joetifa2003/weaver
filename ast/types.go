package ast

type TypeAtom interface {
	typ()
}

type BuiltInType struct {
	Name     string `@("string" | "bool" | "number" | "any")`
	Nullable bool   `@("?")?`
}

func (t *BuiltInType) typ() {}

type ObjectType struct {
	Fields   []*ObjectTypeField `"{" @@ ("," @@)* "}"`
	Nullable bool               `@("?")?`
}

type ObjectTypeField struct {
	Name string `@Ident ":"`
	Type *Type  `@@`
}

func (t *ObjectType) typ() {}

type CustomType struct {
	Name     string `@Ident`
	Nullable bool   `@("?")?`
}

func (t *CustomType) typ() {}

type Type struct {
	Variants []TypeAtom `@@ ("|" @@)*`
}
