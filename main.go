package main

import (
	"fmt"
	"time"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	src := `
		n := 5000000
		even_nums := 0
		odd_nums := 0

		i := 0
		while i < n {
			if i % 2 == 0 {
				even_nums = even_nums + 1
			}	else {
				odd_nums  = odd_nums + 1
			}

			i = i + 1
		}

		even_nums | echo()
		odd_nums  | echo()
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

	irc := ir.NewCompiler()
	if err != nil {
		panic(err)
	}

	ircr, err := irc.Compile(p)
	if err != nil {
		panic(err)
	}

	ct := time.Now()
	c := compiler.New(compiler.WithOptimization(true))
	mainFrame, constants, err := c.Compile(ircr)
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
