package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"config-lang/ast"
	"config-lang/typechecker"
)

// TODO: Maybe call it Weaver
func main() {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Keyword", Pattern: `def|output|let`},
		{Name: "Comment", Pattern: `(?i)rem[^\n]*`},
		{Name: "String", Pattern: `"(\\"|[^"])*"`},
		{Name: "Number", Pattern: `[-+]?(\d*\.)?\d+`},
		{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{Name: "whitespace", Pattern: `[ \t\n]+`},
		{Name: "EOL", Pattern: `[\n\r]+`},
	})

	parseStart := time.Now()
	parser, err := participle.Build[ast.Program](
		participle.Lexer(lex),
		participle.Elide("whitespace"),
		participle.Unquote("String"),
		participle.Union[ast.Stmt](&ast.Def{}, &ast.Output{}, &ast.Let{}, &ast.Assign{}),
		participle.Union[ast.Expr](&ast.Object{}, &ast.String{}, &ast.Number{}, &ast.Ident{}),
	)
	if err != nil {
		panic(err)
	}

	p, err := parser.ParseString("main.tf", `
    let x = "hi"
    x = 1
	  `)
	if err != nil {
		panic(err)
	}
	parseDuration := time.Now().Sub(parseStart).Milliseconds()
	log.Printf("parsing took [%dms]", parseDuration)

	checker := typechecker.New()
	err = checker.Check(p)
	if err != nil {
		panic(err)
	}

	// c := compiler.New()
	// _, err = c.Compile(p)
	// if err != nil {
	// 	panic(err)
	// }

	_ = p

	fmt.Println(parser.String())
}
