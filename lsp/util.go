package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

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

// getLocals returns a list of all local variables contained in node for the given pos.
func getLocals(node ast.Node, pos token.Pos) map[string]*ast.Identifier {
	locals := map[string]*ast.Identifier{}

	ast.Walk(node, func(node ast.Node) bool {
		// If the entire node is before the current position
		if pos >= node.Pos() && pos > node.End() {
			switch node := node.(type) {
			case *ast.LocalStatement:
				for _, ident := range node.Names {
					locals[ident.Literal] = ident
				}
			case *ast.FunctionStatement:
				if ident, ok := node.Left.(*ast.Identifier); ok {
					locals[ident.Literal] = ident
				}
			}
			return false
		}
		// If the node does not intersect the position
		if pos < node.Pos() || pos >= node.End() {
			return false
		}
		switch node := node.(type) {
		case *ast.ForInStatement:
			for _, ident := range node.Names {
				locals[ident.Literal] = ident
			}
		case *ast.ForStatement:
			if node.Name != nil {
				locals[node.Name.Literal] = node.Name
			}
		case *ast.FunctionExpression:
			for _, ident := range node.Params {
				locals[ident.Literal] = ident
			}
		case *ast.FunctionStatement:
			for _, ident := range node.Params {
				locals[ident.Literal] = ident
			}
		default:
			// If the node contains the position
			return node.Pos() <= pos && pos < node.End()
		}

		return true
	})

	return locals
}
