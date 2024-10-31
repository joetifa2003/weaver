package main

import (
	"fmt"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
)

// TODO: Maybe call it Weaver
func main() {
	src := `
	x := 1
	echo x
	`

	p, err := parser.Parse(src)
	if err != nil {
		panic(err)
	}

	c := compiler.New()
	err = c.Compile(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Frames[0].Instructions)
	fmt.Println(opcode.PrintOpcodes(c.Frames[0].Instructions))
}
