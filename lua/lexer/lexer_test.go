package lexer

import (
	"luapls/lua/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperators(t *testing.T) {
	input := "+-++%^=/==>>=<=<~="
	tokens := []token.Token{
		{Type: token.PLUS, Literal: "+", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 1, EndRow: 0}},
		{Type: token.MINUS, Literal: "-", Range: token.Range{StartCol: 1, StartRow: 0, EndCol: 2, EndRow: 0}},
		{Type: token.PLUS, Literal: "+", Range: token.Range{StartCol: 2, StartRow: 0, EndCol: 3, EndRow: 0}},
		{Type: token.PLUS, Literal: "+", Range: token.Range{StartCol: 3, StartRow: 0, EndCol: 4, EndRow: 0}},
		{Type: token.PERCENT, Literal: "%", Range: token.Range{StartCol: 4, StartRow: 0, EndCol: 5, EndRow: 0}},
		{Type: token.CARET, Literal: "^", Range: token.Range{StartCol: 5, StartRow: 0, EndCol: 6, EndRow: 0}},
		{Type: token.ASSIGN, Literal: "=", Range: token.Range{StartCol: 6, StartRow: 0, EndCol: 7, EndRow: 0}},
		{Type: token.SLASH, Literal: "/", Range: token.Range{StartCol: 7, StartRow: 0, EndCol: 8, EndRow: 0}},
		{Type: token.EQUAL, Literal: "==", Range: token.Range{StartCol: 8, StartRow: 0, EndCol: 10, EndRow: 0}},
		{Type: token.GT, Literal: ">", Range: token.Range{StartCol: 10, StartRow: 0, EndCol: 11, EndRow: 0}},
		{Type: token.GEQ, Literal: ">=", Range: token.Range{StartCol: 11, StartRow: 0, EndCol: 13, EndRow: 0}},
		{Type: token.LEQ, Literal: "<=", Range: token.Range{StartCol: 13, StartRow: 0, EndCol: 15, EndRow: 0}},
		{Type: token.LT, Literal: "<", Range: token.Range{StartCol: 15, StartRow: 0, EndCol: 16, EndRow: 0}},
		{Type: token.NEQ, Literal: "~=", Range: token.Range{StartCol: 16, StartRow: 0, EndCol: 18, EndRow: 0}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 18, StartRow: 0, EndCol: 18, EndRow: 0}},
	}
	testLexer(t, input, tokens)
}

func TestKeywords(t *testing.T) {
	input := "local while for"
	tokens := []token.Token{
		{Type: token.LOCAL, Literal: "local", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 5, EndRow: 0}},
		{Type: token.WHILE, Literal: "while", Range: token.Range{StartCol: 6, StartRow: 0, EndCol: 11, EndRow: 0}},
		{Type: token.FOR, Literal: "for", Range: token.Range{StartCol: 12, StartRow: 0, EndCol: 15, EndRow: 0}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 15, StartRow: 0, EndCol: 15, EndRow: 0}},
	}
	testLexer(t, input, tokens)
}

func TestAssignment(t *testing.T) {
	input := `
local foo = 123
2 + 2 == 4
`
	tokens := []token.Token{
		{Type: token.LOCAL, Literal: "local", Range: token.Range{StartCol: 0, StartRow: 1, EndCol: 5, EndRow: 1}},
		{Type: token.IDENT, Literal: "foo", Range: token.Range{StartCol: 6, StartRow: 1, EndCol: 9, EndRow: 1}},
		{Type: token.ASSIGN, Literal: "=", Range: token.Range{StartCol: 10, StartRow: 1, EndCol: 11, EndRow: 1}},
		{Type: token.NUMBER, Literal: "123", Range: token.Range{StartCol: 12, StartRow: 1, EndCol: 15, EndRow: 1}},
		{Type: token.NUMBER, Literal: "2", Range: token.Range{StartCol: 0, StartRow: 2, EndCol: 1, EndRow: 2}},
		{Type: token.PLUS, Literal: "+", Range: token.Range{StartCol: 2, StartRow: 2, EndCol: 3, EndRow: 2}},
		{Type: token.NUMBER, Literal: "2", Range: token.Range{StartCol: 4, StartRow: 2, EndCol: 5, EndRow: 2}},
		{Type: token.EQUAL, Literal: "==", Range: token.Range{StartCol: 6, StartRow: 2, EndCol: 8, EndRow: 2}},
		{Type: token.NUMBER, Literal: "4", Range: token.Range{StartCol: 9, StartRow: 2, EndCol: 10, EndRow: 2}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 0, StartRow: 3, EndCol: 0, EndRow: 3}},
	}
	testLexer(t, input, tokens)
}

