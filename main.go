package main

import (
	"fmt"
	"time"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/value"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	src := `
	isEven := |n| n % 2 == 0	
	add := |a, b| a + b

	a := 1 | add(1) | isEven()
	echo a
	`

	pt := time.Now()
	p, err := parser.Parse(src)
	if err != nil {
		panic(err)
	}
	fmt.Println("parser took: ", time.Since(pt))

	ct := time.Now()
	c := compiler.New()
	mainFrame, constants, err := c.Compile(p)
	if err != nil {
		panic(err)
	}
	fmt.Println("compiler took: ", time.Since(ct))

	fmt.Println(opcode.PrintOpcodes(mainFrame.Instructions))

	for _, c := range constants {
		if c.VType == value.ValueTypeFunction {
			fn := c.GetFunction()
			fmt.Println(opcode.PrintOpcodes(fn.Instructions))
		}
	}

	vt := time.Now()
	vm := vm.New(constants, mainFrame)
	vm.Run()
	fmt.Println("vm took: ", time.Since(vt))
}
