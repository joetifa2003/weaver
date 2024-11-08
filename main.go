package main

import (
	"fmt"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	src := `
	x := 0
	even := 0

	while x < 100 {
		if x % 2 == 0 {
			even = even + 1
		}	

		x = x + 1
	}

	echo x 
	echo even
	`

	p, err := parser.Parse(src)
	if err != nil {
		panic(err)
	}

	c := compiler.New()
	mainFrame, constants, err := c.Compile(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(opcode.PrintOpcodes(mainFrame.Instructions))

	vm := vm.New(constants, mainFrame)
	vm.Run()
}
