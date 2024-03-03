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
	Expr *Expr  `@@`
}

func (t *Output) stmt() {}

type Let struct {
	Name     string   `"let" @Ident `
	TypeNode TypeNode `[":" @Ident]?`
	Expr     *Expr    `"=" @@`
}

func (t *Let) stmt() {}

type Assign struct {
	Name string "@Ident"
	Expr *Expr  `"=" @@`
}

func (t *Assign) stmt() {}

type Block struct {
	Statements []Stmt `"{" @@* "}"`
}

func (t *Block) stmt() {}

type If struct {
	Expr      *Expr `"if" "(" @@ ")"`
	Statement Stmt  `@@`
}

func (t *If) stmt() {}

type Expr struct {
	Equality *Equality `@@`
}

type Equality struct {
	Left  *Comparison `@@`
	Op    string      `( @( "!" "=" | "=" "=" )`
	Right *Equality   `  @@ )*`
}

type Comparison struct {
	Left  *Addition   `@@`
	Op    string      `( @( ">" | ">" "=" | "<" | "<" "=" )`
	Right *Comparison `  @@ )*`
}

type Addition struct {
	Left  *Multiplication `@@`
	Op    string          `( @( "-" | "+" )`
	Right *Addition       `  @@ )*`
}

type Multiplication struct {
	Left  *Unary          `@@`
	Op    string          `( @( "/" | "*" )`
	Right *Multiplication `  @@ )*`
}

type Unary struct {
	Op    string `  ( @( "!" | "-" )`
	Unary *Unary `    @@ )`
	Atom  Atom   `| @@`
}

type Atom interface {
	atom()
}

type Object struct {
	Type   TypeNode `@Ident`
	Fields []*Field `"{" (@@ ("," @@)* )? "}"`
}

func (t *Object) atom() {}

type Field struct {
	Name string `@Ident [":"`
	Expr *Expr  `@@]?`
}

type String struct {
	Value string `@String`
}

func (t *String) atom() {}

// TODO: implement int and float
type Number struct {
	Value int `@Number`
}

func (t *Number) atom() {}

type Ident struct {
	Name string "@Ident"
}

func (t *Ident) atom() {}

type _Bool bool

func (t *_Bool) Capture(values []string) {
	if values[0] == "true" {
		*t = true
	} else if values[0] == "false" {
		*t = false
	}
}

type Bool struct {
	Value _Bool `@("true" | "false")`
}

func (t *Bool) atom() {}
