package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func publishDiagnostics(ctx *glsp.Context, uri protocol.URI) {
	file := files[uri]
	if file == nil {
		return
	}
	diagnostics := []protocol.Diagnostic{}
	for _, err := range file.Errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    file.ToProtocolRange(err.Range),
			Severity: ptr(protocol.DiagnosticSeverityError),
			Message:  err.Message,
		})
	}
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}
