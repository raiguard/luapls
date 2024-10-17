package lsp

import (
	"errors"
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/util"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	if file.AST == nil {
		return nil, errors.New("Attempted to highlight file with no AST")
	}
	nodePath := ast.GetNode(file.AST, file.LineBreaks.ToPos(params.Position))
	if nodePath.Node == nil {
		return nil, nil
	}
	ident, ok := nodePath.Node.(*ast.Identifier)
	if !ok {
		return nil, nil
	}
	// typ, ok := file.Env.Types[ident]
	// if !ok {
	// 	typ = &types.Unknown{}
	// }
	contents := fmt.Sprintf("```lua\n(variable) %s\n```", ident.Token.Literal)
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
		Range:    util.Ptr(file.LineBreaks.ToProtocolRange(ast.Range(ident))),
	}, nil
}
