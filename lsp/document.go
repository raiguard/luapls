package lsp

import (
	"os"
	"slices"
	"strings"
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
	s.legacyFiles[params.TextDocument.URI] = nil
	return nil
}

func (s *Server) createFile(uri protocol.URI) *LegacyFile {
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
	file := &LegacyFile{File: &parserFile, Env: types.NewEnvironment(&parserFile), Path: uri}
	file.Env.ResolveTypes()
	s.legacyFiles[uri] = file
	s.log.Debugf("Parsed and checked file '%s' in %s", uri, time.Since(timer).String())

	return file
}

func (s *Server) parseFile(uri protocol.URI, parent *types.FileNode) *types.FileNode {
	if existing := s.fileGraph.Files[uri]; existing != nil {
		if !slices.Contains(existing.Parents, parent) {
			existing.Parents = append(existing.Parents, parent)
		}
		return existing
	}
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

	fileNode := &types.FileNode{
		AST:         &file.Block,
		LineBreaks:  file.LineBreaks,
		Diagnostics: file.Errors,
		Types:       []*types.Type{},
		Parents:     []*types.FileNode{},
		Children:    []*types.FileNode{},
		Path:        uri,
		Visited:     false,
	}
	if parent != nil {
		fileNode.Parents = append(fileNode.Parents, parent)
	}
	s.fileGraph.Files[uri] = fileNode
	s.log.Debugf("Walking %s", uri)

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
		s.log.Debugf("Found require path %s", pathURI)
		child := s.parseFile(pathURI, fileNode)
		if child != nil {
			fileNode.Children = append(fileNode.Children, child)
		}
		return false // No children to iterate
	})
	return fileNode
}
