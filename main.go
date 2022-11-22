package main

import (
	"fmt"
	"luapls/lua"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var version string = "0.0.1"
var handler protocol.Handler

func main() {
	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen

	server := server.NewServer(&handler, "gopls", true)

	server.RunStdio()
}

func initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	var s lua.Scanner
	s.Init([]byte(params.TextDocument.Text))

	for {
		pos, tok, lit := s.Scan()
		if tok == lua.EOF {
			break
		}
		logToEditor(ctx, fmt.Sprint(pos, tok, " ", lit))
	}

	return nil
}
