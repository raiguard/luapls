package parser

import (
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/stretchr/testify/require"
)

func TestAssignmentStatement(t *testing.T) {
	testStatement(t, "foo = 123", func(stmt ast.AssignmentStatement) {
		require.Equal(t, "foo", stmt.Name.String())
		require.Equal(t, 1, len(stmt.Exps))
		num := requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0])
		require.Equal(t, 123.0, num.Value)
	})
	// TODO: This is invalid Lua, needs two variables
	testStatement(t, "foo = 123, 321", func(stmt ast.AssignmentStatement) {
		require.Equal(t, "foo", stmt.Name.String())
		require.Equal(t, 2, len(stmt.Exps))
		num1 := requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0])
		require.Equal(t, 123.0, num1.Value)
		num2 := requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[1])
		require.Equal(t, 321.0, num2.Value)
	})
}

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
		require.Equal(t, "false", consequenceStmt.Exps[0].String())
	})
}

func TestLocalStatement(t *testing.T) {
	testStatement(t, "local foo = 123", func(stmt ast.LocalStatement) {
		assnStmt := requireTypeConversion[ast.AssignmentStatement](t, stmt.Statement)
		require.Equal(t, "foo", assnStmt.Name.String())
		require.Equal(t, 123.0, requireTypeConversion[ast.NumberLiteral](t, assnStmt.Exps[0]).Value)
	})
	testStatement(t, "local foo = 'lorem ipsum'", func(stmt ast.LocalStatement) {
		assnStmt := requireTypeConversion[ast.AssignmentStatement](t, stmt.Statement)
		require.Equal(t, "foo", assnStmt.Name.String())
		require.Equal(t, "lorem ipsum", requireTypeConversion[ast.StringLiteral](t, assnStmt.Exps[0]).Value)
	})
}

func testStatement[T any](t *testing.T, input string, tester func(T)) {
	l := lexer.New(input)
	p := New(l)
	stmt := p.parseStatement()
	checkParserErrors(t, p)
	tester(requireTypeConversion[T](t, stmt))
}
