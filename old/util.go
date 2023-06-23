package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func logToEditor(ctx *glsp.Context, msg string) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: msg},
	)
}

func withinRange(rng *protocol.Range, pos *protocol.Position) bool {
    // TODO: Multiline tokens (raw strings)
	return rng.Start.Line == pos.Line && rng.Start.Character <= pos.Character && rng.End.Character > pos.Character
}
