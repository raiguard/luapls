package lsp

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}

	// Iterate the node's children
	// Find any local declarations in the children
	// Add the identifiers to a locals list
	// Recurse above on any node that doesn't create a new scope
	// Generate completions for those identifiers

	pos := file.ToPos(params.Position)

	locals := map[string]*ast.Identifier{}

	ast.Walk(&file.Block, func(node ast.Node) bool {
		// If the entire node is before the current position
		if pos >= node.Pos() && pos > node.End() {
			switch node := node.(type) {
			case *ast.LocalStatement:
				for _, ident := range node.Names {
					locals[ident.Literal] = ident
				}
			case *ast.FunctionStatement:
				if ident, ok := node.Left.(*ast.Identifier); ok {
					locals[ident.Literal] = ident
				}
			}
			return false
		}
		// If the node does not intersect the position
		if pos < node.Pos() || pos >= node.End() {
			return false
		}
		switch node := node.(type) {
		case *ast.ForInStatement:
			for _, ident := range node.Names {
				locals[ident.Literal] = ident
			}
		case *ast.ForStatement:
			if node.Name != nil {
				locals[node.Name.Literal] = node.Name
			}
		case *ast.FunctionExpression:
			for _, ident := range node.Params {
				locals[ident.Literal] = ident
			}
		case *ast.FunctionStatement:
			for _, ident := range node.Params {
				locals[ident.Literal] = ident
			}
		default:
			// If the node contains the position
			return node.Pos() <= pos && pos < node.End()
		}

		return true
	})

	completions := []protocol.CompletionItem{}

	for _, ident := range locals {
		completions = append(completions, protocol.CompletionItem{
			Label: ident.Literal,
			Kind:  ptr(protocol.CompletionItemKindVariable),
		})
	}

	return completions, nil
}
