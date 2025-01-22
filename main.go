package main

import (
	"fmt"
	"time"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	src := `
		x := [1, 2, 3]
		x[0] = 2
		echo(x[0])
	`

	// src := `
	// for i := 0; i < 100; i = i + 1 {
	// 	[1, 2, 3]
	// 		| map(|x| x + 1)
	// 		| filter(|x| x % 2 == 0)
	// }
	// `

	// src := `
	// n := 10000000
	// even_nums := 0
	// odd_nums := 0
	//
	// is_even := |x| x % 2 == 0
	//
	// for i := 0; i < n; i = i + 1 {
	// 	if is_even(i) {
	// 		even_nums = even_nums + 1
	// 	} else {
	// 		odd_nums = odd_nums + 1
	// 	}
	// }
	//
	// echo even_nums
	// echo odd_nums
	// `

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
		if c.VType == vm.ValueTypeFunction {
			fn := c.GetFunction()
			fmt.Println(opcode.PrintOpcodes(fn.Instructions))
		}
	}

	vt := time.Now()
	vm := vm.New(constants, mainFrame.Instructions, len(mainFrame.Vars))
	vm.Run()
	fmt.Println("vm took: ", time.Since(vt))
}
