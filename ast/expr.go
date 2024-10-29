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
