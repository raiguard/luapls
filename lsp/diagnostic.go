package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) publishDiagnostics(ctx *glsp.Context, uri protocol.URI) {
	file := s.getFile(uri)
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
	s.validateLocals(file, &diagnostics)
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func (s *Server) validateLocals(file *parser.File, diagnostics *[]protocol.Diagnostic) {
	ast.Walk(&file.Block, func(node ast.Node) bool {
		ident, ok := node.(*ast.Identifier)
		if !ok {
			return true
		}

		locals := getLocals(&file.Block, ident.Pos(), true)
		if locals[ident.Literal] == nil {
			(*diagnostics) = append((*diagnostics), protocol.Diagnostic{
				Range: file.ToProtocolRange(ast.Range(ident)),
				// TODO: Configurable severity
				Severity: ptr(protocol.DiagnosticSeverityInformation),
				Message:  fmt.Sprintf("Unknown variable '%s'", ident.Literal),
			})
		}

		return true
	})
}
