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
		students := []

		for i := 0; i < 1000000; i = i + 1 {
			students
				|> push({name: string(i), age: 10})	
			students 
				|> push({name: string(i), age: 20})	
			students 
				|> push({name: string(i), age: 30})	
		}

		valid_names := []
		invalid_count := 0

		for i := 0; i < len(students); i = i + 1 {
			match students[i] {
				{name: n, age: a} if a >= 10 && a <= 20 => {
					valid_names 
						|> push(a)
				},
				else => {
					invalid_count = invalid_count + 1
				}
			}
		}

		valid_names 
			|> len() 
			|> echo()
		invalid_count 
			|> echo()
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
	res := ""
	for _, s := range ircr {
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