func TestNumbers(t *testing.T) {
	input := "3 3.0 3.1416 314.16e-2 0.31416E1 0xff 0x0.1E 0xA23p-4 0X1.921FB54442D18P+1"
	tokens := []token.Token{
		{Type: token.NUMBER, Literal: "3", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 1, EndRow: 0}},
		{Type: token.NUMBER, Literal: "3.0", Range: token.Range{StartCol: 2, StartRow: 0, EndCol: 5, EndRow: 0}},
		{Type: token.NUMBER, Literal: "3.1416", Range: token.Range{StartCol: 6, StartRow: 0, EndCol: 12, EndRow: 0}},
		{Type: token.NUMBER, Literal: "314.16e-2", Range: token.Range{StartCol: 13, StartRow: 0, EndCol: 22, EndRow: 0}},
		{Type: token.NUMBER, Literal: "0.31416E1", Range: token.Range{StartCol: 23, StartRow: 0, EndCol: 32, EndRow: 0}},
		{Type: token.NUMBER, Literal: "0xff", Range: token.Range{StartCol: 33, StartRow: 0, EndCol: 37, EndRow: 0}},
		{Type: token.NUMBER, Literal: "0x0.1E", Range: token.Range{StartCol: 38, StartRow: 0, EndCol: 44, EndRow: 0}},
		{Type: token.NUMBER, Literal: "0xA23p-4", Range: token.Range{StartCol: 45, StartRow: 0, EndCol: 53, EndRow: 0}},
		{Type: token.NUMBER, Literal: "0X1.921FB54442D18P+1", Range: token.Range{StartCol: 54, StartRow: 0, EndCol: 74, EndRow: 0}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 74, StartRow: 0, EndCol: 74, EndRow: 0}},
	}
	testLexer(t, input, tokens)
}

func TestStrings(t *testing.T) {
	input := "'314.16e-2' \"foo bar\""
	tokens := []token.Token{
		{Type: token.STRING, Literal: "'314.16e-2'", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 11, EndRow: 0}},
		{Type: token.STRING, Literal: "\"foo bar\"", Range: token.Range{StartCol: 12, StartRow: 0, EndCol: 21, EndRow: 0}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 21, StartRow: 0, EndCol: 21, EndRow: 0}},
	}
	testLexer(t, input, tokens)
}

func TestRawStrings(t *testing.T) {
	input := `[====[asdkflasdkfjs]===]
  19"38'44'"al]====]`
	tokens := []token.Token{
		{Type: token.RAWSTRING, Literal: "[====[asdkflasdkfjs]===]\n  19\"38'44'\"al]====]", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 20, EndRow: 1}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 20, StartRow: 1, EndCol: 20, EndRow: 1}},
	}
	testLexer(t, input, tokens)
}

func TestLabel(t *testing.T) {
	input := "::foo_BAR_314::"
	tokens := []token.Token{
		{Type: token.LABEL, Literal: "::", Range: token.Range{StartCol: 0, StartRow: 0, EndCol: 2, EndRow: 0}},
		{Type: token.IDENT, Literal: "foo_BAR_314", Range: token.Range{StartCol: 2, StartRow: 0, EndCol: 13, EndRow: 0}},
		{Type: token.LABEL, Literal: "::", Range: token.Range{StartCol: 13, StartRow: 0, EndCol: 15, EndRow: 0}},
		{Type: token.EOF, Literal: "", Range: token.Range{StartCol: 15, StartRow: 0, EndCol: 15, EndRow: 0}},
	}
	testLexer(t, input, tokens)
}

func testLexer(t *testing.T, input string, tokens []token.Token) {
	l := New(input)
	for _, expected := range tokens {
		actual := l.NextToken()
		assert.Equal(t, expected.Type, actual.Type)
		assert.Equal(t, expected.Literal, actual.Literal)
		assert.Equal(t, expected.Range.StartCol, actual.Range.StartCol)
		assert.Equal(t, expected.Range.StartRow, actual.Range.StartRow)
		assert.Equal(t, expected.Range.EndCol, actual.Range.EndCol)
		assert.Equal(t, expected.Range.EndRow, actual.Range.EndRow)
	}
}
