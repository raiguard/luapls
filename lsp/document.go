package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	parseFile(ctx, params.TextDocument.URI, params.TextDocument.Text)
	publishDiagnostics(ctx, params.TextDocument.URI)
	return nil
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			parseFile(ctx, params.TextDocument.URI, change.Text)
			publishDiagnostics(ctx, params.TextDocument.URI)
		}
	}
	return nil
}

func textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}
	node := getInnermostNode(&file.Block, file.ToPos(params.Position))
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
	node := getInnermostNode(&file.Block, file.ToPos(params.Position))
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

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	return reserved, nil
}
