package ast

import (
	"fmt"

	"github.com/raiguard/luapls/lua/token"
)

type Error struct {
	Message string
	Range   token.Range
}

func (pe *Error) String() string {
	return fmt.Sprintf("%s: %s", &pe.Range, pe.Message)
}
