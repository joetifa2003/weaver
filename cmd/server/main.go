package main

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joetifa2003/weaver/builtin"
	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/", func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					c.String(200, r.Error())
				case string:
					c.String(200, r)
				default:
					c.String(200, "error")
				}
			}
		}()

		srcData, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(200, err.Error())
		}

		p, err := parser.Parse(string(srcData))
		if err != nil {
			return c.String(200, err.Error())
		}

		irc := ir.NewCompiler()

		ircr, err := irc.Compile(p)
		if err != nil {
			return c.String(200, err.Error())
		}

		irout, err := os.Create("irout.wvr")
		defer irout.Close()
		_, err = irout.WriteString(ircr.String())
		if err != nil {
			return c.String(200, err.Error())
		}

		compiler := compiler.New(builtin.StdReg, compiler.WithOptimization(true))
		instructions, vars, constants, err := compiler.Compile(ircr)
		if err != nil {
			return c.String(200, err.Error())
		}

		executor := vm.NewExecutor(constants)
		val := executor.Run(
			vm.Frame{
				Instructions: instructions,
				NumVars:      vars,
				HaltAfter:    true,
			},
			0,
		)

		if val.IsError() {
			return c.String(200, val.GetError().Error())
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
