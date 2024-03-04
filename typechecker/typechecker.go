package typechecker

import (
	"fmt"

	"config-lang/ast"
)

type Binding struct {
	Name string
	Type Type
}

type TypeChecker struct {
	defs     map[string]*ast.Def
	bindings [][]Binding
	src      string
}

func New(src string) *TypeChecker {
	return &TypeChecker{
		defs:     map[string]*ast.Def{},
		bindings: [][]Binding{{}},
		src:      src,
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
		letType := t.astToType(n.Type)
		if letType == nil {
			var err error
			letType, err = t.exprType(n.Expr)
			if err != nil {
				return err
			}
		}

		exprType, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}

		err = t.expectType(letType, exprType)
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
		err = t.expectType(ty, exprType)
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
		err = t.expectType(BoolType{}, exprType)
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
func (t *TypeChecker) exprType(n interface{}) (Type, error) {
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

		err = t.expectType(lhs, rhs)
		if err != nil {
			return nil, err
		}

		return BoolType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

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

		err = t.expectType(NumberType{}, lhs)
		if err != nil {
			return nil, err
		}
		err = t.expectType(NumberType{}, rhs)
		if err != nil {
			return nil, err
		}

		return BoolType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

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

		err = t.expectType(NumberType{}, lhs)
		if err != nil {
			return nil, err
		}
		err = t.expectType(NumberType{}, rhs)
		if err != nil {
			return nil, err
		}

		return NumberType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

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

		err = t.expectType(NumberType{}, lhs)
		if err != nil {
			return nil, err
		}
		err = t.expectType(NumberType{}, rhs)
		if err != nil {
			return nil, err
		}

		return NumberType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

	case *ast.Unary:
		if n.Unary == nil {
			return t.exprType(n.Atom)
		}

		panic("handle unary")
		return BoolType{}, nil

	case *ast.String:
		return StringType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

	case *ast.Number:
		return NumberType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

	case *ast.Bool:
		return BoolType{
			pos:    n.Pos,
			endPos: n.EndPos,
		}, nil

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

func (t *TypeChecker) expectType(expected Type, typ Type) error {
	if (typ.Is(AnyType{})) {
		return nil
	}
	if !typ.Is(expected) {
		return NewTypeError(t.src, expected, typ)
	}

	return nil
}

func (t *TypeChecker) bind(name string, typ Type) {
	idx := len(t.bindings) - 1
	t.bindings[idx] = append(t.bindings[idx], Binding{
		Name: name,
		Type: typ,
	})
}

func (t *TypeChecker) get(name string) (Type, error) {
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

func (t *TypeChecker) astToType(astType ast.Type) Type {
	if astType == nil {
		return nil
	}

	switch n := astType.(type) {
	case *ast.BuiltInType:
		switch n.Name {
		case "string":
			return StringType{}
		case "number":
			return NumberType{}
		case "bool":
			return BoolType{}
		case "any":
			return AnyType{}
		}
	}

	panic(fmt.Sprintf("unimplemented type %T", astType))
}
