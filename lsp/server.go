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

type FileGraph struct {
	Roots []*FileNode
	Files map[protocol.URI]*FileNode
}

func (fg *FileGraph) Traverse(visitor func(fn *FileNode) bool) {
	for _, root := range fg.Roots {
		root.Traverse(visitor)
	}
	for _, file := range fg.Files {
		file.visited = false
	}
}

// TODO: Atomics to allow multithreading
type FileNode struct {
	// AST is discarded after type checking is complete, unless the file is open in the editor.
	AST        *ast.Block
	LineBreaks []int

	Diagnostics []ast.Error
	Types       []*types.Type

	Parents  []*FileNode
	Children []*FileNode

	Path    string
	visited bool
}

func (fn *FileNode) Traverse(visitor func(fn *FileNode) bool) {
	if fn == nil || fn.visited {
		return
	}
	fn.visited = true
	if visitor(fn) {
		for _, child := range fn.Children {
			child.Traverse(visitor)
		}
	}
}

// Server contains the state for the LSP session.
type Server struct {
	legacyFiles map[string]*LegacyFile
	fileGraph   FileGraph
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
		fileGraph: FileGraph{
			Roots: []*FileNode{},
			Files: map[string]*FileNode{},
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

		s.fileGraph.Traverse(func(file *FileNode) bool {
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
