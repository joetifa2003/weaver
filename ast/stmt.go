package ast

type Statement interface{ stmt() }

type LetStmt struct {
	Name string
	Expr Expr
}

func (t LetStmt) stmt() {}

type BlockStmt struct {
	Statements []Statement
}

func (t BlockStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      Statement
}

func (t WhileStmt) stmt() {}

type ForStmt struct {
	InitStmt  Statement
	Condition Expr
	Increment Expr
	Body      Statement
}

func (t ForStmt) stmt() {}

type IfStmt struct {
	Condition   Expr
	Body        Statement
	Alternative *Statement
}

func (t IfStmt) stmt() {}

type ExprStmt struct {
	Expr Expr
}

func (t ExprStmt) stmt() {}

type ContinueStmt struct{}

func (t ContinueStmt) stmt() {}

type BreakStmt struct{}

func (t BreakStmt) stmt() {}

type MatchStmt struct {
	Expr  Expr
	Cases []MatchCase
}

func (t MatchStmt) stmt() {}

type MatchCase struct {
	Condition MatchCaseCondition
	ExtraCond *Expr
	Body      Statement
}

type MatchCaseCondition interface {
	matchCaseCondition()
}

type MatchCaseInt struct {
	Value int
}

func (t MatchCaseInt) matchCaseCondition() {}

type MatchCaseFloat struct {
	Value float64
}

func (t MatchCaseFloat) matchCaseCondition() {}

type MatchCaseString struct {
	Value string
}

func (t MatchCaseString) matchCaseCondition() {}

type MatchCaseArray struct {
	Conditions []MatchCaseCondition
}

func (t MatchCaseArray) matchCaseCondition() {}

type MatchCaseObject struct {
	KVs map[string]*MatchCaseCondition
}

func (t MatchCaseObject) matchCaseCondition() {}

type MatchCaseIdent struct {
	Name string
}

func (t MatchCaseIdent) matchCaseCondition() {}

type MatchCaseRange struct {
	Begin Expr
	End   Expr
}

func (t MatchCaseRange) matchCaseCondition() {}

type MatchCaseTypeError struct {
	Message MatchCaseCondition
	Data    MatchCaseCondition
}

func (t MatchCaseTypeError) matchCaseCondition() {}

type MatchCaseTypeString struct {
	Cond *MatchCaseCondition
}

func (t MatchCaseTypeString) matchCaseCondition() {}

type MatchCaseTypeNumber struct {
	Cond *MatchCaseCondition
}

func (t MatchCaseTypeNumber) matchCaseCondition() {}

type MatchCaseOr struct {
	Conditions []MatchCaseCondition
}

func (t MatchCaseOr) matchCaseCondition() {}
