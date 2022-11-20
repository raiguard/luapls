package main

import (
	"errors"
	"luapls/lua"

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
	return rng.Start.Line == pos.Line && rng.Start.Character <= pos.Character && rng.End.Character > pos.Character
}

func findToken(uri protocol.DocumentUri, pos protocol.Position) (*lua.Token, error) {
	// Find	the	current	token
	tokens := files[uri]
	if tokens == nil {
		return nil, errors.New("Invalid	file URI")
	}
	for i := range tokens {
		tok := &tokens[i]
		if withinRange(&tok.Range, &pos) {
			return tok, nil
		}
	}
	return nil, nil
}
