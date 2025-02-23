package ir

import (
	"fmt"
	"strings"
)

type Program struct {
	VarCount   int
	Statements []Statement
}

type Statement interface {
	stmt()
	String(indent int) string
}

type BlockStmt struct {
	Statements []Statement
}

func (t BlockStmt) stmt() {}

func (t BlockStmt) String(i int) string {
	res := make([]string, 0, len(t.Statements))
	for _, stmt := range t.Statements {
		res = append(res, stmt.String(i+1))
	}
	return fmt.Sprintf("%s{\n%s\n%s}", strings.Repeat("\t", i), strings.Join(res, "\n"), strings.Repeat("\t", i))
}

type LoopStmt struct {
	Body Statement
}

func (t LoopStmt) stmt() {}

func (t LoopStmt) String(i int) string {
	return fmt.Sprintf("%sloop %s", strings.Repeat("\t", i), t.Body.String(i+1))
}

type IfStmt struct {
	Condition   Expr
	Body        Statement
	Alternative *Statement
}

func (t IfStmt) stmt() {}

func (t IfStmt) String(i int) string {
	var res []string
	if t.Alternative != nil {
		res = append(res, fmt.Sprintf(
			"%sif %s {\n%s\n%s} else {\n%s\n%s}",
			strings.Repeat("\t", i),
			t.Condition.String(i),
			t.Body.String(i+1),
			strings.Repeat("\t", i),
			(*t.Alternative).String(i+1),
			strings.Repeat("\t", i),
		),
		)
	} else {
		res = append(res, fmt.Sprintf("%sif %s {\n%s\n%s}", strings.Repeat("\t", i), t.Condition.String(i), t.Body.String(i+1), strings.Repeat("\t", i)))
	}
	return strings.Join(res, "\n")
}

type ReturnStmt struct {
	Expr Expr
}

func (t ReturnStmt) stmt() {}

func (t ReturnStmt) String(i int) string {
	return fmt.Sprintf("%sreturn %s", strings.Repeat("\t", i), t.Expr.String(i))
}

type ExpressionStmt struct {
	Expr Expr
}

func (t ExpressionStmt) stmt() {}

func (t ExpressionStmt) String(i int) string {
	return fmt.Sprintf("%s%s", strings.Repeat("\t", i), t.Expr.String(i))
}

type ContinueStmt struct {
	IncrementStatement Statement
}

func (t ContinueStmt) stmt() {}

func (t ContinueStmt) String(i int) string {
	return strings.Repeat("\t", i) + "continue"
}

type BreakStmt struct{}

func (t BreakStmt) stmt() {}

func (t BreakStmt) String(i int) string {
	return strings.Repeat("\t", i) + "break"
}
