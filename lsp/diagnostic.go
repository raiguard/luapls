package lsp

import (
	"github.com/raiguard/luapls/lua/types"
	"github.com/raiguard/luapls/util"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) publishDiagnostics(ctx *glsp.Context, file *types.File) {
	diagnostics := []protocol.Diagnostic{}
	for _, err := range file.Diagnostics {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range:    file.LineBreaks.ToProtocolRange(err.Range),
			Severity: util.Ptr(protocol.DiagnosticSeverityError),
			Message:  err.Message,
		})
	}
	// for _, err := range file.Env.Errors {
	// 	diagnostics = append(diagnostics, protocol.Diagnostic{
	// 		Range:    file.LineBreaks.ToProtocolRange(err.Range),
	// 		Severity: util.Ptr(protocol.DiagnosticSeverityWarning),
	// 		Message:  err.Message,
	// 	})
	// }
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         file.URI,
		Diagnostics: diagnostics,
	})
}
