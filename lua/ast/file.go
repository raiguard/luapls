package ast

import "github.com/raiguard/luapls/lua/token"

type File struct {
	Block      Block
	Errors     []Error
	LineBreaks token.LineBreaks
}
