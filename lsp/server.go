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
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
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

var files = Mutex[map[string]*parser.File]{
	inner: map[string]*parser.File{},
	mu:    sync.Mutex{},
}
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
	publishDiagnostics(ctx, params.TextDocument.URI)
	return nil
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			parseFile(ctx, params.TextDocument.URI, change.Text)
			publishDiagnostics(ctx, params.TextDocument.URI)
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
	node := getInnermostNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	return []protocol.DocumentHighlight{
		{Range: file.ToProtocolRange(ast.Range(node))},
	}, nil
}

func textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	file := files.Lock()[params.TextDocument.URI]
	defer files.Unlock()
	if file == nil {
		return nil, nil
	}
	node := getInnermostNode(&file.Block, file.ToPos(params.Position))
	if node == nil {
		return nil, nil
	}
	return &protocol.Hover{
		Contents: fmt.Sprintf(
			"# %T\n\nRange: `{%d, %d, %d}` `{%d, %d, %d}`",
			node,
			file.ToProtocolRange(ast.Range(node)).Start.Line,
			file.ToProtocolRange(ast.Range(node)).Start.Character,
			node.Pos(),
			file.ToProtocolRange(ast.Range(node)).End.Line,
			file.ToProtocolRange(ast.Range(node)).End.Character,
			node.End(),
		),
		Range: ptr(file.ToProtocolRange(ast.Range(node))),
	}, nil
}

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	return reserved, nil
}

func publishDiagnostics(ctx *glsp.Context, uri protocol.URI) {
	file := files.Lock()[uri]
	defer files.Unlock()
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
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func getInnermostNode(n ast.Node, pos token.Pos) ast.Node {
	var node ast.Node
	ast.Walk(n, func(n ast.Node) bool {
		if n.Pos() <= pos && pos < n.End() {
			node = n
			return true
		}
		return false
	})
	return node
}

func parseFile(ctx *glsp.Context, filename, src string) {
	file := parser.New(src).ParseFile()
	files.Lock()[filename] = &file
	files.Unlock()
}

func logToEditor(ctx *glsp.Context, format string, args ...any) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: fmt.Sprintf(format, args...)},
	)
}
