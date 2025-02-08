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
		x := 1
		x == 1 |> assert()
	`

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
	// res := ""
	// for _, s := range ircr {
	// 	res += s.String(0) + "\n"
	// }
	// iro, err := os.Create("ir.js")
	// if err != nil {
	// 	panic(err)
	// }
	// defer iro.Close()
	// iro.WriteString(res)
	fmt.Println("ir took: ", time.Since(irt))

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
