package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	node, _ := ast.GetNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	return &protocol.Hover{
		Contents: fmt.Sprintf(
			"# %T\n\nRange: `{%d, %d, %d}` `{%d, %d, %d}`",
			node,
			file.ToProtocolRange(ast.Range(node)).Start.Line,
			file.ToProtocolRange(ast.Range(node)).Start.Character,
			node.Pos(),
			file.ToProtocolRange(ast.Range(node)).End.Line,
			file.ToProtocolRange(ast.Range(node)).End.Character,
			node.End(),
		),
		Range: ptr(file.ToProtocolRange(ast.Range(node))),
	}, nil
}
