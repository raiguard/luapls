package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"testing"

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
	statements := block.Statements

	tests := []struct {
		name string
		num  float64
	}{
		{"foo", 123},
		{"bar", 456},
	}

	require.Equal(t, len(tests), len(statements))
	for i, test := range tests {
		assnStmt, ok := statements[i].(*ast.AssignmentStatement)
		require.True(t, ok)
		require.Equal(t, test.name, assnStmt.Name.String())
		lit, ok := assnStmt.Value.(*ast.NumberLiteral)
		require.True(t, ok)
		require.Equal(t, test.num, lit.Value)
	}
}

// func testNodeList(t *testing.T, expected []ast.Node, actual []ast.Node) {
// 	require.Equal(t, len(expected), len(actual))
// 	for i := 0; i < len(expected); i += 1 {
// 		expected := expected[i]
// 		actual := actual[i]
// 		require.Equal(t, expected.Pos
// 	}
// }
