package ast

type Program struct {
	Statements []Stmt `@@*`
}

type Stmt interface {
	stmt()
}

// Def define type
type Def struct {
	Name  string  `"def" @Ident "{"`
	Props []*Prop `(@@ ("," @@)* )? "}"`
}

func (t *Def) stmt() {}

type Prop struct {
	Name string   `@Ident ":"`
	Type TypeNode `@Ident`
}

type Output struct {
	Name string `"output" @String`
	Expr Expr   `@@`
}

func (t *Output) stmt() {}

type Let struct {
	Name     string   `"let" @Ident `
	TypeNode TypeNode `[":" @Ident]?`
	Expr     Expr     `"=" @@`
}

func (t *Let) stmt() {}

type Assign struct {
	Name string "@Ident"
	Expr Expr   `"=" @@`
}

func (t *Assign) stmt() {}

type Expr interface {
	expr()
}

type Object struct {
	Type   TypeNode `@Ident`
	Fields []*Field `"{" (@@ ("," @@)* )? "}"`
}

func (t *Object) expr() {}

type Field struct {
	Name string `@Ident [":"`
	Expr Expr   `@@]?`
}

type String struct {
	Value string `@String`
}

func (t *String) expr() {}

// TODO: implement int and float
type Number struct {
	Value int `@Number`
}

func (t *Number) expr() {}

type Ident struct {
	Name string "@Ident"
}

func (t *Ident) expr() {}
