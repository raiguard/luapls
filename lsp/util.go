package lsp

import (
	"encoding/json"
	"net/url"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func toJSON(v any) string {
	res, err := json.Marshal(v)
	if err != nil {
		return "<ERROR>"
	}
	return string(res)
}

// ptr returns a pointer to the given value.
func ptr[T any](value T) *T {
	return &value
}

func (s *Server) uriToPath(uri protocol.URI) string {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		s.log.Errorf("Failed to parse file URI: %s", err)
	}
	return u.Path
}

// getLocals returns a list of all local variables contained in node for the given pos.
func getLocals(node ast.Node, pos token.Pos, includeSelf bool) map[string]*ast.Identifier {
	locals := map[string]*ast.Identifier{}

	ast.Walk(node, func(node ast.Node) bool {
		isAfter := node.Pos() > pos && pos < node.End()
		if isAfter {
			return false
		}
		isBefore := node.Pos() <= pos && pos > node.End()
		isInside := node.Pos() <= pos && pos < node.End()
		switch node := node.(type) {
		case *ast.ForInStatement:
			if isInside {
				for _, ident := range node.Names {
					locals[ident.Literal] = ident
				}
			}
		case *ast.ForStatement:
			if isInside {
				if node.Name != nil {
					locals[node.Name.Literal] = node.Name
				}
			}
		case *ast.FunctionExpression:
			if isInside {
				for _, ident := range node.Params {
					locals[ident.Literal] = ident
				}
			}
		case *ast.FunctionStatement:
			if isInside {
				for _, ident := range node.Params {
					locals[ident.Literal] = ident
				}
			}
			if isBefore || includeSelf {
				if ident, ok := node.Left.(*ast.Identifier); ok {
					locals[ident.Literal] = ident
				}
			}
		case *ast.LocalStatement:
			if isBefore || includeSelf {
				for _, ident := range node.Names {
					locals[ident.Literal] = ident
				}
			}
		default:
			return isInside
		}

		return true
	})

	return locals
}
