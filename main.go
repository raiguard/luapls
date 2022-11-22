package main

import (
	"fmt"
	"io/fs"
	"luapls/lua"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var version string = "0.0.1"
var handler protocol.Handler

var files = map[string][]lua.TokenExt{}
var rootPath string = ""

func main() {
	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen
	handler.TextDocumentDocumentHighlight = textDocumentHighlight

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	rootPath = *params.RootPath

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	var toParse []string
	filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".lua") {
			toParse = append(toParse, path)
		}
		return nil
	})
	before := time.Now()
	for _, path := range toParse {
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		parseFile(path, src)
	}
	logToEditor(ctx, fmt.Sprint("Initial scan (", len(toParse), " files): ", time.Since(before)))
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
	parseFile(params.TextDocument.URI, []byte(params.TextDocument.Text))
	return nil
}

func textDocumentHighlight(ctc *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	if tokens, ok := files[params.TextDocument.URI]; ok {
		for i := 0; i < len(tokens); i++ {
			token := &tokens[i]
			if withinRange(&token.Range, &params.Position) {
				return []protocol.DocumentHighlight{{
					Range: token.Range,
				}}, nil
			}
		}
	}
	return nil, nil
}

func parseFile(filename string, src []byte) {
	var p lua.Parser
	p.Init(filename, src)
	p.Parse()
	files[filename] = p.Tokens
}
