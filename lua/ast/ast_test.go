package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &Block{
		&AssignmentStatement{
			Vars: []Expression{
				&Identifier{"foo"},
			},
			Exps: []Expression{
				&NumberLiteral{"123", 123},
			},
		},
	}

	require.Equal(t, program.String(), "foo = 123")
}
