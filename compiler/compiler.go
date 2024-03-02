package compiler

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
//
// 	"config-lang/ast"
// 	"config-lang/value"
// )

// type def struct {
// 	name  string
// 	props []prop
// }
//
// type prop struct {
// 	name string
// 	typ  string
// }
//
// type Compiler struct {
// 	defs     []def
// 	bindings map[string]ast.Expr
// }
//
// func New() *Compiler {
// 	return &Compiler{
// 		bindings: map[string]ast.Expr{},
// 	}
// }
//
// func (t *Compiler) Compile(p *ast.Program) (string, error) {
// 	res := strings.Builder{}
//
// 	for _, s := range p.Statements {
// 		out, err := t.compileStmt(s)
// 		if err != nil {
// 			return "", err
// 		}
// 		_, err = res.WriteString(out)
// 		if err != nil {
// 			return "", err
// 		}
// 	}
//
// 	return res.String(), nil
// }
//
// func (t *Compiler) compileStmt(s ast.Stmt) (string, error) {
// 	switch s := s.(type) {
//
// 	case *ast.Def:
// 		props := make([]prop, 0, len(s.Props))
// 		for _, p := range s.Props {
// 			props = append(props, prop{
// 				name: p.Name,
// 				typ:  p.Type.Value,
// 			})
// 		}
// 		t.defs = append(t.defs, def{name: s.Name, props: props})
//
// 	case *ast.Output:
// 		expr, err := t.compileExpr(s.Expr)
// 		if err != nil {
// 			return "", err
// 		}
//
// 		outPath, err := filepath.Abs(fmt.Sprintf("./%s", s.Name))
// 		if err != nil {
// 			return "", err
// 		}
//
// 		f, err := os.Create(outPath)
// 		if err != nil {
// 			return "", err
// 		}
// 		defer f.Close()
//
// 		_, err = f.WriteString(expr.String())
// 		if err != nil {
// 			panic(err)
// 		}
//
// 	case *ast.Let:
// 		t.bindings[s.Name] = s.Expr
//
// 	default:
// 		panic(fmt.Sprintf("Unimplemented Statement: %T", s))
// 	}
//
// 	return "", nil
// }
//
// func (t *Compiler) compileExpr(s ast.Expr) (value.Value, error) {
// 	switch s := s.(type) {
//
// 	case *ast.Object:
// 		obj := map[string]value.Value{}
// 		for _, f := range s.Fields {
// 			if f.Expr == nil {
// 				f.Expr = &ast.Ident{Name: f.Name}
// 			}
// 			filedExpr, err := t.compileExpr(f.Expr)
// 			if err != nil {
// 				return value.Value{}, err
// 			}
//
// 			obj[f.Name] = filedExpr
// 		}
//
// 		return value.NewObject(obj), nil
//
// 	case *ast.String:
// 		return value.NewString(s.Value), nil
//
// 	case *ast.Number:
// 		return value.NewInt(s.Value), nil
//
// 	case *ast.Ident:
// 		expr := t.bindings[s.Name]
// 		return t.compileExpr(expr)
//
// 	default:
// 		panic(fmt.Sprintf("Unimplemented Expr: %T", s))
// 	}
// }
