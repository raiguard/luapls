package parser

import (
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/stretchr/testify/require"
)

func TestAssignmentStatement(t *testing.T) {
	testStatement(t, "foo = 123", func(stmt ast.AssignmentStatement) {
		require.Equal(t, 1, len(stmt.Vars))
		require.Equal(t, "foo", stmt.Vars[0].String())
		require.Equal(t, 1, len(stmt.Exps))
		num := requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0])
		require.Equal(t, 123.0, num.Value)
	})
	testStatement(t, "foo, bar = 123, 321", func(stmt ast.AssignmentStatement) {
		require.Equal(t, 2, len(stmt.Vars))
		require.Equal(t, "foo", stmt.Vars[0].String())
		require.Equal(t, "bar", stmt.Vars[1].String())
		require.Equal(t, 2, len(stmt.Exps))
		require.Equal(t, 123.0, requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0]).Value)
		require.Equal(t, 321.0, requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[1]).Value)
	})
}

func TestBreakStatement(t *testing.T) {
	testStatement(t, "break", func(stmt ast.BreakStatement) {})
}

func TestDoStatement(t *testing.T) {
	testStatement(t, "do i = 1 end", func(stmt ast.DoStatement) {
		require.Equal(t, 1, len(stmt.Block.Statements))
	})
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

		block := stmt.Block
		require.Equal(t, 2, len(block.Statements))
	})

	input2 := `
		if 1 + 2 == 4 then
			math_is_true = false
		end
	`
	testStatement(t, input2, func(stmt ast.IfStatement) {
		exp := requireTypeConversion[ast.BinaryExpression](t, stmt.Condition)
		require.Equal(t, "((1 + 2) == 4)", exp.String())

		block := stmt.Block
		require.Equal(t, 1, len(block.Statements))
		blockStmt := requireTypeConversion[ast.AssignmentStatement](t, block.Statements[0])
		require.Equal(t, "math_is_true", blockStmt.Vars[0].String())
		require.Equal(t, "false", blockStmt.Exps[0].String())
	})
}

func TestLocalStatement(t *testing.T) {
	testStatement(t, "local foo = 123", func(stmt ast.LocalStatement) {
		require.Equal(t, 1, len(stmt.Names))
		require.Equal(t, "foo", stmt.Names[0].String())
		require.Equal(t, 1, len(stmt.Exps))
		require.Equal(t, 123.0, requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0]).Value)
	})
	testStatement(t, "local foo, bar = 123, 321", func(stmt ast.LocalStatement) {
		require.Equal(t, 2, len(stmt.Names))
		require.Equal(t, "foo", stmt.Names[0].String())
		require.Equal(t, "bar", stmt.Names[1].String())
		require.Equal(t, 2, len(stmt.Exps))
		require.Equal(t, 123.0, requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[0]).Value)
		require.Equal(t, 321.0, requireTypeConversion[ast.NumberLiteral](t, stmt.Exps[1]).Value)
	})
}

func TestRepeatStatement(t *testing.T) {
	testStatement(t, "repeat i = i + 1 until i == 10", func(stmt ast.RepeatStatement) {
		require.Equal(t, 1, len(stmt.Block.Statements))
		require.Equal(t, "(i == 10)", stmt.Condition.String())
	})
}

func TestWhileStatement(t *testing.T) {
	testStatement(t, "while i < 10 do i = i + 1 end", func(stmt ast.WhileStatement) {
		require.Equal(t, "(i < 10)", stmt.Condition.String())
		require.Equal(t, 1, len(stmt.Block.Statements))
	})
}

func testStatement[T any](t *testing.T, input string, tester func(T)) {
	l := lexer.New(input)
	p := New(l)
	stmt := p.parseStatement()
	checkParserErrors(t, p)
	tester(requireTypeConversion[T](t, stmt))
}
