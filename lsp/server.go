package lsp

import (
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/types"
	"github.com/raiguard/luapls/util"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserv "github.com/tliron/glsp/server"

	// To enable logging
	_ "github.com/tliron/commonlog/simple"
)

const LS_NAME = "luapls"

type LegacyFile struct {
	File *ast.File
	Env  types.Environment
	Path string
}

// Server contains the state for the LSP session.
type Server struct {
	legacyFiles map[string]*LegacyFile
	fileGraph   types.FileGraph
	handler     protocol.Handler
	log         commonlog.Logger
	rootPath    string
	server      *glspserv.Server

	config Config

	isInitialized bool
}

func Run(logLevel int) {
	commonlog.Configure(logLevel, util.Ptr("/tmp/luapls.log"))

	s := Server{
		legacyFiles: map[string]*LegacyFile{},
		fileGraph: types.FileGraph{
			Roots: []*types.FileNode{},
			Files: map[string]*types.FileNode{},
		}}

	s.handler.Initialize = s.initialize
	s.handler.Initialized = s.initialized
	s.handler.WorkspaceDidChangeConfiguration = s.didChangeConfiguration
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

	s.updateConfig(params.InitializationOptions)

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &protocol.InitializeResultServerInfo{Name: LS_NAME},
	}, nil
}

func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	go func() {
		before := time.Now()
		for _, path := range *s.config.Roots {
			uri, err := pathToURI(path)
			if err != nil {
				s.log.Errorf("%s", err)
				continue
			}
			// TODO: Check for duplicate roots
			file := s.parseFile(uri, nil)
			if file != nil {
				s.fileGraph.Roots = append(s.fileGraph.Roots, file)
			}
		}
		s.isInitialized = true
		s.log.Debugf("Initialized in %s", time.Since(before).String())

		for _, file := range s.legacyFiles {
			s.publishDiagnostics(ctx, file)
		}

		s.fileGraph.Traverse(func(file *types.FileNode) bool {
			s.log.Debugf("FILE: %s", file.Path)
			s.log.Debug("  PARENTS:")
			for _, parent := range file.Parents {
				s.log.Debugf("    %s", parent.Path)
			}
			s.log.Debug("  CHILDREN:")
			for _, child := range file.Children {
				s.log.Debugf("    %s", child.Path)
			}
			return true
		})

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

func (s *Server) getFile(uri protocol.URI) *LegacyFile {
	if !s.isInitialized {
		return nil
	}
	existing := s.legacyFiles[uri]
	if existing != nil {
		return existing
	}
	return nil
}
