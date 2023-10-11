package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func getInnermostNode(n ast.Node, pos token.Pos) ast.Node {
	var node ast.Node
	ast.Walk(n, func(n ast.Node) bool {
		if n.Pos() <= pos && pos < n.End() {
			node = n
			return true
		}
		return false
	})
	return node
}

func parseFile(ctx *glsp.Context, filename, src string) {
	file := parser.New(src).ParseFile()
	files[filename] = &file
}

func logToEditor(ctx *glsp.Context, format string, args ...any) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: fmt.Sprintf(format, args...)},
	)
}
