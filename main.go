package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"config-lang/ast"
	"config-lang/typechecker"
)

// TODO: Maybe call it Weaver
func main() {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Keyword", Pattern: `def|output|let|true|false|string|number|any|bool|fn`},
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
		participle.Union[ast.Stmt](&ast.Def{}, &ast.Output{}, &ast.Let{}, &ast.Assign{}, &ast.Block{}, &ast.If{}, &ast.Fn{}, &ast.Return{}),
		participle.Union[ast.Atom](&ast.Paren{}, &ast.Object{}, &ast.String{}, &ast.Number{}, &ast.Bool{}, &ast.Call{}, &ast.Ident{}),
		participle.Union[ast.TypeAtom](&ast.BuiltInType{}, &ast.ObjectType{}, &ast.CustomType{}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(parser.String())

	src := `
  fn something(a: number, b: string): number | string {

  }

  let x: number | string | bool = something(1, "2") 
  `

	p, err := parser.ParseString("main.tf", src)
	if err != nil {
		panic(err)
	}
	parseDuration := time.Now().Sub(parseStart).Milliseconds()
	log.Printf("parsing took [%dms]", parseDuration)

	parserOut, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("ast.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(parserOut)
	if err != nil {
		panic(err)
	}

	checker := typechecker.New(src)
	err = checker.Check(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// c := compiler.New()
	// _, err = c.Compile(p)
	// if err != nil {
	// 	panic(err)
	// }

}
