package parser

import (
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/stretchr/testify/require"
)

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

		binaryExp := requireTypeConversion[ast.BinaryExpression](t, exp)

		testNumberLiteral(t, binaryExp.Left, tt.leftValue)

		require.Equal(t, tt.operator, binaryExp.Operator.String())

		testNumberLiteral(t, binaryExp.Right, tt.rightValue)
	}
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"false", false},
		{"true", true},
	}
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		lit := p.parseBooleanLiteral()
		checkParserErrors(t, p)
		require.Equal(t, test.expected, lit.Value)
	}
}

func TestFunctionExpression(t *testing.T) {
	testExpression(t, "function(a, b) print(a + b) end", func(exp ast.FunctionExpression) {
		require.Equal(t, 2, len(exp.Params))
		require.Equal(t, "a", exp.Params[0].String())
		require.Equal(t, "b", exp.Params[1].String())
		require.Equal(t, 1, len(exp.Body))
	})
}

func TestOperatorPrecedence(t *testing.T) {
	// TODO: Directly check parser output instead of String() output
	testStatements(t, []statementTest{
		{"i = 1 + 2 - 3 * -4 / 5 % 6 ^ 7 .. 8", "i = (((1 + 2) - (((3 * (-4)) / 5) % (6 ^ 7))) .. 8)"},
		{"i = 2 + 2 + 2", "i = ((2 + 2) + 2)"},
		{"i = 2 ^ 2 ^ 2", "i = (2 ^ (2 ^ 2))"},
		{"i = 2 .. 2 .. 2", "i = (2 .. (2 .. 2))"},
	})
}

func TestTableLiteral(t *testing.T) {
	input := "{ 'bar', baz, 2 + 2, foo = lorem, ['12345-54321'] = baz, }"
	p := New(lexer.New(input))
	tbl := p.parseTableLiteral()
	checkParserErrors(t, p)
	require.NotNil(t, tbl)
	require.Equal(t, 5, len(tbl.Fields))
}

func testExpression[T any](t *testing.T, input string, tester func(T)) {
	l := lexer.New(input)
	p := New(l)
	stmt := p.parseExpression(LOWEST)
	checkParserErrors(t, p)
	tester(requireTypeConversion[T](t, stmt))
}
