package lsp

import (
	"encoding/json"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func toJSON(v any) string {
	res, err := json.Marshal(v)
	if err != nil {
		return "<ERROR>"
	}
	return string(res)
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
				for _, ident := range node.Names.Pairs {
					locals[ident.Node.Token.Literal] = ident.Node
				}
			}
		case *ast.ForStatement:
			if isInside {
				if node.Name != nil {
					locals[node.Name.Token.Literal] = node.Name
				}
			}
		case *ast.FunctionExpression:
			if isInside {
				for _, ident := range node.Params.Pairs {
					locals[ident.Node.Token.Literal] = ident.Node
				}
			}
		case *ast.FunctionStatement:
			if isInside {
				for _, ident := range node.Params.Pairs {
					locals[ident.Node.Token.Literal] = ident.Node
				}
			}
			if isBefore || includeSelf {
				if ident, ok := node.Name.(*ast.Identifier); ok {
					locals[ident.Token.Literal] = ident
				}
			}
		case *ast.LocalStatement:
			if isBefore || includeSelf {
				for _, ident := range node.Names.Pairs {
					locals[ident.Node.Token.Literal] = ident.Node
				}
			}
		default:
			return isInside
		}

		return true
	})

	return locals
}
