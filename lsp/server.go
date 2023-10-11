package lsp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var version string = "0.0.1"
var handler protocol.Handler

// Go doesn't have a generics-based mutex yet, so let's roll our own!
type Mutex[T any] struct {
	inner T
	mu    sync.Mutex
}

func (m *Mutex[T]) Lock() T {
	m.mu.Lock()
	return m.inner
}

func (m *Mutex[T]) Unlock() {
	m.mu.Unlock()
}

var files = Mutex[map[string]*ast.File]{
	inner: map[string]*ast.File{},
	mu:    sync.Mutex{},
}
var rootPath string = ""

func Run() {
	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen
	handler.TextDocumentDidChange = textDocumentDidChange
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
	var parseWg sync.WaitGroup
	for _, path := range toParse {
		src, err := os.ReadFile(path)
		if err != nil {
			logToEditor(ctx, "%s", err)
			continue
		}
		parseWg.Add(1)
		go func(path string) {
			parseFile(ctx, path, string(src))
			parseWg.Done()
		}(path)
	}
	parseWg.Wait()
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

func textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	parseFile(ctx, params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			parseFile(ctx, params.TextDocument.URI, change.Text)
		}
	}
	return nil
}

func textDocumentHighlight(ctx *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	file := files.Lock()[params.TextDocument.URI]
	defer files.Unlock()
	if file == nil {
		return nil, nil
	}
	pos := file.ToPos(params.Position)
	for _, stmt := range file.Block.Stmts {
		if stmt.Pos() <= pos && pos < stmt.End() {
			kind := protocol.DocumentHighlightKindText
			return []protocol.DocumentHighlight{
				{Kind: &kind, Range: toProtocolRange(file, stmt)},
			}, nil
		}
	}
	return nil, nil
}

func toProtocolRange(file *ast.File, node ast.Node) protocol.Range {
	return protocol.Range{
		Start: file.ToProtocolPos(node.Pos()),
		End:   file.ToProtocolPos(node.End()),
	}
}

func parseFile(ctx *glsp.Context, filename, src string) {
	p := parser.New(lexer.New(src))
	file := p.ParseFile()
	if len(p.Errors()) > 0 {
		logToEditor(ctx, "Errors parsing %s:", filename)
		for _, err := range p.Errors() {
			logToEditor(ctx, "%s", err)
		}
	}
	files.Lock()[filename] = &file
	files.Unlock()
}

func logToEditor(ctx *glsp.Context, format string, args ...any) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: fmt.Sprintf(format, args...)},
	)
}
