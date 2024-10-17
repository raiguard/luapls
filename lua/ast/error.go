package ast

import (
	"fmt"

	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Error struct {
	Message  string
	Range    token.Range
	Severity protocol.DiagnosticSeverity
}

func (pe *Error) String() string {
	return fmt.Sprintf("%s: %s", &pe.Range, pe.Message)
}
