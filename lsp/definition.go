package lsp

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDefinition(ctx *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}

	pos := file.ToPos(params.Position)

	node, _ := ast.GetNode(&file.Block, pos)
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil, nil
	}

	locals := getLocals(&file.Block, pos, true)

	def := locals[ident.Literal]
	if def == nil {
		return nil, nil
	}

	return &protocol.Location{
		URI:   params.TextDocument.URI,
		Range: file.ToProtocolRange(ast.Range(def)),
	}, nil
}
