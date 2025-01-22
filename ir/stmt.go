package ir

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

type LoopStmt struct {
	Condition Expr
	Body      Statement
}

func (t LoopStmt) stmt() {}

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
