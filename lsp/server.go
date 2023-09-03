package lsp

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "luapls"

var version string = "0.0.1"
var handler protocol.Handler

var files = map[string][]token.Token{}
var rootPath string = ""

func Run() {
	handler.Initialize = initialize
	handler.Initialized = initialized
	handler.Shutdown = shutdown
	handler.SetTrace = setTrace
	handler.TextDocumentDidOpen = textDocumentDidOpen
	handler.TextDocumentDidChange = textDocumentDidChange
	handler.TextDocumentDocumentHighlight = textDocumentHighlight
	handler.TextDocumentHover = textDocumentHover

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
		lexFile(path, string(src))
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
	lexFile(params.TextDocument.URI, params.TextDocument.Text)
	return nil
}

func textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			before := time.Now()
			lexFile(params.TextDocument.URI, change.Text)
			logToEditor(ctx, fmt.Sprint("Rescan duration: ", time.Since(before)))
		}
	}
	return nil
}

func textDocumentHighlight(ctc *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	if tokens, ok := files[params.TextDocument.URI]; ok {
		for i := 0; i < len(tokens); i++ {
			tok := &tokens[i]
			if withinRange(&tok.Range, &params.Position) {
				return []protocol.DocumentHighlight{{
					Range: toProtocolRange(&tok.Range),
				}}, nil
			}
		}
	}
	return nil, nil
}

func textDocumentHover(ctc *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if tokens, ok := files[params.TextDocument.URI]; ok {
		for i := 0; i < len(tokens); i++ {
			tok := &tokens[i]
			if withinRange(&tok.Range, &params.Position) {
				return &protocol.Hover{
					Contents: fmt.Sprintf("Type: %s\n\nLiteral:\n```\n%s\n```", token.TokenStr[tok.Type], tok.Literal),
					// Range:    &token.Range,
				}, nil
			}
		}
	}
	return nil, nil
}

func lexFile(filename, src string) {
	l := lexer.New(src)
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF || tok.Type == token.INVALID {
			break
		}
		files[filename] = append(files[filename], tok)
	}
}

func logToEditor(ctx *glsp.Context, msg string) {
	ctx.Notify(
		protocol.ServerWindowLogMessage,
		protocol.LogMessageParams{Type: protocol.MessageTypeLog, Message: msg},
	)
}

func withinRange(rng *token.Range, pos *protocol.Position) bool {
	startCol, endCol := rng.StartCol, rng.EndCol
	startRow, endRow := rng.StartRow, rng.EndRow
	if startRow < pos.Line {
		startCol = 0
		startRow = pos.Line
	}
	if endRow > pos.Line {
		endCol = math.MaxUint32
		endRow = pos.Line
	}
	return startRow == pos.Line && endRow == pos.Line && startCol <= pos.Character && endCol > pos.Character
}

func toProtocolRange(rng *token.Range) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{Line: rng.StartRow, Character: rng.StartCol},
		End:   protocol.Position{Line: rng.EndRow, Character: rng.EndCol},
	}
}