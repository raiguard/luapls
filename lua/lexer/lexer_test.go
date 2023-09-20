package lexer

import (
	"github.com/raiguard/luapls/lua/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperators(t *testing.T) {
	input := "+-++%^=/==>>=<=<~="
	tokens := []token.Token{
		{Type: token.PLUS, Literal: "+", Pos: 0},
		{Type: token.MINUS, Literal: "-", Pos: 1},
		{Type: token.PLUS, Literal: "+", Pos: 2},
		{Type: token.PLUS, Literal: "+", Pos: 3},
		{Type: token.PERCENT, Literal: "%", Pos: 4},
		{Type: token.POW, Literal: "^", Pos: 5},
		{Type: token.ASSIGN, Literal: "=", Pos: 6},
		{Type: token.SLASH, Literal: "/", Pos: 7},
		{Type: token.EQUAL, Literal: "==", Pos: 8},
		{Type: token.GT, Literal: ">", Pos: 10},
		{Type: token.GEQ, Literal: ">=", Pos: 11},
		{Type: token.LEQ, Literal: "<=", Pos: 13},
		{Type: token.LT, Literal: "<", Pos: 15},
		{Type: token.NEQ, Literal: "~=", Pos: 16},
		{Type: token.EOF, Literal: "", Pos: 18},
	}
	testLexer(t, input, tokens)
}

func TestKeywords(t *testing.T) {
	input := "local while for"
	tokens := []token.Token{
		{Type: token.LOCAL, Literal: "local", Pos: 0},
		{Type: token.WHILE, Literal: "while", Pos: 6},
		{Type: token.FOR, Literal: "for", Pos: 12},
		{Type: token.EOF, Literal: "", Pos: 15},
	}
	testLexer(t, input, tokens)
}

func TestAssignment(t *testing.T) {
	input := `
local foo = 123
2 + 2 == 4`
	tokens := []token.Token{
		{Type: token.LOCAL, Literal: "local", Pos: 1},
		{Type: token.IDENT, Literal: "foo", Pos: 7},
		{Type: token.ASSIGN, Literal: "=", Pos: 11},
		{Type: token.NUMBER, Literal: "123", Pos: 13},
		{Type: token.NUMBER, Literal: "2", Pos: 17},
		{Type: token.PLUS, Literal: "+", Pos: 19},
		{Type: token.NUMBER, Literal: "2", Pos: 21},
		{Type: token.EQUAL, Literal: "==", Pos: 23},
		{Type: token.NUMBER, Literal: "4", Pos: 26},
		{Type: token.EOF, Literal: "", Pos: 27},
	}
	testLexer(t, input, tokens)
}

func TestNumbers(t *testing.T) {
	input := "3 3.0 3.1416 314.16e-2 0.31416E1 0xff 0x0.1E 0xA23p-4 0X1.921FB54442D18P+1"
	tokens := []token.Token{
		{Type: token.NUMBER, Literal: "3", Pos: 0},
		{Type: token.NUMBER, Literal: "3.0", Pos: 2},
		{Type: token.NUMBER, Literal: "3.1416", Pos: 6},
		{Type: token.NUMBER, Literal: "314.16e-2", Pos: 13},
		{Type: token.NUMBER, Literal: "0.31416E1", Pos: 23},
		{Type: token.NUMBER, Literal: "0xff", Pos: 33},
		{Type: token.NUMBER, Literal: "0x0.1E", Pos: 38},
		{Type: token.NUMBER, Literal: "0xA23p-4", Pos: 45},
		{Type: token.NUMBER, Literal: "0X1.921FB54442D18P+1", Pos: 54},
		{Type: token.EOF, Literal: "", Pos: 74},
	}
	testLexer(t, input, tokens)
}

func TestStrings(t *testing.T) {
	input := "'314.16e-2' \"foo bar\""
	tokens := []token.Token{
		{Type: token.STRING, Literal: "'314.16e-2'", Pos: 0},
		{Type: token.STRING, Literal: "\"foo bar\"", Pos: 12},
		{Type: token.EOF, Literal: "", Pos: 21},
	}
	testLexer(t, input, tokens)
}

func TestRawStrings(t *testing.T) {
	input := `[====[asdkflasdkfjs]===]
  19"38'44'"al]====]
[===[[[==[=]]==]====][=[]=]===]`
	tokens := []token.Token{
		{Type: token.RAWSTRING, Literal: "[====[asdkflasdkfjs]===]\n  19\"38'44'\"al]====]", Pos: 0},
		{Type: token.RAWSTRING, Literal: "[===[[[==[=]]==]====][=[]=]===]", Pos: 46},
		{Type: token.EOF, Literal: "", Pos: 77},
	}
	testLexer(t, input, tokens)
}

func TestLabel(t *testing.T) {
	input := "::foo_BAR_314::"
	tokens := []token.Token{
		{Type: token.LABEL, Literal: "::", Pos: 0},
		{Type: token.IDENT, Literal: "foo_BAR_314", Pos: 2},
		{Type: token.LABEL, Literal: "::", Pos: 13},
		{Type: token.EOF, Literal: "", Pos: 15},
	}
	testLexer(t, input, tokens)
}

func testLexer(t *testing.T, input string, tokens []token.Token) {
	l := New(input)
	for _, expected := range tokens {
		actual := l.NextToken()
		// Compare strings for better test output
		assert.Equal(t, expected.Type.String(), actual.Type.String())
		assert.Equal(t, expected.Literal, actual.Literal)
		assert.Equal(t, expected.Pos, actual.Pos)
	}
}
