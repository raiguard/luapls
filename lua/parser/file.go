package parser

import (
	"github.com/raiguard/luapls/lua/ast"
)

type File struct {
	Block      ast.Block
	Errors     []ParserError
	LineBreaks []int
	// TODO: Global exports, etc.
}
