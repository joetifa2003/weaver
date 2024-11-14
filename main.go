package main

import (
	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	// defer profile.Start().Stop()

	src := `
	evenNums := 0
	x := 0
	while x < 10000000 {
		if x % 2 == 0 {
			evenNums = evenNums + 1
		}

		x = x + 1
	}

	echo evenNums
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

	// fmt.Println(opcode.PrintOpcodes(mainFrame.Instructions))
	//
	// for _, c := range constants {
	// 	if c.VType == value.ValueTypeFunction {
	// 		fn := c.GetFunction()
	// 		fmt.Println(opcode.PrintOpcodes(fn.Instructions))
	// 	}
	// }
	//
	vm := vm.New(constants, mainFrame)
	vm.Run()
}
