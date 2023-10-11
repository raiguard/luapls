package lsp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var handler protocol.Handler

var files = map[string]*parser.File{}
var rootPath string = ""

// ptr returns a pointer to the given value.
func ptr[T any](value T) *T {
	inner := value
	return &inner
}

var reserved []protocol.CompletionItem

func Run() {
	reserved = []protocol.CompletionItem{}
	for literal := range token.Reserved {
		reserved = append(reserved, protocol.CompletionItem{
			Label: literal,
			Kind:  ptr(protocol.CompletionItemKindKeyword),
		})
	}

	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen
	handler.TextDocumentDidChange = textDocumentDidChange
	handler.TextDocumentDocumentHighlight = textDocumentHighlight
	handler.TextDocumentHover = textDocumentHover
	handler.TextDocumentCompletion = textDocumentCompletion

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	rootPath = *params.RootPath

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &protocol.InitializeResultServerInfo{Name: lsName},
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
			logToEditor(ctx, "%s", err)
			continue
		}
		parseFile(ctx, path, string(src))
	}
	logToEditor(ctx, "Initial parse (%d files): %s", len(toParse), time.Since(before))
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
