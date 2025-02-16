package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	src := `
		forEach := |arr, f| {
			for i := 0; i < len(arr); i++ {
				f(arr[i])
			}
		}	

		forEach([1, 2, 3], |x| {
			echo(x)
		})
	`

	// defer profile.Start(profile.MemProfile).Stop()

	// src := `
	// 	x := |i| {
	// 		if i == 0 {
	// 			return 0
	// 		}
	//
	// 		x(i - 1)
	// 	}
	//
	// 	x(10)
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

	irt := time.Now()
	ircr, err := irc.Compile(p)
	if err != nil {
		panic(err)
	}
	res := ""
	for _, s := range ircr.Statements {
		res += s.String(0) + "\n"
	}
	iro, err := os.Create("ir.js")
	if err != nil {
		panic(err)
	}
	defer iro.Close()
	iro.WriteString(res)
	fmt.Println("ir took: ", time.Since(irt))

	ct := time.Now()
	c := compiler.New(compiler.WithOptimization(true))
	instructions, vars, constants, err := c.Compile(ircr)
	if err != nil {
		panic(err)
	}
	fmt.Println("compiler took: ", time.Since(ct))

	fmt.Println(opcode.PrintOpcodes(instructions))

	for _, c := range constants {
		if c.VType == vm.ValueTypeFunction {
			fn := c.GetFunction()
			fmt.Println(opcode.PrintOpcodes(fn.Instructions))
		}
	}

	vt := time.Now()
	vm := vm.New(constants, instructions, vars)
	vm.Run()
	fmt.Println("vm took: ", time.Since(vt))
}
