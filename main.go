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
	arr := [[0, 1, 3]]

	obj := {
		"a": "Hello world",
		"b": arr
	}

	len(arr) 		|> echo()
	len(arr[0]) |> echo()
	len(obj) 		|> echo()

	nums := [1, 2, 3]
	`

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
