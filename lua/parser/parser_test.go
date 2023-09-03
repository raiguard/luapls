package parser

import (
	"fmt"
	"testing"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"

	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	input := `
		foo = 123
		bar = 456
	`

	l := lexer.New(input)
	p := New(l)

	block := p.ParseBlock()
	stmts := block.Statements

	tests := []struct {
		name string
		num  float64
	}{
		{"foo", 123},
		{"bar", 456},
	}

	require.Equal(t, len(tests), len(stmts))
	for i, test := range tests {
		assnStmt := requireTypeConversion[ast.AssignmentStatement](t, stmts[i])
		require.Equal(t, test.name, assnStmt.Name.String())
		lit := requireTypeConversion[ast.NumberLiteral](t, assnStmt.Value)
		require.Equal(t, test.num, lit.Value)
	}
}

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
	integ := requireTypeConversion[ast.NumberLiteral](t, il)
	require.Equal(t, value, integ.Value)
	require.Equal(t, fmt.Sprintf("%.0f", value), integ.Token.Literal)
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
		require.Equal(t, 1, len(block.Statements))

		require.Equal(t, test.expected, block.String())
	}
}

func requireTypeConversion[T any](t *testing.T, val any) T {
	res, ok := val.(*T)
	if !ok {
		fmt.Println("Failed")
	}
	require.True(t, ok)
	return *res
}
