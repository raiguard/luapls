package lsp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TODO: Listen for modifications to files that were previously parsed

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	parseFile(ctx, params.TextDocument.URI, params.TextDocument.Text)
	publishDiagnostics(ctx, params.TextDocument.URI)
	return nil
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			before := time.Now()
			parseFile(ctx, params.TextDocument.URI, change.Text)
			publishDiagnostics(ctx, params.TextDocument.URI)
			logToEditor(ctx, fmt.Sprint("Reparse duration:", time.Since(before).String()))
		}
	}
	return nil
}

func textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}
	node, _ := ast.GetNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	return []protocol.DocumentHighlight{
		{Range: file.ToProtocolRange(ast.Range(node))},
	}, nil
}

func textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}
	node, _ := ast.GetNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	return &protocol.Hover{
		Contents: fmt.Sprintf(
			"# %T\n\nRange: `{%d, %d, %d}` `{%d, %d, %d}`",
			node,
			file.ToProtocolRange(ast.Range(node)).Start.Line,
			file.ToProtocolRange(ast.Range(node)).Start.Character,
			node.Pos(),
			file.ToProtocolRange(ast.Range(node)).End.Line,
			file.ToProtocolRange(ast.Range(node)).End.Character,
			node.End(),
		),
		Range: ptr(file.ToProtocolRange(ast.Range(node))),
	}, nil
}

func textDocumentSelectionRange(ctx *glsp.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	file := files[params.TextDocument.URI]
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
