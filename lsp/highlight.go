package lsp

import (
	"errors"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, errors.New("File not found")
	}
	if file.AST == nil {
		return nil, errors.New("Attempted to highlight file that has no AST")
	}
	nodePath := ast.GetSemanticNode(file.AST, file.LineBreaks.ToPos(params.Position))
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
	highlights = append(highlights, protocol.DocumentHighlight{Range: file.LineBreaks.ToProtocolRange(ast.Range(nodePath.Node))})

	// TODO:
	// def := file.Env.FindDefinition(nodePath)
	// if def != nil {
	// 	highlights = append(highlights, protocol.DocumentHighlight{Range: file.LineBreaks.ToProtocolRange(ast.Range(def))})
	// }

	return highlights, nil
}
