package jit

import (
	"bytes"
	"fmt"

	"github.com/joetifa2003/weaver/ast"
)

func Generate(program *ast.Program) (string, error) {
	var out bytes.Buffer

	fmt.Fprintln(&out, "package main")
	fmt.Fprintln(&out, "")
	fmt.Fprintln(&out, `import "fmt"`)
	fmt.Fprintln(&out, "")
	fmt.Fprintln(&out, "type Value = interface{}")
	fmt.Fprintln(&out, "type Object = map[string]Value")
	fmt.Fprintln(&out, "type Function = func(...Value) Value")
	fmt.Fprintln(&out, "")
	fmt.Fprintln(&out, `
func echo(args ...Value) Value {
    fmt.Println(args[0])
    return nil
}
`)
	fmt.Fprintln(&out, "func main() {")
	g := &generator{out: &out}
	for _, stmt := range program.Statements {
		if err := g.genStmt(stmt); err != nil {
			return "", err
		}
	}
	fmt.Fprintln(&out, "}")

	// For now, we are not formatting the output, but it's a good idea to do it
	// to make the generated code readable.
	// formatted, err := format.Source(out.Bytes())
	// if err != nil {
	// 	return "", err
	// }

	return out.String(), nil
}

type generator struct {
	out *bytes.Buffer
}

func (g *generator) genStmt(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case ast.LetStmt:
		return g.genLetStmt(s)
	case ast.ExprStmt:
		return g.genExprStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", s)
	}
}

func (g *generator) genLetStmt(stmt ast.LetStmt) error {
	fmt.Fprintf(g.out, "var %s Value = ", stmt.Name)
	if err := g.genExpr(stmt.Expr); err != nil {
		return err
	}
	fmt.Fprintln(g.out, "")
	return nil
}

func (g *generator) genExprStmt(stmt ast.ExprStmt) error {
	if err := g.genExpr(stmt.Expr); err != nil {
		return err
	}
	fmt.Fprintln(g.out, "")
	return nil
}

func (g *generator) genExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case ast.LambdaExpr:
		return g.genLambdaExpr(e)
	case ast.PostFixExpr:
		return g.genPostfixExpr(e)
	case ast.IdentExpr:
		return g.genIdentExpr(e)
	case ast.ObjectExpr:
		return g.genObjectExpr(e)
	case ast.StringExpr:
		return g.genStringExpr(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", e)
	}
}

func (g *generator) genLambdaExpr(expr ast.LambdaExpr) error {
	fmt.Fprintf(g.out, "Function(func(args ...Value) Value {\n")
	for i, param := range expr.Params {
		fmt.Fprintf(g.out, "%s := args[%d]\n", param, i)
	}
	fmt.Fprintf(g.out, "return ")
	if err := g.genExpr(expr.Expr); err != nil {
		return err
	}
	fmt.Fprintf(g.out, "})")
	return nil
}

func (g *generator) genPostfixExpr(expr ast.PostFixExpr) error {
	if err := g.genExpr(expr.Expr); err != nil {
		return err
	}

	for _, op := range expr.Ops {
		switch o := op.(type) {
		case ast.CallOp:
			// A bit of a hack to handle the built-in echo function.
			// A proper implementation would have a symbol table to distinguish
			// between built-in functions and user-defined functions.
			isBuiltin := false
			if ident, ok := expr.Expr.(ast.IdentExpr); ok {
				if ident.Name == "echo" {
					isBuiltin = true
				}
			}

			if !isBuiltin {
				fmt.Fprintf(g.out, ".(Function)")
			}

			fmt.Fprintf(g.out, "(")
			for i, arg := range o.Args {
				if err := g.genExpr(arg); err != nil {
					return err
				}
				if i < len(o.Args)-1 {
					fmt.Fprintf(g.out, ", ")
				}
			}
			fmt.Fprintf(g.out, ")")
		case ast.DotOp:
			fmt.Fprintf(g.out, ".(Object)[\"%s\"]", o.Index)
		default:
			return fmt.Errorf("unsupported postfix operator: %T", o)
		}
	}

	return nil
}

func (g *generator) genIdentExpr(expr ast.IdentExpr) error {
	fmt.Fprintf(g.out, "%s", expr.Name)
	return nil
}

func (g *generator) genObjectExpr(expr ast.ObjectExpr) error {
	fmt.Fprintf(g.out, "Object{")
	var keys []string
	for k := range expr.KVs {
		keys = append(keys, k)
	}
	for i, k := range keys {
		v := expr.KVs[k]
		fmt.Fprintf(g.out, "\"%s\": ", k)
		if err := g.genExpr(v); err != nil {
			return err
		}
		if i < len(keys)-1 {
			fmt.Fprintf(g.out, ", ")
		}
	}
	fmt.Fprintf(g.out, "}")
	return nil
}

func (g *generator) genStringExpr(expr ast.StringExpr) error {
	fmt.Fprintf(g.out, "\"%s\"", expr.Value)
	return nil
}
