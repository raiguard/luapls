package lsp

import (
	"time"

	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TODO: Incremental changes
func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil
	}
	s.publishDiagnostics(ctx, file)
	return nil
}

func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil
	}
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			before := time.Now()
			newFile := parser.New(change.Text).ParseFile()
			file.AST = &newFile.Block
			file.LineBreaks = newFile.LineBreaks
			file.Diagnostics = newFile.Errors
			// TODO: Re-add type errors
			s.log.Debugf("Reparse duration: %s", time.Since(before).String())
			s.publishDiagnostics(ctx, file)
		}
	}
	return nil
}

func (s *Server) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.legacyFiles[params.TextDocument.URI] = nil
	return nil
}
