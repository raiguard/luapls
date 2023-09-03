package parser

import (
	"fmt"
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"

	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	input := `
		foo = 123
		bar = 456
	`

	l := lexer.New(input)
	p := New(l)

	block := p.ParseBlock()
	stmts := block.Statements

	tests := []struct {
		name string
		num  float64
	}{
		{"foo", 123},
		{"bar", 456},
	}

	require.Equal(t, len(tests), len(stmts))
	for i, test := range tests {
		assnStmt, ok := stmts[i].(*ast.AssignmentStatement)
		require.True(t, ok)
		require.Equal(t, test.name, assnStmt.Name.String())
		lit, ok := assnStmt.Value.(*ast.NumberLiteral)
		require.True(t, ok)
		require.Equal(t, test.num, lit.Value)
	}
}

func TestLocalStatement(t *testing.T) {
	input := `
		local foo = 123
		local bar = 456
		local baz = "lorem ipsum"
		local complex = "\"dolor sit amet"
	`

	l := lexer.New(input)
	p := New(l)

	block := p.ParseBlock()
	stmts := block.Statements

	tests := []struct {
		name  string
		value any
	}{
		{"foo", 123.0},
		{"bar", 456.0},
		{"baz", "lorem ipsum"},
		{"complex", "\\\"dolor sit amet"},
	}

	require.Equal(t, len(tests), len(stmts))

	for i, test := range tests {
		localStmt, ok := stmts[i].(*ast.LocalStatement)
		require.True(t, ok)
		assnStmt, ok := localStmt.Statement.(*ast.AssignmentStatement)
		require.True(t, ok)
		require.Equal(t, test.name, assnStmt.Name.String())
		value := assnStmt.Value
		switch value.(type) {
		case *ast.NumberLiteral:
			require.Equal(t, test.value, value.(*ast.NumberLiteral).Value)
		case *ast.StringLiteral:
			require.Equal(t, test.value, value.(*ast.StringLiteral).Value)
		default:
			require.FailNow(t, "Untested token type %s", value.String())
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := `
		if foo then
			foo = 123
			bar = "baz"
		end
	`

	l := lexer.New(input)
	p := New(l)

	block := p.ParseBlock()
	stmts := block.Statements

	require.Equal(t, 1, len(stmts))

	ifStmt, ok := stmts[0].(*ast.IfStatement)
	require.True(t, ok)

	lit, ok := ifStmt.Condition.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, "foo", lit.String())

	consequence := ifStmt.Consequence
	require.Equal(t, 2, len(consequence.Statements))
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  float64
		operator   string
		rightValue float64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 ~= 5", 5, "~=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		exp := p.parseExpression(LOWEST)
		checkParserErrors(t, p)

		infix, ok := exp.(*ast.InfixExpression)
		require.True(t, ok)

		testNumberLiteral(t, infix.Left, tt.leftValue)

		require.Equal(t, tt.operator, infix.Operator)

		testNumberLiteral(t, infix.Right, tt.rightValue)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range p.errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testNumberLiteral(t *testing.T, il ast.Expression, value float64) {
	integ, ok := il.(*ast.NumberLiteral)
	require.True(t, ok)
	require.Equal(t, value, integ.Value)
	require.Equal(t, fmt.Sprintf("%.0f", value), integ.Token.Literal)
}

// func testNodeList(t *testing.T, expected []ast.Node, actual []ast.Node) {
// 	require.Equal(t, len(expected), len(actual))
// 	for i := 0; i < len(expected); i += 1 {
// 		expected := expected[i]
// 		actual := actual[i]
// 		require.Equal(t, expected.Pos
// 	}
// }
