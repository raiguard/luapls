package lsp

import (
	"fmt"

	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func parseFile(ctx *glsp.Context, filename, src string) {
	file := parser.New(src).ParseFile()
	files[filename] = &file
}

func logToEditor(ctx *glsp.Context, format string, args ...any) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: fmt.Sprintf(format, args...)},
	)
}
