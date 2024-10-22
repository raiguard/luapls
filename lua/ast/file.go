package ast

import "github.com/raiguard/luapls/lua/token"

type File struct {
	Block       Block
	Diagnostics []Diagnostic
	LineBreaks  token.LineBreaks
}
