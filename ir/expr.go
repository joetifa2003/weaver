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
	Labels      []string
}

func (t FrameExpr) expr() {}

func (t FrameExpr) String(indent int) string {
	var b strings.Builder
	indentStr := strings.Repeat("\t", indent)
	innerIndentStr := strings.Repeat("\t", indent+1)

	b.WriteString(fmt.Sprintf("@frame(vars:%d, params:%d, freeVars:%d) {\n", t.VarCount, t.ParamsCount, len(t.FreeVars)))
	for _, stmt := range t.Body {
		b.WriteString(innerIndentStr)
		b.WriteString(stmt.String(indent + 1))
		b.WriteString("\n")
	}
	b.WriteString(indentStr)
	b.WriteString("}")

	return b.String()
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
	// Binary expressions usually fit on one line, indentation level doesn't propagate inwards
	res := make([]string, 0, len(t.Operands))
	for _, expr := range t.Operands {
		res = append(res, expr.String(0)) // Operands don't get extra indent
	}
	// Join operands with the operator
	return fmt.Sprintf("(%s)", strings.Join(res, fmt.Sprintf(" %s ", t.Operator.String())))
}

type NilExpr struct{}

func (t NilExpr) expr() {}

func (t NilExpr) String(indent int) string {
	return "nil"
}

type IfExpr struct {
	Condition Expr
	TrueExpr  Expr
	FalseExpr Expr
}

func (t IfExpr) expr() {}

func (t IfExpr) String(indent int) string {
	return fmt.Sprintf("(%s ? %s | %s)", t.Condition.String(0), t.TrueExpr.String(0), t.FalseExpr.String(0))
}

type ReturnExpr struct {
	Expr Expr
}

func (t ReturnExpr) expr() {}

func (t ReturnExpr) String(i int) string {
	// Return expression itself doesn't add indent, but contains an expression
	return fmt.Sprintf("return %s", t.Expr.String(0)) // Contained expr doesn't get extra indent
}

type IndexExpr struct {
	Expr  Expr
	Index Expr
}

func (t IndexExpr) expr() {}

func (t IndexExpr) String(i int) string {
	// Index expression itself doesn't add indent
	return fmt.Sprintf("%s[%s]", t.Expr.String(0), t.Index.String(0)) // Contained exprs don't get extra indent
}

type CallExpr struct {
	Expr Expr
	Args []Expr
}

func (t CallExpr) expr() {}

func (t CallExpr) String(i int) string {
	// Call expression itself doesn't add indent
	res := make([]string, 0, len(t.Args))
	for _, arg := range t.Args {
		res = append(res, arg.String(0)) // Args don't get extra indent
	}
	return fmt.Sprintf("%s(%s)", t.Expr.String(0), strings.Join(res, ", ")) // Callee doesn't get extra indent
}

type TryExpr struct {
	Expr Expr
}

func (t TryExpr) expr() {}

func (t TryExpr) String(i int) string {
	return fmt.Sprintf("try %s", t.Expr.String(0))
}

type RaiseExpr struct {
	Expr Expr
}

func (t RaiseExpr) expr() {}

func (t RaiseExpr) String(i int) string {
	return fmt.Sprintf("raise %s", t.Expr.String(0))
}
