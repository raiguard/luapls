package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Get type of node at cursor and show it
// - Find most recent definition of the variable
// - Resolve the type of that expression

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	node, parents := ast.GetNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	// TODO: Literals
	ident, _ := node.(*ast.Identifier)
	if ident == nil {
		return nil, nil
	}

	def := getDefinition(node, parents)
	if def == nil {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: fmt.Sprintf(
			"```lua\n%T: %s\n```",
			node,
			def.Type.String(),
		),
		Range: ptr(file.ToProtocolRange(ast.Range(node))),
	}, nil
}
