package parser

// import (
// 	"testing"

// 	"github.com/raiguard/luapls/lua/ast"
// 	"github.com/raiguard/luapls/lua/token"
// 	"github.com/stretchr/testify/require"
// )

// func TestInfixExpression(t *testing.T) {
// 	infixTests := []struct {
// 		input    string
// 		left     float64
// 		operator token.TokenType
// 		right    float64
// 	}{
// 		{"5 + 5", 5, token.PLUS, 5},
// 		{"5 - 5", 5, token.MINUS, 5},
// 		{"5 * 5", 5, token.MUL, 5},
// 		{"5 / 5", 5, token.SLASH, 5},
// 		{"5 > 5", 5, token.GT, 5},
// 		{"5 < 5", 5, token.LT, 5},
// 		{"5 == 5", 5, token.EQUAL, 5},
// 		{"5 ~= 5", 5, token.NEQ, 5},
// 	}

// 	for _, tt := range infixTests {
// 		testExpression(t, tt.input, func(exp ast.InfixExpression) {
// 			require.Equal(t, tt.left, requireTypeConversion[ast.NumberLiteral](t, exp.Left).Value)
// 			require.Equal(t, tt.operator, exp.Operator)
// 			require.Equal(t, tt.right, requireTypeConversion[ast.NumberLiteral](t, exp.Right).Value)
// 		})
// 	}
// }

// func TestBooleanLiteral(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected bool
// 	}{
// 		{"false", false},
// 		{"true", true},
// 	}
// 	for _, test := range tests {
// 		p := New(test.input)
// 		lit := p.parseBooleanLiteral()
// 		checkParserErrors(t, p)
// 		require.Equal(t, test.expected, lit.Value)
// 	}
// }

// func TestFunctionExpression(t *testing.T) {
// 	testExpression(t, "function(a, b) print(a + b) end", func(exp ast.FunctionExpression) {
// 		require.Equal(t, 2, len(exp.Params))
// 		require.Equal(t, "a", exp.Params[0].String())
// 		require.Equal(t, "b", exp.Params[1].String())
// 		require.Equal(t, 1, len(exp.Body.Stmts))
// 	})
// }

// func TestFunctionCall(t *testing.T) {
// 	testExpression(t, "print(foo)", func(exp ast.FunctionCall) {
// 		require.Equal(t, "print", exp.Left.String())
// 		require.Equal(t, 1, len(exp.Args))
// 		require.Equal(t, "foo", exp.Args[0].String())
// 	})
// 	testExpression(t, "print(foo)(bar)", func(exp ast.FunctionCall) {
// 		require.Equal(t, "print(foo)", exp.Left.String())
// 		require.Equal(t, 1, len(exp.Args))
// 		require.Equal(t, "bar", exp.Args[0].String())
// 	})
// }

// func TestOperatorPrecedence(t *testing.T) {
// 	testStatements(t, []statementTest{
// 		{"i = 1 + 2 - 3 * -4 / 5 % 6 ^ 7 .. 8", "i = (((1 + 2) - (((3 * (-4)) / 5) % (6 ^ 7))) .. 8)"},
// 		{"i = 2 + 2 + 2", "i = ((2 + 2) + 2)"},
// 		{"i = 2 ^ 2 ^ 2", "i = (2 ^ (2 ^ 2))"},
// 		{"i = 2 .. 2 .. 2", "i = (2 .. (2 .. 2))"},
// 	})
// }

// func TestTableLiteral(t *testing.T) {
// 	input := "{ 'bar', baz, 2 + 2, foo = lorem, ['12345-54321'] = baz, }"
// 	p := New(input)
// 	tbl := p.parseTableLiteral()
// 	checkParserErrors(t, p)
// 	require.NotNil(t, tbl)
// 	require.Equal(t, 5, len(tbl.Fields))
// }

// func TestIndexExpression(t *testing.T) {
// 	testExpression(t, "foo.bar[baz]", func(exp ast.IndexExpression) {
// 		require.Equal(t, "foo.bar", exp.Left.String())
// 		require.Equal(t, "baz", exp.Inner.String())
// 		require.Equal(t, true, exp.IsBrackets)
// 		require.Equal(t, false, exp.Left.(*ast.IndexExpression).IsBrackets)
// 	})
// 	testExpression(t, "foo[bar][baz]", func(exp ast.IndexExpression) {
// 		require.Equal(t, "foo[bar]", exp.Left.String())
// 		require.Equal(t, "baz", exp.Inner.String())
// 		require.Equal(t, true, exp.IsBrackets)
// 		require.Equal(t, true, exp.Left.(*ast.IndexExpression).IsBrackets)
// 	})
// }

// func testExpression[T any](t *testing.T, input string, tester func(T)) {
// 	p := New(input)
// 	stmt := p.parseExpression(LOWEST, true)
// 	checkParserErrors(t, p)
// 	tester(requireTypeConversion[T](t, stmt))
// }
