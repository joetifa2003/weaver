package ir

import (
	"fmt"
	"strings"
)

type Expr interface {
	expr()
	String() string
}

type Var struct {
	Idx   int
	Scope VarScope
}

func (v Var) String() string {
	switch v.Scope {
	case VarScopeLocal:
		return fmt.Sprintf("v%d_%s", v.Idx, v.Scope)
	case VarScopeFree:
		return fmt.Sprintf("f%d_%s", v.Idx, v.Scope)
	default:
		panic(fmt.Sprintf("unknown scope %s", v.Scope))
	}
}

type IntExpr struct {
	Value int
}

func (t IntExpr) expr() {}

func (t IntExpr) String() string {
	return fmt.Sprintf("%d", t.Value)
}

type FloatExpr struct {
	Value float64
}

func (t FloatExpr) expr() {}

func (t FloatExpr) String() string {
	return fmt.Sprintf("%f", t.Value)
}

type BoolExpr struct {
	Value bool
}

func (t BoolExpr) expr() {}

func (t BoolExpr) String() string {
	if t.Value {
		return "true"
	}
	return "false"
}

type StringExpr struct {
	Value string
}

func (t StringExpr) expr() {}

func (t StringExpr) String() string {
	return fmt.Sprintf("\"%s\"", t.Value)
}

type BuiltInExpr struct {
	Name string
}

func (t BuiltInExpr) expr() {}

func (t BuiltInExpr) String() string {
	return fmt.Sprintf("@%s", t.Name)
}

type LoadExpr struct {
	Var
}

func (t LoadExpr) expr() {}

func (t LoadExpr) String() string {
	return fmt.Sprintf("%s", t.Var)
}

type IdxAssignExpr struct {
	Assignee Expr
	Index    Expr
	Value    Expr
}

func (t IdxAssignExpr) expr() {}

func (t IdxAssignExpr) String() string {
	return fmt.Sprintf("%s[%s] = %s", t.Assignee.String(), t.Index.String(), t.Value.String())
}

type VarAssignExpr struct {
	Var   Var
	Value Expr
}

func (t VarAssignExpr) expr() {}

func (t VarAssignExpr) String() string {
	return fmt.Sprintf("%s = %s", t.Var, t.Value.String())
}

type VarIncrementExpr struct {
	Var Var
}

func (t VarIncrementExpr) expr() {}

func (t VarIncrementExpr) String() string {
	return fmt.Sprintf("%s++", t.Var)
}

type UnaryOp int

const (
	UnaryOpNot UnaryOp = iota
	UnaryOpNegate
)

func (t UnaryOp) String() string {
	switch t {
	case UnaryOpNot:
		return "!"
	case UnaryOpNegate:
		return "-"
	default:
		panic(fmt.Sprintf("unimplemented %T", t))
	}
}

type UnaryExpr struct {
	Operator UnaryOp
	Expr     Expr
}

func (t UnaryExpr) expr() {}

func (t UnaryExpr) String() string {
	return fmt.Sprintf("%s(%s)", t.Operator.String(), t.Expr.String())
}

type ArrayExpr struct {
	Exprs []Expr
}

func (t ArrayExpr) expr() {}

func (t ArrayExpr) String() string {
	var res []string
	for _, expr := range t.Exprs {
		res = append(res, expr.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

type ObjectExpr struct {
	KVs map[string]Expr
}

func (t ObjectExpr) expr() {}

func (t ObjectExpr) String() string {
	var res []string
	for key, expr := range t.KVs {
		res = append(res, fmt.Sprintf("%s: %s", key, expr.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(res, ", "))
}

type FrameExpr struct {
	VarCount    int
	ParamsCount int
	FreeVars    []Var
	Body        []Statement
}

func (t FrameExpr) expr() {}

func (t FrameExpr) String() string {
	return "frame"
}

type BinaryOp int

const (
	BinaryOpAdd BinaryOp = iota
	BinaryOpSub
	BinaryOpMul
	BinaryOpDiv
	BinaryOpMod
	BinaryOpEq
	BinaryOpNeq
	BinaryOpGt
	BinaryOpLt
	BinaryOpGte
	BinaryOpLte
	BinaryOpAnd
	BinaryOpOr
)

func (t BinaryOp) String() string {
	switch t {
	case BinaryOpAdd:
		return "+"
	case BinaryOpSub:
		return "-"
	case BinaryOpMul:
		return "*"
	case BinaryOpDiv:
		return "/"
	case BinaryOpMod:
		return "%"
	case BinaryOpEq:
		return "=="
	case BinaryOpNeq:
		return "!="
	case BinaryOpGt:
		return ">"
	case BinaryOpLt:
		return "<"
	case BinaryOpGte:
		return ">="
	case BinaryOpLte:
		return "<="
	case BinaryOpAnd:
		return "&&"
	case BinaryOpOr:
		return "||"
	default:
		panic(fmt.Sprintf("unimplemented %T", t))
	}
}

type BinaryExpr struct {
	Operator BinaryOp
	Operands []Expr
}

func (t BinaryExpr) expr() {}

func (t BinaryExpr) String() string {
	var res []string
	for i, expr := range t.Operands {
		res = append(res, expr.String())
		if i != len(t.Operands)-1 {
			res = append(res, t.Operator.String())
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(res, " "))
}

type PostFixExpr struct {
	Expr Expr
	Ops  []PostFixOp
}

func (t PostFixExpr) expr() {}

func (t PostFixExpr) String() string {
	var res []string
	for _, op := range t.Ops {
		res = append(res, op.String())
	}
	return fmt.Sprintf("%s%s", t.Expr.String(), strings.Join(res, ""))
}

type PostFixOp interface {
	postFixOp()
	String() string
}

type IndexOp struct {
	Index Expr
}

func (t IndexOp) postFixOp() {}

func (t IndexOp) String() string {
	return fmt.Sprintf("[%s]", t.Index.String())
}

type CallOp struct {
	Args []Expr
}

func (t CallOp) postFixOp() {}

func (t CallOp) String() string {
	var res []string
	for _, arg := range t.Args {
		res = append(res, arg.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(res, ", "))
}

type NilExpr struct{}

func (t NilExpr) expr() {}

func (t NilExpr) String() string {
	return "nil"
}
