package ir

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface {
	expr()
	String(indent int) string
}

type Var struct {
	Idx   int
	Scope VarScope
	Ref   bool
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

func (t IntExpr) String(indent int) string {
	return strconv.Itoa(t.Value)
}

type FloatExpr struct {
	Value float64
}

func (t FloatExpr) expr() {}

func (t FloatExpr) String(indent int) string {
	return fmt.Sprintf("%f", t.Value)
}

type BoolExpr struct {
	Value bool
}

func (t BoolExpr) expr() {}

func (t BoolExpr) String(indent int) string {
	if t.Value {
		return "true"
	}
	return "false"
}

type StringExpr struct {
	Value string
}

func (t StringExpr) expr() {}

func (t StringExpr) String(indent int) string {
	return fmt.Sprintf("\"%s\"", t.Value)
}

type BuiltInExpr struct {
	Name string
}

func (t BuiltInExpr) expr() {}

func (t BuiltInExpr) String(indent int) string {
	return "@" + t.Name
}

type LoadExpr struct {
	Var
}

func (t LoadExpr) expr() {}

func (t LoadExpr) String(indent int) string {
	return t.Var.String()
}

type IdxAssignExpr struct {
	Assignee Expr
	Index    Expr
	Value    Expr
}

func (t IdxAssignExpr) expr() {}

func (t IdxAssignExpr) String(indent int) string {
	return fmt.Sprintf("%s[%s] = %s", t.Assignee.String(indent), t.Index.String(indent), t.Value.String(indent))
}

type VarAssignExpr struct {
	Var   Var
	Value Expr
}

func (t VarAssignExpr) expr() {}

func (t VarAssignExpr) String(indent int) string {
	if t.Var.Ref {
		return fmt.Sprintf("%s_ref = %s", t.Var, t.Value.String(indent))
	}
	return fmt.Sprintf("%s = %s", t.Var, t.Value.String(indent))
}

type VarIncrementExpr struct {
	Var Var
}

func (t VarIncrementExpr) expr() {}

func (t VarIncrementExpr) String(indent int) string {
	return fmt.Sprintf("%s++", t.Var)
}

type VarDecrementExpr struct {
	Var Var
}

func (t VarDecrementExpr) expr() {}

func (t VarDecrementExpr) String(indent int) string {
	return fmt.Sprintf("%s--", t.Var)
}

type UnaryOp int

const (
	UnaryOpNot UnaryOp = iota
	UnaryOpNegate
)

func (t UnaryOp) String(index int) string {
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

func (t UnaryExpr) String(indent int) string {
	return fmt.Sprintf("%s(%s)", t.Operator.String(indent), t.Expr.String(indent))
}

type ArrayExpr struct {
	Exprs []Expr
}

func (t ArrayExpr) expr() {}

func (t ArrayExpr) String(indent int) string {
	res := make([]string, 0, len(t.Exprs))
	for _, expr := range t.Exprs {
		res = append(res, expr.String(indent))
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

type ObjectExpr struct {
	KVs map[string]Expr
}

func (t ObjectExpr) expr() {}

func (t ObjectExpr) String(indent int) string {
	res := make([]string, 0, len(t.KVs))
	for key, expr := range t.KVs {
		res = append(res, fmt.Sprintf("%s: %s", key, expr.String(indent)))
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

func (t FrameExpr) String(indent int) string {
	stmts := make([]string, 0, len(t.Body))
	for _, stmt := range t.Body {
		stmts = append(stmts, stmt.String(indent+1))
	}

	return fmt.Sprintf(
		"@frame(vars:%d, freeVars:%d,) %s",
		t.VarCount,
		len(t.FreeVars),
		strings.Join(stmts, "\n"),
	)
}

type BinaryOp int

const (
	BinaryOpUnknown BinaryOp = iota
	BinaryOpAdd
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

func (t BinaryExpr) String(indent int) string {
	res := make([]string, 0, len(t.Operands))
	for i, expr := range t.Operands {
		res = append(res, expr.String(indent))
		if i != len(t.Operands)-1 {
			res = append(res, t.Operator.String())
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(res, " "))
}

type NilExpr struct{}

func (t NilExpr) expr() {}

func (t NilExpr) String(indent int) string {
	return "nil"
}

type ModuleLoadExpr struct {
	Name string
	Load string
}

func (t ModuleLoadExpr) expr() {}

func (t ModuleLoadExpr) String(indent int) string {
	return fmt.Sprintf("%s:%s", t.Name, t.Load)
}

type IfExpr struct {
	Condition Expr
	TrueExpr  Expr
	FalseExpr Expr
}

func (t IfExpr) expr() {}

func (t IfExpr) String(indent int) string {
	return fmt.Sprintf("if (%s) %s else %s", t.Condition.String(indent), t.TrueExpr.String(indent), t.FalseExpr.String(indent))
}

type ReturnExpr struct {
	Expr Expr
}

func (t ReturnExpr) expr() {}

func (t ReturnExpr) String(i int) string {
	return fmt.Sprintf("%sreturn %s", strings.Repeat("\t", i), t.Expr.String(i))
}

type IndexExpr struct {
	Expr  Expr
	Index Expr
}

func (t IndexExpr) expr() {}

func (t IndexExpr) String(i int) string {
	return fmt.Sprintf("%s[%s]", t.Expr.String(i), t.Index.String(i))
}

type CallExpr struct {
	Expr Expr
	Args []Expr
}

func (t CallExpr) expr() {}

func (t CallExpr) String(i int) string {
	res := make([]string, 0, len(t.Args))
	for _, arg := range t.Args {
		res = append(res, arg.String(i))
	}
	return fmt.Sprintf("%s(%s)", t.Expr.String(i), strings.Join(res, ", "))
}
