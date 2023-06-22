package lexer

import (
	"luapls/lua/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	input := "+-++%^=/==>>=<=<"

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
	}

	l := New(input)
	for _, tok := range tokens {
		realTok := l.NextToken()
		assert.Equal(t, realTok.Type, tok.Type)
		assert.Equal(t, realTok.Literal, tok.Literal)
		assert.Equal(t, realTok.Range.StartCol, tok.Range.StartCol)
		assert.Equal(t, realTok.Range.StartRow, tok.Range.StartRow)
		assert.Equal(t, realTok.Range.EndCol, tok.Range.EndCol)
		assert.Equal(t, realTok.Range.EndRow, tok.Range.EndRow)
	}
}
