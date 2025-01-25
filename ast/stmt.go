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

type ReturnStmt struct {
	Expr Expr
}

func (t ReturnStmt) stmt() {}

type ExprStmt struct {
	Expr Expr
}

func (t ExprStmt) stmt() {}

type ContinueStmt struct{}

func (t ContinueStmt) stmt() {}

type BreakStmt struct{}

func (t BreakStmt) stmt() {}
