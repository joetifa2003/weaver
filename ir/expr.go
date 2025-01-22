package ir

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

type IdentExpr struct {
	Name string
}

func (t IdentExpr) expr() {}

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

type FunctionExpr struct {
	Params []string
	Body   Statement
}

func (t FunctionExpr) expr() {}

type BinaryExpr struct {
	Operands []Expr
	Operator string
}

func (t BinaryExpr) expr() {}

type PostFixExpr struct {
	Expr Expr
	Ops  []PostFixOp
}

func (t PostFixExpr) expr() {}

type PostFixOp interface {
	postFixOp()
}

type IndexExpr struct {
	Index Expr
}

func (t IndexExpr) postFixOp() {}
