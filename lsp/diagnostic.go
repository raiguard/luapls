package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) publishDiagnostics(ctx *glsp.Context, file *File) {
	diagnostics := []protocol.Diagnostic{}
	for _, err := range file.File.Errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    file.File.ToProtocolRange(err.Range),
			Severity: ptr(protocol.DiagnosticSeverityError),
			Message:  err.Message,
		})
	}
	for _, err := range file.Env.Errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    file.File.ToProtocolRange(err.Range),
			Severity: ptr(protocol.DiagnosticSeverityWarning),
			Message:  err.Message,
		})
	}
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         file.Path,
		Diagnostics: diagnostics,
	})
}
