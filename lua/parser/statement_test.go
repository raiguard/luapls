package parser

import (
	"testing"

	"github.com/raiguard/luapls/lua/ast"
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
		require.Equal(t, 1, len(stmt.Body.Stmts))
	})
}

func TestForStatement(t *testing.T) {
	testStatement(t, "for i = 1, 100 do j = j + i end", func(stmt ast.ForStatement) {
		require.Equal(t, "i", stmt.Name.String())
		require.Equal(t, "1", stmt.Start.String())
		require.Equal(t, "100", stmt.Finish.String())
		require.Nil(t, stmt.Step)
		require.Equal(t, 1, len(stmt.Body.Stmts))
	})
}

func TestForInStatement(t *testing.T) {
	testStatement(t, "for key, value in tbl do j = j + i end", func(stmt ast.ForInStatement) {
		require.Equal(t, 2, len(stmt.Names))
		require.Equal(t, "key", stmt.Names[0].String())
		require.Equal(t, "value", stmt.Names[1].String())
		require.Equal(t, 1, len(stmt.Exps))
		require.Equal(t, "tbl", stmt.Exps[0].String())
		require.Equal(t, 1, len(stmt.Body.Stmts))
	})
}

func TestFunctionCallStatement(t *testing.T) {
	testStatement(t, "foo(1, 2)", func(stmt ast.ExpressionStatement) {
		fc, ok := stmt.Exp.(*ast.FunctionCall)
		require.True(t, ok)
		require.Equal(t, "foo", fc.Left.String())
		require.Equal(t, 2, len(fc.Args))
		require.Equal(t, "1", fc.Args[0].String())
		require.Equal(t, "2", fc.Args[1].String())
	})
}

func TestFunctionStatement(t *testing.T) {
	testStatement(t, "function foo(key, value) end", func(stmt ast.FunctionStatement) {
		require.Equal(t, "foo", stmt.Left.String())
		require.Equal(t, 2, len(stmt.Params))
		require.Equal(t, "key", stmt.Params[0].String())
		require.Equal(t, "value", stmt.Params[1].String())
		require.Equal(t, 0, len(stmt.Body.Stmts))
		require.Equal(t, false, stmt.IsLocal)
	})
}

func TestGotoStatement(t *testing.T) {
	testStatement(t, "goto continue", func(stmt ast.GotoStatement) {
		require.Equal(t, "continue", stmt.Name.String())
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
		require.Equal(t, 1, len(stmt.Clauses))
		clause := stmt.Clauses[0]
		lit := requireTypeConversion[ast.Identifier](t, clause.Condition)
		require.Equal(t, "foo", lit.String())

		block := clause.Body
		require.Equal(t, 2, len(block.Stmts))
	})

	input2 := `
		if 1 + 2 == 4 then
			math_is_true = false
		elseif 2 + 2 == 4 then
			math_is_true = true
		end
	`
	testStatement(t, input2, func(stmt ast.IfStatement) {
		require.Equal(t, 2, len(stmt.Clauses))
		require.Equal(t, 1, len(stmt.Clauses[0].Body.Stmts))
		require.Equal(t, 1, len(stmt.Clauses[1].Body.Stmts))
	})
}

func TestLabelStatement(t *testing.T) {
	testStatement(t, "::continue::", func(stmt ast.LabelStatement) {
		require.Equal(t, "continue", stmt.Name.Literal)
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
		require.Equal(t, 1, len(stmt.Body.Stmts))
		require.Equal(t, "(i == 10)", stmt.Condition.String())
	})
}

func TestReturnStatement(t *testing.T) {
	testStatement(t, "return foo, bar", func(stmt ast.ReturnStatement) {
		require.Equal(t, 2, len(stmt.Exps))
		require.Equal(t, "foo", stmt.Exps[0].String())
		require.Equal(t, "bar", stmt.Exps[1].String())
	})
}

func TestWhileStatement(t *testing.T) {
	testStatement(t, "while i < 10 do i = i + 1 end", func(stmt ast.WhileStatement) {
		require.Equal(t, "(i < 10)", stmt.Condition.String())
		require.Equal(t, 1, len(stmt.Body.Stmts))
	})
}

func testStatement[T any](t *testing.T, input string, tester func(T)) {
	p := New(input)
	stmt := p.parseStatement()
	checkParserErrors(t, p)
	tester(requireTypeConversion[T](t, stmt))
}
