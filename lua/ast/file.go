package ast

import (
	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type File struct {
	Block       *Block
	Diagnostics []Diagnostic
	LineBreaks  token.LineBreaks
	URI         protocol.URI
}
