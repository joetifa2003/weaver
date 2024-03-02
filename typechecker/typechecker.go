package typechecker

import (
	"fmt"

	"config-lang/ast"
)

type TypeChecker struct {
	defs     map[string]*ast.Def
	bindings map[string]ast.Type
}

func New() *TypeChecker {
	return &TypeChecker{
		defs:     map[string]*ast.Def{},
		bindings: map[string]ast.Type{},
	}
}

func (t *TypeChecker) Check(p *ast.Program) error {
	for _, stmt := range p.Statements {
		err := t.checkStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeChecker) checkStmt(n ast.Stmt) error {
	switch n := n.(type) {
	case *ast.Def:
		t.defs[n.Name] = n

	case *ast.Let:
		// infer
		if n.TypeNode.Type == nil {
			n.TypeNode.Type = t.exprType(n.Expr)
		}

		letType := n.TypeNode.Type
		exprType := t.exprType(n.Expr)
		err := t.matchTypes(letType, exprType)
		if err != nil {
			return err
		}

		t.bindings[n.Name] = letType

	case *ast.Assign:
		ty, ok := t.bindings[n.Name]
		if !ok {
			return fmt.Errorf("unknown variable %s", n.Name)
		}
		exprType := t.exprType(n.Expr)
		err := t.matchTypes(ty, exprType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TypeChecker) exprType(n ast.Expr) ast.Type {
	switch n := n.(type) {
	case *ast.Number:
		return ast.IntType{}

	case *ast.String:
		return ast.StringType{}

	default:
		panic(fmt.Sprintf("TypeChecker.exprType: unimplemented %T", n))
	}
}

func (t *TypeChecker) matchTypes(t1 ast.Type, t2 ast.Type) error {
	if (t1 == ast.AnyType{} || t2 == ast.AnyType{}) {
		return nil
	}

	if t1 != t2 {
		return fmt.Errorf("type mismatch")
	}
	return nil
}
