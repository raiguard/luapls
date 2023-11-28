package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	file := files[params.TextDocument.URI]
	if file == nil {
		return nil, nil
	}

	pos := file.ToPos(params.Position)

	completions := []protocol.CompletionItem{}

	for _, ident := range getLocals(&file.Block, pos) {
		completions = append(completions, protocol.CompletionItem{
			Label: ident.Literal,
			Kind:  ptr(protocol.CompletionItemKindVariable),
		})
	}

	return completions, nil
}
