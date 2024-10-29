package typechecker

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/joetifa2003/weaver/ast"
)

type Binding struct {
	Name string
	Type Type
}

type TypeChecker struct {
	defs     map[string]Type
	bindings [][]Binding
	src      string
}

func New(src string) *TypeChecker {
	return &TypeChecker{
		defs:     map[string]Type{},
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
		t.defs[n.Name] = t.astToType(n.Type)

	case *ast.Let:
		// infer
		letType := t.astToType(n.Type)
		if letType == nil {
			var err error
			letType, _, err = t.exprType(n.Expr)
			if err != nil {
				return err
			}
		}

		exprType, exprPos, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}

		err = t.expectType(letType, exprType, exprPos)
		if err != nil {
			return err
		}

		t.bind(n.Name, letType)

	case *ast.Assign:
		ty, err := t.get(n.Name)
		if err != nil {
			return err
		}

		exprType, exprPos, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}
		err = t.expectType(ty, exprType, exprPos)
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
		exprType, exprPos, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}
		err = t.expectType(BoolType{}, exprType, exprPos)
		if err != nil {
			return err
		}
		err = t.checkStmt(n.Statement)
		if err != nil {
			return err
		}

	case *ast.Fn:
		for _, s := range n.Statements {
			err := t.checkStmt(s)
			if err != nil {
				return err
			}
		}

		args := make([]Type, len(n.Args))
		for i, arg := range n.Args {
			args[i] = t.astToType(arg.Type)
		}
		returnType := t.astToType(n.ReturnType)

		t.bind(n.Name, FnType{
			Args:       args,
			ReturnType: returnType,
		})

	case *ast.Echo:
		_, _, err := t.exprType(n.Expr)
		if err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("TypeChecker.checkStmt: unimplemented %T", n))
	}

	return nil
}

// TODO: Don't allow things like true + true
func (t *TypeChecker) exprType(n interface{}) (Type, lexer.Position, error) {
	switch n := n.(type) {
	case *ast.Expr:
		return t.exprType(n.Equality)

	case *ast.Equality:
		lhs, lhsPos, err := t.exprType(n.Left)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		if n.Right == nil {
			return lhs, lhsPos, nil
		}

		rhs, rhsPos, err := t.exprType(n.Right)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		err = t.expectType(lhs, rhs, rhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		return BoolType{}, n.Pos, nil

	case *ast.Comparison:
		lhs, lhsPos, err := t.exprType(n.Left)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		if n.Right == nil {
			return lhs, lhsPos, nil
		}

		rhs, rhsPos, err := t.exprType(n.Right)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		err = t.expectType(NumberType{}, lhs, lhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}
		err = t.expectType(NumberType{}, rhs, rhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		return BoolType{}, n.Pos, nil

	case *ast.Addition:
		lhs, lhsPos, err := t.exprType(n.Left)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		if n.Right == nil {
			return lhs, lhsPos, nil
		}

		rhs, rhsPos, err := t.exprType(n.Right)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		err = t.expectType(NumberType{}, lhs, lhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}
		err = t.expectType(NumberType{}, rhs, rhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		return BoolType{}, n.Pos, nil

	case *ast.Multiplication:
		lhs, lhsPos, err := t.exprType(n.Left)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		if n.Right == nil {
			return lhs, lhsPos, nil
		}

		rhs, rhsPos, err := t.exprType(n.Right)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		err = t.expectType(NumberType{}, lhs, lhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}
		err = t.expectType(NumberType{}, rhs, rhsPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		return BoolType{}, n.Pos, nil

	case *ast.Unary:
		if n.Unary == nil {
			return t.exprType(n.Atom)
		}

		panic("handle unary")
		return BoolType{}, n.Pos, nil

	case *ast.String:
		return StringType{}, n.Pos, nil

	case *ast.Number:
		return NumberType{}, n.Pos, nil

	case *ast.Bool:
		return BoolType{}, n.Pos, nil

	case *ast.Ident:
		identType, err := t.get(n.Name)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		return identType, n.Pos, nil

	case *ast.Object:
		res := ObjectType{
			Fields: map[string]Type{},
		}

		for _, f := range n.Fields {
			var err error
			res.Fields[f.Name], _, err = t.exprType(f.Expr)
			if err != nil {
				return nil, lexer.Position{}, err
			}
		}

		return res, n.Pos, nil

	case *ast.Paren:
		return t.exprType(n.Expr)

	case *ast.Call:
		// TODO: CHECK ARGS TYPES
		identTyp, identPos, err := t.exprType(n.Name)
		if err != nil {
			return nil, lexer.Position{}, err
		}
		err = softExpect[FnType](t, identTyp, identPos)
		if err != nil {
			return nil, lexer.Position{}, err
		}

		fn := identTyp.(FnType)
		if len(fn.Args) != len(n.Args) {
			// TODO: Create an error type called, PositionalError
			// to be a more generic type errror
			return nil, identPos, errors.New("Invalid number of args")
		}

		for i, arg := range n.Args {
			typ, pos, err := t.exprType(arg)
			if err != nil {
				return typ, pos, err
			}

			err = t.expectType(fn.Args[i], typ, pos)
			if err != nil {
				return typ, pos, err
			}
		}

		rt := fn.ReturnType
		return rt, identPos, nil

	default:
		panic(fmt.Sprintf("TypeChecker.exprType: unimplemented %T", n))
	}
}

func (t *TypeChecker) expectType(expected Type, typ Type, pos lexer.Position) error {
	if _, ok := expected.(AnyType); ok {
		return nil
	}

	if _, ok := typ.(AnyType); ok {
		return nil
	}

	if !typ.IsAssignableTo(expected) {
		return NewTypeError(t.src, expected, typ, pos)
	}

	return nil
}

func softExpect[E Type](t *TypeChecker, typ Type, pos lexer.Position) error {
	_, ok := typ.(E)
	if !ok {
		var zero E
		return NewTypeError(t.src, zero, typ, pos)
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

// TODO: Handle nullable types
func (t *TypeChecker) astToType(astType *ast.Type) Type {
	if astType == nil {
		return nil
	}

	res := make([]Type, len(astType.Variants))
	for i, n := range astType.Variants {
		switch n := n.(type) {
		case *ast.BuiltInType:
			switch n.Name {
			case "string":
				res[i] = StringType{BaseType: BaseType{nullable: n.Nullable}}
			case "number":
				res[i] = NumberType{BaseType: BaseType{nullable: n.Nullable}}
			case "bool":
				res[i] = BoolType{BaseType: BaseType{nullable: n.Nullable}}
			case "any":
				res[i] = AnyType{BaseType: BaseType{nullable: n.Nullable}}
			}

		case *ast.CustomType:
			// TODO: return a proper error here if things are not defined
			res[i] = t.defs[n.Name]

		case *ast.ObjectType:
			obj := ObjectType{
				Fields: map[string]Type{},
			}

			for _, f := range n.Fields {
				obj.Fields[f.Name] = t.astToType(f.Type)
			}

			obj.nullable = n.Nullable

			res[i] = obj
		}
	}

	if len(res) == 1 {
		return res[0]
	}

	return Variant{
		Variants: res,
	}
}
