package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/types"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	nodePath := ast.GetNode(&file.File.Block, file.File.ToPos(params.Position))
	if nodePath.Node == nil {
		return nil, nil
	}
	ident, ok := nodePath.Node.(*ast.Identifier)
	if !ok {
		return nil, nil
	}
	typ, ok := file.Env.Types[ident]
	if !ok {
		typ = &types.Unknown{}
	}
	contents := fmt.Sprintf("```lua\n(variable) %s: %s\n```", ident.Literal, typ)
	// comments := ident.GetComments()
	// i := len(nodePath.Parents) - 1
	// for comments == "" && i >= 0 {
	// 	comments = nodePath.Parents[i].GetComments()
	// 	if _, ok := nodePath.Parents[i].(ast.Statement); ok {
	// 		break
	// 	}
	// 	i--
	// }
	// if comments != "" {
	// 	contents += "\n\n" + comments
	// }
	return &protocol.Hover{
		Contents: contents,
		Range:    ptr(file.File.ToProtocolRange(ast.Range(ident))),
	}, nil
}
