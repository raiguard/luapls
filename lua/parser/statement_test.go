package parser

import (
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/stretchr/testify/require"
)

func TestBreakStatement(t *testing.T) {
	testStatement(t, "break", func(stmt ast.BreakStatement) {})
}

func TestGotoStatement(t *testing.T) {
	testStatement(t, "goto continue", func(stmt ast.GotoStatement) {
		require.Equal(t, "continue", stmt.Label.String())
	})
}

func TestIfStatement(t *testing.T) {
	input := `
		if foo then
			foo = 123
			bar = "baz"
		end
	`
	testStatement(t, input, func(stmt ast.IfStatement) {
		lit := requireTypeConversion[ast.Identifier](t, stmt.Condition)
		require.Equal(t, "foo", lit.String())

		consequence := stmt.Consequence
		require.Equal(t, 2, len(consequence.Statements))
	})

	input2 := `
		if 1 + 2 == 4 then
			math_is_true = false
		end
	`
	testStatement(t, input2, func(stmt ast.IfStatement) {
		exp := requireTypeConversion[ast.BinaryExpression](t, stmt.Condition)
		require.Equal(t, "((1 + 2) == 4)", exp.String())

		consequence := stmt.Consequence
		require.Equal(t, 1, len(consequence.Statements))
		consequenceStmt := requireTypeConversion[ast.AssignmentStatement](t, consequence.Statements[0])
		require.Equal(t, "math_is_true", consequenceStmt.Name.String())
		require.Equal(t, "false", consequenceStmt.Value.String())
	})
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
		localStmt := requireTypeConversion[ast.LocalStatement](t, stmts[i])
		assnStmt := requireTypeConversion[ast.AssignmentStatement](t, localStmt.Statement)
		require.Equal(t, test.name, assnStmt.Name.String())
		value := assnStmt.Value
		switch value.(type) {
		case *ast.NumberLiteral:
			require.Equal(t, test.value, value.(*ast.NumberLiteral).Value)
		case *ast.StringLiteral:
			require.Equal(t, test.value, value.(*ast.StringLiteral).Value)
		default:
			require.Fail(t, "Untested token type %s", value.String())
		}
	}
}

func testStatement[T any](t *testing.T, input string, tester func(T)) {
	l := lexer.New(input)
	p := New(l)
	stmt := p.parseStatement()
	checkParserErrors(t, p)
	tester(requireTypeConversion[T](t, stmt))
}
