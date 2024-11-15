package main

import (
	"fmt"
	"time"

	"github.com/pkg/profile"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	defer profile.Start().Stop()
	src := `
	isEven := |n| n % 2 == 0	
	evenNums := 0	
	n := 10000000

	i := 0
	while i < n {
		if isEven(i) {
			evenNums = evenNums + 1
		}
		
		i = i + 1
	}
	
	echo evenNums
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

	// fmt.Println(opcode.PrintOpcodes(mainFrame.Instructions))
	//
	// for _, c := range constants {
	// 	if c.VType == value.ValueTypeFunction {
	// 		fn := c.GetFunction()
	// 		fmt.Println(opcode.PrintOpcodes(fn.Instructions))
	// 	}
	// }

	vt := time.Now()
	vm := vm.New(constants, mainFrame)
	vm.Run()
	fmt.Println("vm took: ", time.Since(vt))
}
