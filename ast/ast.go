package ast

import "github.com/alecthomas/participle/v2/lexer"

type Program struct {
	Statements []Stmt `@@*`
}

type Stmt interface {
	stmt()
}

type Fn struct {
	Pos        lexer.Position
	Name       string `"fn" @Ident`
	Args       []*Arg `"(" (@@ ("," @@)* )? ")"`
	ReturnType Type   `":" @@`
	Statements []Stmt `"{" @@* "}"`
	EndPos     lexer.Position
}

func (t Fn) stmt() {}

type Arg struct {
	Name string `@Ident`
	Type Type   `":" @@`
}

// Def define type
type Def struct {
	Name string `"def" @Ident`
	Type Type   `@@`
}

func (t *Def) stmt() {}

type Prop struct {
	Name string `@Ident ":"`
	Type Type   `@@`
}

type Output struct {
	Name string `"output" @String`
	Expr *Expr  `@@`
}

func (t *Output) stmt() {}

type Let struct {
	Name string `"let" @Ident `
	Type Type   `[":" @@]?`
	Expr *Expr  `"=" @@`
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

type Return struct {
	Expr *Expr `"return" @@`
}

func (t *Return) stmt() {}

type Expr struct {
	Equality *Equality `@@`
}

type Equality struct {
	Pos    lexer.Position
	Left   *Comparison `@@`
	Op     string      `( @( "!" "=" | "=" "=" )`
	Right  *Equality   `  @@ )*`
	EndPos lexer.Position
}

type Comparison struct {
	Pos    lexer.Position
	Left   *Addition   `@@`
	Op     string      `( @( ">" | ">" "=" | "<" | "<" "=" )`
	Right  *Comparison `  @@ )*`
	EndPos lexer.Position
}

type Addition struct {
	Pos    lexer.Position
	Left   *Multiplication `@@`
	Op     string          `( @( "-" | "+" )`
	Right  *Addition       `  @@ )*`
	EndPos lexer.Position
}

type Multiplication struct {
	Pos    lexer.Position
	Left   *Unary          `@@`
	Op     string          `( @( "/" | "*" )`
	Right  *Multiplication `  @@ )*`
	EndPos lexer.Position
}

type Unary struct {
	Pos    lexer.Position
	Op     string `  ( @( "!" | "-" )`
	Unary  *Unary `    @@ )`
	Atom   Atom   `| @@`
	EndPos lexer.Position
}

type Atom interface {
	atom()
}

type Paren struct {
	Pos    lexer.Position
	Expr   *Expr `  "(" @@ ")"`
	EndPos lexer.Position
}

func (t Paren) atom() {}

type Object struct {
	Pos lexer.Position

	Fields []*Field `"{" (@@ ("," @@)* )? "}"`

	EndPos lexer.Position
}

func (t *Object) atom() {}

type Field struct {
	Pos lexer.Position

	Name string `@Ident [":"`
	Expr *Expr  `@@]?`

	EndPos lexer.Position
}

type String struct {
	Pos lexer.Position

	Value string `@String`

	EndPos lexer.Position
}

func (t *String) atom() {}

type Number struct {
	Pos lexer.Position

	Value int `@Number`

	EndPos lexer.Position
}

func (t *Number) atom() {}

type Ident struct {
	Pos lexer.Position

	Name string "@Ident"

	EndPos lexer.Position
}

func (t *Ident) atom() {}

type Call struct {
	Pos lexer.Position

	Name *Ident  "@@"
	Args []*Expr `"(" (@@ ("," @@)* )? ")"`

	EndPos lexer.Position
}

func (t *Call) atom() {}

type Bool struct {
	Pos lexer.Position

	Value _Bool `@("true" | "false")`

	EndPos lexer.Position
}

type _Bool bool

func (t *_Bool) Capture(values []string) {
	if values[0] == "true" {
		*t = true
	} else if values[0] == "false" {
		*t = false
	}
}

func (t *Bool) atom() {}
