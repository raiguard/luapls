package types

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/tliron/commonlog"
)

type Environment struct {
	Files    FileGraph
	RootPath string

	log commonlog.Logger
}

func NewEnvironment() *Environment {
	e := &Environment{
		Files: FileGraph{Files: map[string]*File{}, Roots: []*File{}},
		log:   commonlog.GetLogger("luapls.environment"),
	}

	return e
}

// Init parses all Lua files in the root directory and builds the type graph.
func (e *Environment) Init() {
	filepath.Walk(e.RootPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".lua") {
			e.AddFile(path)
		}
		return nil
	})
}

func (e *Environment) AddFile(path string) *File {
	if existing := e.Files.Files[path]; existing != nil {
		return existing
	}
	src, err := os.ReadFile(path)
	if err != nil {
		e.log.Errorf("Failed to parse file %s: %s", path, err)
		return nil
	}
	timer := time.Now()
	astFile := parser.New(string(src)).ParseFile()
	e.log.Debugf("Parsed file '%s' in %s", path, time.Since(timer).String())

	file := &File{
		AST:         &astFile.Block,
		LineBreaks:  astFile.LineBreaks,
		Diagnostics: astFile.Errors,
		Types:       []*Type{},
		Parents:     []*File{},
		Children:    []*File{},
		Path:        path,
		Visited:     false,
	}
	e.Files.Files[path] = file

	ast.Walk(&astFile.Block, func(n ast.Node) bool {
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
		childPath := strings.ReplaceAll(pathNode.Token.Literal[1:len(pathNode.Token.Literal)-1], ".", "/")
		if !strings.HasSuffix(childPath, ".lua") {
			childPath += ".lua"
		}
		e.log.Debugf("Found require path %s", childPath)
		child := e.AddFile(childPath)
		if child != nil {
			file.Children = append(file.Children, child)
			if !slices.Contains(child.Parents, file) {
				child.Parents = append(child.Parents, file)
			}
		}
		return false // No children to iterate
	})
	return file
}
