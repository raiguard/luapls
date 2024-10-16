package lsp

import (
	"os"
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/types"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TODO: Incremental changes
func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	file := s.createFile(params.TextDocument.URI)
	if file == nil {
		return nil
	}
	s.publishDiagnostics(ctx, file)
	return nil
}

func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	file := s.getFile(params.TextDocument.URI)
	if file == nil {
		return nil
	}
	for _, change := range params.ContentChanges {
		if change, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			before := time.Now()
			newFile := parser.New(change.Text).ParseFile()
			file.File = &newFile
			file.Env = types.NewEnvironment(&newFile)
			file.Env.ResolveTypes()
			s.log.Debugf("Reparse duration: %s", time.Since(before).String())
			s.publishDiagnostics(ctx, file)
		}
	}
	return nil
}

func (s *Server) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.files[params.TextDocument.URI] = nil
	return nil
}

func (s *Server) createFile(uri protocol.URI) *File {
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
	s.log.Debugf("Parsed and checked file '%s' in %s", uri, time.Since(timer).String())

	return file
}

func (s *Server) parseFile(uri protocol.URI) *ast.File {
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
	file := parser.New(string(src)).ParseFile()
	s.log.Debugf("Parsed file '%s' in %s", uri, time.Since(timer).String())

	return &file
}
