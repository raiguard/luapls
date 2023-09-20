package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &Block{
		Stmts: []Statement{&AssignmentStatement{
			Vars: []Expression{
				&Identifier{"foo", 0},
			},
			Exps: []Expression{
				&NumberLiteral{"123", 123, 4},
			},
		}},
		pos: 0,
	}

	require.Equal(t, program.String(), "foo = 123")
}
