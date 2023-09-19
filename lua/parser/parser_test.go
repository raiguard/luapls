package parser

import (
	"fmt"
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"

	"github.com/stretchr/testify/require"
)

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range p.errors {
		t.Errorf("parser error: %q", msg)
	}
	t.Fail()
}

func testNumberLiteral(t *testing.T, il ast.Expression, value float64) {
	num := requireTypeConversion[ast.NumberLiteral](t, il)
	require.Equal(t, value, num.Value)
	require.Equal(t, fmt.Sprintf("%.0f", value), num.String())
}

type statementTest struct {
	input, expected string
}

func testStatements(t *testing.T, tests []statementTest) {
	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		block := p.ParseBlock()
		require.NotNil(t, block)
		require.Equal(t, 1, len(block))

		require.Equal(t, test.expected, block.String())
	}
}

func requireTypeConversion[T any](t *testing.T, val any) T {
	res, ok := val.(*T)
	require.True(t, ok)
	return *res
}
