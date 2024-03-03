package typechecker

import (
	"errors"
	"fmt"

	"config-lang/ast"
)

type Binding struct {
	Name string
	Type ast.Type
}

type TypeChecker struct {
	defs     map[string]*ast.Def
	bindings [][]Binding
}

func New() *TypeChecker {
	return &TypeChecker{
		defs:     map[string]*ast.Def{},
		bindings: [][]Binding{{}},
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
			var err error
			n.TypeNode.Type, err = t.exprType(n.Expr)
			if err != nil {
				return err
			}
		}

		letType := n.TypeNode.Type
		exprType, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}

		err = t.matchTypes(letType, exprType, nil)
		if err != nil {
			return err
		}

		t.bind(n.Name, letType)

	case *ast.Assign:
		ty, err := t.get(n.Name)
		if err != nil {
			return err
		}

		exprType, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}
		err = t.matchTypes(ty, exprType, nil)
		if err != nil {
			return err
		}

	case *ast.Block:
		t.push()
		for _, stmt := range n.Statements {
			err := t.checkStmt(stmt)
			if err != nil {
				return err
			}
		}
		t.pop()

	case *ast.If:
		exprType, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}
		err = t.matchTypes(exprType, ast.BoolType{}, nil)
		if err != nil {
			return err
		}
		err = t.checkStmt(n.Statement)
		if err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("TypeChecker.checkStmt: unimplemented %T", n))
	}

	return nil
}

// TODO: Don't allow things like true + true
func (t *TypeChecker) exprType(n interface{}) (ast.Type, error) {
	switch n := n.(type) {
	case *ast.Expr:
		return t.exprType(n.Equality)

	case *ast.Equality:
		lhs, err := t.exprType(n.Left)
		if err != nil {
			return nil, err
		}

		if n.Right == nil {
			return lhs, nil
		}

		rhs, err := t.exprType(n.Right)
		if err != nil {
			return nil, err
		}

		err = t.matchTypes(lhs, rhs, nil)
		if err != nil {
			return nil, err
		}

		return ast.BoolType{}, nil

	case *ast.Comparison:
		lhs, err := t.exprType(n.Left)
		if err != nil {
			return nil, err
		}

		if n.Right == nil {
			return lhs, nil
		}

		rhs, err := t.exprType(n.Right)
		if err != nil {
			return nil, err
		}

		err = t.matchTypes(lhs, rhs, ast.IntType{})
		if err != nil {
			return nil, err
		}

		return ast.BoolType{}, nil

	case *ast.Addition:
		lhs, err := t.exprType(n.Left)
		if err != nil {
			return nil, err
		}

		if n.Right == nil {
			return lhs, nil
		}

		rhs, err := t.exprType(n.Right)
		if err != nil {
			return nil, err
		}

		err = t.matchTypes(lhs, rhs, ast.IntType{})
		if err != nil {
			return nil, err
		}

		return ast.IntType{}, nil

	case *ast.Multiplication:
		lhs, err := t.exprType(n.Left)
		if err != nil {
			return nil, err
		}

		if n.Right == nil {
			return lhs, nil
		}

		rhs, err := t.exprType(n.Right)
		if err != nil {
			return nil, err
		}

		err = t.matchTypes(lhs, rhs, ast.IntType{})
		if err != nil {
			return nil, err
		}

		return ast.IntType{}, nil

	case *ast.Unary:
		if n.Unary == nil {
			return t.exprType(n.Atom)
		}

		panic("handle unary")
		return ast.IntType{}, nil

	case *ast.String:
		return ast.StringType{}, nil

	case *ast.Number:
		return ast.IntType{}, nil

	case *ast.Bool:
		return ast.BoolType{}, nil

	case *ast.Ident:
		identType, err := t.get(n.Name)
		if err != nil {
			return nil, err
		}

		return identType, nil

	default:
		panic(fmt.Sprintf("TypeChecker.exprType: unimplemented %T", n))
	}
}

var (
	ErrMismatch = errors.New("Type mismatch")
)

func (t *TypeChecker) matchTypes(t1 ast.Type, t2 ast.Type, eq ast.Type) error {
	if (t1 == ast.AnyType{} || t2 == ast.AnyType{}) {
		return nil
	}
	if t1 != t2 {
		return fmt.Errorf("%w: expected %T to be %T", ErrMismatch, t2, t1)
	}

	if eq != nil {
		if t1 != eq {
			return fmt.Errorf("%w: expected %T to be %T", ErrMismatch, t1, eq)
		}
	}

	return nil
}

func (t *TypeChecker) bind(name string, typ ast.Type) {
	idx := len(t.bindings) - 1
	t.bindings[idx] = append(t.bindings[idx], Binding{
		Name: name,
		Type: typ,
	})
}

func (t *TypeChecker) get(name string) (ast.Type, error) {
	for i := len(t.bindings) - 1; i >= 0; i-- {
		for _, b := range t.bindings[i] {
			if b.Name == name {
				return b.Type, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown variable %s", name)
}

func (t *TypeChecker) push() {
	t.bindings = append(t.bindings, []Binding{})
}

func (t *TypeChecker) pop() {
	t.bindings = t.bindings[:len(t.bindings)-1]
}
