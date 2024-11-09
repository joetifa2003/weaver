package ast

type Statement interface{ stmt() }

type LetStmt struct {
	Name string
	Expr Expr
}

func (t LetStmt) stmt() {}

type EchoStmt struct {
	Expr Expr
}

func (t EchoStmt) stmt() {}

type BlockStmt struct {
	Statements []Statement
}

func (t BlockStmt) stmt() {}

type AssignStmt struct {
	Name string
	Expr Expr
}

func (t AssignStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      Statement
}

func (t WhileStmt) stmt() {}

type IfStmt struct {
	Condition Expr
	Body      Statement
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
