package lsp

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentDefinition(ctx *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}

	pos := file.File.ToPos(params.Position)

	nodePath := ast.GetNode(&file.File.Block, pos)

	def := file.Env.FindDefinition(nodePath)
	if def == nil {
		return nil, nil
	}

	return &protocol.Location{
		URI:   params.TextDocument.URI,
		Range: file.File.ToProtocolRange(ast.Range(def)),
	}, nil
}
