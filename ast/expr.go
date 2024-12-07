package ast

type Expr interface{ expr() }

type IntExpr struct {
	Value int
}

func (t IntExpr) expr() {}

type FloatExpr struct {
	Value float64
}

func (t FloatExpr) expr() {}

type BoolExpr struct {
	Value bool
}

func (t BoolExpr) expr() {}

type StringExpr struct {
	Value string
}

func (t StringExpr) expr() {}

type BinaryExpr struct {
	Operands []Expr
	Operator string
}

func (t BinaryExpr) expr() {}

type IdentExpr struct {
	Name string
}

func (t IdentExpr) expr() {}

type FunctionExpr struct {
	Params []string
	Body   Statement
}

func (t FunctionExpr) expr() {}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (t CallExpr) expr() {}

type PipeExpr struct {
	Exprs []Expr
}

func (t PipeExpr) expr() {}

type AssignExpr struct {
	Assignee Expr
	Expr     Expr
}

func (t AssignExpr) expr() {}

type UnaryExpr struct {
	Operator string
	Expr     Expr
}

func (t UnaryExpr) expr() {}

type ArrayExpr struct {
	Exprs []Expr
}

func (t ArrayExpr) expr() {}

type ObjectExpr struct {
	KVs map[string]Expr
}

func (t ObjectExpr) expr() {}

type PostFixExpr struct {
	Expr Expr
	Ops  []PostFixOp
}

func (t PostFixExpr) expr() {}

type PostFixOp interface {
	postFixOp()
}

type ArrayIndexExpr struct {
	Index Expr
}

func (t ArrayIndexExpr) postFixOp() {}
