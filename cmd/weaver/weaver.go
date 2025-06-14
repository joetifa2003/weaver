package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/joetifa2003/weaver/builtin"
	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
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

					absPath, err := filepath.Abs(cc.Args().Get(0))
					if err != nil {
						return err
					}

					src := string(srcData)
					if len(src) == 0 {
						return errors.New("empty file")
					}

					p, err := parser.Parse(src)
					if err != nil {
						return err
					}

					irc := ir.NewCompiler()

					ircr, err := irc.Compile(absPath, p)
					if err != nil {
						return err
					}
					f, err := os.Create("ir.wvr")
					if err != nil {
						return err
					}
					f.WriteString(ircr.String())

					c := compiler.New(builtin.StdReg)
					instructions, vars, constants, err := c.Compile(ircr)
					if err != nil {
						return err
					}

					executor := vm.NewExecutor(builtin.StdReg)

					v := vm.New(executor)

					val, _ := v.RunFunction(
						vm.NewFunction(vm.FunctionValue{
							NumVars:      vars,
							Instructions: instructions,
							Constants:    constants,
							Path:         absPath,
						}),
					)
					if val.VType == vm.ValueTypeError {
						fmt.Println(val.String())
					}

					return nil
				},
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
