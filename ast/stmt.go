package ast

type Statement interface{ stmt() }

type LetStmt struct {
	Name string
	Expr Expr
}

func (t LetStmt) stmt() {}
