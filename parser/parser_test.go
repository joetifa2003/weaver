package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joetifa2003/weaver/ast"
	"github.com/joetifa2003/weaver/internal/pargo"
)

func TestIntExpr(t *testing.T) {
	assert := require.New(t)

	p := intExpr()
	expr, err := pargo.Parse(p, newLexer(), "123")
	require.NoError(t, err)

	intExpr, ok := expr.(ast.IntExpr)
	assert.True(ok)
	assert.Equal(123, intExpr.Value)
}

func TestFloatExpr(t *testing.T) {
	assert := require.New(t)

	p := floatExpr()
	expr, err := pargo.Parse(p, newLexer(), "123.456")
	require.NoError(t, err)

	floatExpr, ok := expr.(ast.FloatExpr)
	assert.True(ok)
	assert.Equal(123.456, floatExpr.Value)
}

func TestBooleanExpr(t *testing.T) {
	assert := require.New(t)

	p := booleanExpr()

	expr, err := pargo.Parse(p, newLexer(), "true")
	require.NoError(t, err)

	boolExpr, ok := expr.(ast.BoolExpr)
	assert.True(ok)
	assert.True(boolExpr.Value)

	expr, err = pargo.Parse(p, newLexer(), "false")
	require.NoError(t, err)

	boolExpr, ok = expr.(ast.BoolExpr)
	assert.True(ok)
	assert.False(boolExpr.Value)

	expr, err = pargo.Parse(p, newLexer(), "foo")
	require.Error(t, err)
}

func TestStringExpr(t *testing.T) {
	assert := require.New(t)

	p := stringExpr()

	expr, err := pargo.Parse(p, newLexer(), `"foo"`)
	require.NoError(t, err)

	stringExpr, ok := expr.(ast.StringExpr)
	assert.True(ok)
	assert.Equal("foo", stringExpr.Value)
}

func TestVarDeclStmt(t *testing.T) {
	assert := require.New(t)

	p := varDeclStmt()

	stmt, err := pargo.Parse(p, newLexer(), "foo := 123")
	require.NoError(t, err)

	letStmt, ok := stmt.(ast.LetStmt)
	assert.True(ok)
	assert.Equal("foo", letStmt.Name)
	assert.Equal(ast.IntExpr{Value: 123}, letStmt.Expr)
}

func TestBinaryExpr(t *testing.T) {
	assert := require.New(t)

	p := addExpr()

	expr, err := pargo.Parse(p, newLexer(), "123 + 456")
	require.NoError(t, err)

	binaryExpr, ok := expr.(ast.BinaryExpr)
	assert.True(ok)
	assert.Equal(
		ast.BinaryExpr{
			Operands: []ast.Expr{
				ast.IntExpr{Value: 123},
				ast.IntExpr{Value: 456},
			},
			Operator: "+",
		},
		binaryExpr,
	)

	expr, err = pargo.Parse(p, newLexer(), "123 + 456 * 789")
	require.NoError(t, err)

	binaryExpr, ok = expr.(ast.BinaryExpr)
	assert.True(ok)
	assert.Equal(
		ast.BinaryExpr{
			Operands: []ast.Expr{
				ast.IntExpr{Value: 123},
				ast.BinaryExpr{
					Operands: []ast.Expr{
						ast.IntExpr{Value: 456},
						ast.IntExpr{Value: 789},
					},
					Operator: "*",
				},
			},
			Operator: "+",
		},
		binaryExpr,
	)
}
