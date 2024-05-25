package lsp

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	node, _ := ast.GetNode(&file.File.Block, file.File.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	highlights := []protocol.DocumentHighlight{}
	// TODO: Labels
	if _, ok := node.(ast.LeafNode); ok {
		highlights = append(highlights, protocol.DocumentHighlight{Range: file.File.ToProtocolRange(ast.Range(node))})
	}
	if ident, ok := node.(*ast.Identifier); ok {
		def := file.Env.FindDefinition(ident, true)
		if def != nil {
			highlights = append(highlights, protocol.DocumentHighlight{Range: file.File.ToProtocolRange(ast.Range(def))})
		}
	}
	return highlights, nil
}
