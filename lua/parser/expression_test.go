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

		require.Equal(t, tt.operator, binaryExp.Token.Literal)

		testNumberLiteral(t, binaryExp.Right, tt.rightValue)
	}
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
