package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/opcode"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	cmd := cli.Command{
		Name:  "weaver",
		Usage: "programming language",
		Commands: []*cli.Command{
			{
				Name:        "run",
				Usage:       "run a file",
				Description: "run [file]",
				Action: func(ctx context.Context, cc *cli.Command) error {
					srcData, err := os.ReadFile(cc.Args().Get(0))
					if err != nil {
						return err
					}

					src := string(srcData)

					p, err := parser.Parse(src)
					if err != nil {
						return err
					}

					irc := ir.NewCompiler()
					if err != nil {
						return err
					}

					ircr, err := irc.Compile(p)
					if err != nil {
						return err
					}

					c := compiler.New(compiler.WithOptimization(true))
					instructions, vars, constants, err := c.Compile(ircr)
					if err != nil {
						return err
					}

					fmt.Println(opcode.PrintOpcodes(instructions))
					for _, c := range constants {
						if c.VType == vm.ValueTypeFunction {
							fn := c.GetFunction()
							fmt.Println(opcode.PrintOpcodes(fn.Instructions))
						}
					}

					vm := vm.New(constants, instructions, vars)
					vm.Run()

					return nil
				},
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}

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

}
