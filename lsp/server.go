package lsp

import (
	"strings"

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

type File struct {
	File *ast.File
	Env  types.Environment
	Path string
}

// Server contains the state for the LSP session.
type Server struct {
	files    map[string]*File
	handler  protocol.Handler
	log      commonlog.Logger
	rootPath string
	server   *glspserv.Server

	config Config

	isInitialized bool
}

func Run(logLevel int) {
	commonlog.Configure(logLevel, util.Ptr("/tmp/luapls.log"))

	s := Server{files: map[string]*File{}}

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
		s.isInitialized = true
		s.log.Debug("Initialized")
		toParse := []string{}
		for _, path := range *s.config.Roots {
			uri, err := pathToURI(path)
			if err != nil {
				s.log.Errorf("%s", err)
			}
			toParse = append(toParse, uri)
		}
		parsed := map[string]bool{}
		for i := 0; i < len(toParse); i++ {
			uri := toParse[i]
			// TODO: Normalize paths
			if parsed[uri] {
				continue
			}
			file := s.parseFile(uri)
			if file == nil {
				continue
			}
			parsed[uri] = true
			s.log.Debugf("walking %s", uri)

			ast.Walk(&file.Block, func(n ast.Node) bool {
				fc, ok := n.(*ast.FunctionCall)
				if !ok {
					return true
				}
				ident, ok := fc.Name.(*ast.Identifier)
				if !ok || ident.Token.Literal != "require" {
					return true
				}
				if len(fc.Args.Pairs) != 1 {
					return true
				}
				pathNode, ok := fc.Args.Pairs[0].Node.(*ast.StringLiteral)
				if !ok {
					return false // There are no children to iterate at this point
				}
				// TODO: Clean this up
				pathString := strings.ReplaceAll(pathNode.Token.Literal[1:len(pathNode.Token.Literal)-1], ".", "/")
				if !strings.HasSuffix(pathString, ".lua") {
					pathString += ".lua"
				}
				pathURI, err := pathToURI(pathString)
				if err != nil {
					s.log.Errorf("%s", err)
					return false
				}
				s.log.Debugf("found require path %s", pathURI)
				toParse = append(toParse, pathURI)
				return false // No children to iterate
			})
		}

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
	return nil
}
