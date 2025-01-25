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

type IdxAssignExpr struct {
	Assignee Expr
	Index    Expr
	Value    Expr
}

func (t IdxAssignExpr) expr() {}

type VarAssignExpr struct {
	Name  string
	Value Expr
}

func (t VarAssignExpr) expr() {}

type UnaryOp int

const (
	UnaryOpNot UnaryOp = iota
	UnaryOpNegate
)

type UnaryExpr struct {
	Operator UnaryOp
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

type BinaryExpr struct {
	Operands []Expr
	Operator BinaryOp
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

type IndexOp struct {
	Index Expr
}

func (t IndexOp) postFixOp() {}

type CallOp struct {
	Args []Expr
}

func (t CallOp) postFixOp() {}
