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

	pos := file.ToPos(params.Position)

	node, parents := ast.GetNode(&file.Block, pos)
	def := getDefinition(node, parents)
	if def == nil {
		return nil, nil
	}

	return &protocol.Location{
		URI:   params.TextDocument.URI,
		Range: file.ToProtocolRange(ast.Range(def.Definition)),
	}, nil
}

func getDefinition(node ast.Node, parents []ast.Node) *ast.VariableDeclaration {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	for i := len(parents) - 1; i >= 0; i-- {
		parent := parents[i]
		block, ok := parent.(*ast.Block)
		if !ok {
			continue
		}
		def := block.Locals[ident.Literal]
		if def == nil {
			continue
		}
		// TODO: Handle shadowing
		if def.Definition.Pos() > node.Pos() {
			continue
		}
		return def
	}
	return nil
}
