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
	nodePath := ast.GetNode(&file.File.Block, file.File.ToPos(params.Position))
	if nodePath.Node == nil {
		return nil, nil
	}
	highlights := []protocol.DocumentHighlight{}
	// TODO: Labels
	_, ok := nodePath.Node.(ast.LeafNode)
	if !ok {
		return nil, nil
	}
	// TODO: References
	highlights = append(highlights, protocol.DocumentHighlight{Range: file.File.ToProtocolRange(ast.Range(nodePath.Node))})

	def := file.Env.FindDefinition(nodePath)
	if def != nil {
		highlights = append(highlights, protocol.DocumentHighlight{Range: file.File.ToProtocolRange(ast.Range(def))})
	}

	return highlights, nil
}
