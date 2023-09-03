package ast

import (
	"testing"

	"github.com/raiguard/luapls/lua/token"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &Block{
		Statements: []Statement{
			&AssignmentStatement{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foo",
					Range:   token.Range{StartCol: 0, StartRow: 0, EndCol: 3, EndRow: 0},
				},
				Vars: []Identifier{
					{
						Type:    token.IDENT,
						Literal: "foo",
						Range:   token.Range{StartCol: 6, StartRow: 0, EndCol: 9, EndRow: 3},
					},
				},
				Exps: []Expression{
					&NumberLiteral{
						Token: token.Token{
							Type:    token.NUMBER,
							Literal: "123",
							Range:   token.Range{},
						},
						Value: 123,
					},
				},
			},
		},
	}

	require.Equal(t, program.String(), "foo = 123")
}
