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
		assnStmt, ok := stmts[i].(*ast.AssignmentStatement)
		require.True(t, ok)
		require.Equal(t, test.name, assnStmt.Name.String())
		lit, ok := assnStmt.Value.(*ast.NumberLiteral)
		require.True(t, ok)
		require.Equal(t, test.num, lit.Value)
	}
}

func TestLocalStatement(t *testing.T) {
	input := `
		local foo = 123
		local bar = 456
		local baz = "lorem ipsum"
	`

	l := lexer.New(input)
	p := New(l)

	block := p.ParseBlock()
	stmts := block.Statements

	tests := []struct {
		name  string
		value any
	}{
		{"foo", 123.0},
		{"bar", 456.0},
		{"baz", "lorem ipsum"},
	}

	require.Equal(t, len(tests), len(stmts))

	for i, test := range tests {
		localStmt, ok := stmts[i].(*ast.LocalStatement)
		require.True(t, ok)
		assnStmt, ok := localStmt.Statement.(*ast.AssignmentStatement)
		require.True(t, ok)
		require.Equal(t, test.name, assnStmt.Name.String())
		value := assnStmt.Value
		switch value.(type) {
		case *ast.NumberLiteral:
			require.Equal(t, test.value, value.(*ast.NumberLiteral).Value)
		case *ast.StringLiteral:
			require.Equal(t, test.value, value.(*ast.StringLiteral).Value)
		default:
			require.FailNow(t, "Untested token type %s", value.String())
		}
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
