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
				Vars: []Identifier{
					{
						Type:    token.IDENT,
						Literal: "foo",
						Range:   token.Range{StartCol: 6, StartRow: 0, EndCol: 9, EndRow: 3},
					},
				},
				Exps: []Expression{
					&NumberLiteral{123},
				},
			},
		},
	}

	require.Equal(t, program.String(), "foo = 123")
}
