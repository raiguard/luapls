package main

import (
	"luapls/lua"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var version string = "0.0.1"
var handler protocol.Handler

var files = map[protocol.DocumentUri][]lua.Token{}

func main() {
	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen
	handler.TextDocumentHover = textDocumentHover
	handler.TextDocumentDocumentHighlight = textDocumentHighlight

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
	tokens := lua.Tokenize(strings.Split(params.TextDocument.Text, "\n"))
	files[params.TextDocument.URI] = tokens
	return nil
}

func textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	token, err := findToken(params.TextDocument.URI, params.Position)
	if err != nil {
		return nil, err
	}
	if token != nil {
		// Show token information
		return &protocol.Hover{
			Contents: lua.StrToken(token),
			Range:    &token.Range,
		}, nil
	}

	return nil, nil
}

func textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	token, err := findToken(params.TextDocument.URI, params.Position)
	if err != nil {
		return nil, err
	}
	if token != nil {
		return []protocol.DocumentHighlight{{
			Range: token.Range,
		}}, nil
	}

	return nil, nil
}
