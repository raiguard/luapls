package lsp

import (
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TODO: Incremental changes
func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	// TODO: Check if file was not yet parsed
	s.publishDiagnostics(ctx, params.TextDocument.URI)
	return nil
}

func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			before := time.Now()
			newFile := parser.New(change.Text).ParseFile()
			s.files[params.TextDocument.URI] = &newFile
			s.log.Debugf("Reparse duration: %s", time.Since(before).String())
			s.publishDiagnostics(ctx, params.TextDocument.URI)
		}
	}
	return nil
}

func (s *Server) textDocumentSelectionRange(ctx *glsp.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil, nil
	}
	ranges := []protocol.SelectionRange{}
	for _, position := range params.Positions {
		node, parents := ast.GetNode(&file.Block, file.ToPos(position))
		ranges = append(ranges, protocol.SelectionRange{
			Range: file.ToProtocolRange(ast.Range(node)),
		})
		curRange := &ranges[len(ranges)-1]
		for i := len(parents) - 1; i >= 0; i-- {
			parent := &parents[i]
			if parent, ok := (*parent).(*ast.TableField); ok {
				// The table field will have the same selection range as the value node.
				if parent.Key == nil {
					continue
				}
			}
			parentRange := protocol.SelectionRange{Range: file.ToProtocolRange(ast.Range(*parent))}
			curRange.Parent = &parentRange
			curRange = &parentRange
		}
	}
	return ranges, nil
}
