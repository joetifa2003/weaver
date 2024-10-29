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
