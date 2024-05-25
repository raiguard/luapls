package lsp

import (
	"os"
	"time"

	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/types"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserv "github.com/tliron/glsp/server"

	// To enable logging
	_ "github.com/tliron/commonlog/simple"
)

const LS_NAME = "luapls"

type File struct {
	File *parser.File
	Env  types.Environment
	Path string
}

// Type Server contains the state for the LSP session.
type Server struct {
	files    map[string]*File
	handler  protocol.Handler
	log      commonlog.Logger
	rootPath string
	server   *glspserv.Server

	isInitialized bool
}

func Run(logLevel int) {
	commonlog.Configure(logLevel, ptr("/tmp/luapls.log"))

	s := Server{
		files: map[string]*File{},
	}

	s.handler.Initialize = s.initialize
	s.handler.Initialized = s.initialized
	s.handler.Shutdown = s.shutdown
	s.handler.SetTrace = s.setTrace
	s.handler.TextDocumentDidOpen = s.textDocumentDidOpen
	s.handler.TextDocumentDidChange = s.textDocumentDidChange
	s.handler.TextDocumentDidClose = s.textDocumentDidClose
	s.handler.TextDocumentDocumentHighlight = s.textDocumentHighlight
	s.handler.TextDocumentHover = s.textDocumentHover
	s.handler.TextDocumentDefinition = s.textDocumentDefinition

	s.server = glspserv.NewServer(&s.handler, LS_NAME, logLevel > 2)

	s.log = s.server.Log

	s.server.RunStdio()
}

func (s *Server) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()
	s.rootPath = *params.RootPath

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &protocol.InitializeResultServerInfo{Name: LS_NAME},
	}, nil
}

func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	go func() {
		s.isInitialized = true
		for _, file := range s.files {
			s.publishDiagnostics(ctx, file)
		}
	}()
	return nil
}

func (s *Server) shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *Server) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *Server) getFile(uri protocol.URI) *File {
	if !s.isInitialized {
		return nil
	}
	existing := s.files[uri]
	if existing != nil {
		return existing
	}

	// Otherwise, create, parse, and check it
	path, err := uriToPath(uri)
	if err != nil {
		s.log.Errorf("%s", err)
		return nil
	}
	src, err := os.ReadFile(path)
	if err != nil {
		s.log.Errorf("Failed to parse file %s: %s", uri, err)
		return nil
	}
	timer := time.Now()
	parserFile := parser.New(string(src)).ParseFile()
	file := &File{File: &parserFile, Env: types.NewEnvironment(&parserFile), Path: uri}
	file.Env.ResolveTypes()
	s.files[uri] = file
	s.log.Debugf("Parsed file '%s' in %s", uri, time.Since(timer).String())

	return file
}
