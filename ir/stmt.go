package ir

import (
	"fmt"
	"strings"
)

type Program struct {
	VarCount   int
	Statements []Statement
	Labels     []string
}

type Statement interface {
	stmt()
	String(indent int) string
}

func (p *Program) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# VarCount: %d\n", p.VarCount))
	for _, stmt := range p.Statements {
		b.WriteString(stmt.String(0))
		b.WriteString("\n")
	}
	return b.String()
}

type BlockStmt struct {
	Statements []Statement
}

func (t BlockStmt) stmt() {}

func (t BlockStmt) String(i int) string {
	var b strings.Builder
	indentStr := strings.Repeat("\t", i)
	innerIndentStr := strings.Repeat("\t", i+1)

	b.WriteString("{\n")
	for _, stmt := range t.Statements {
		b.WriteString(innerIndentStr)
		b.WriteString(stmt.String(i + 1))
		b.WriteString("\n")
	}
	b.WriteString(indentStr)
	b.WriteString("}")

	return b.String()
}

type LoopStmt struct {
	Body Statement
}

func (t LoopStmt) stmt() {}

func (t LoopStmt) String(i int) string {
	return fmt.Sprintf("loop %s", t.Body.String(i))
}

type IfStmt struct {
	Condition   Expr
	Body        Statement
	Alternative *Statement
}

func (t IfStmt) stmt() {}

func (t IfStmt) String(i int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("if (%s) ", t.Condition.String(i)))
	b.WriteString(t.Body.String(i))

	if t.Alternative != nil {
		b.WriteString(" else ")
		b.WriteString((*t.Alternative).String(i))
	}

	return b.String()
}

type ExpressionStmt struct {
	Expr Expr
}

func (t ExpressionStmt) stmt() {}

func (t ExpressionStmt) String(i int) string {
	return t.Expr.String(i)
}

type ContinueStmt struct {
	IncrementStatement Statement
}

func (t ContinueStmt) stmt() {}

func (t ContinueStmt) String(i int) string {
	return "continue"
}

type BreakStmt struct{}

func (t BreakStmt) stmt() {}

func (t BreakStmt) String(i int) string {
	return "break"
}

type LabelStmt struct {
	Name string
}

func (t LabelStmt) stmt() {}

func (t LabelStmt) String(i int) string {
	return fmt.Sprintf("label %s", t.Name)
}

type GotoStmt struct {
	Name string
}

func (t GotoStmt) stmt() {}

func (t GotoStmt) String(i int) string {
	return fmt.Sprintf("goto %s", t.Name)
}
