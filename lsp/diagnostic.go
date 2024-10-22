package lsp

import (
	"github.com/raiguard/luapls/lua/ast"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) publishDiagnostics(ctx *glsp.Context, file *ast.File) {
	diagnostics := []protocol.Diagnostic{}
	for _, err := range file.Diagnostics {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    file.LineBreaks.ToProtocolRange(err.Range),
			Severity: &err.Severity,
			Message:  err.Message,
		})
	}
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         file.URI,
		Diagnostics: diagnostics,
	})
}
