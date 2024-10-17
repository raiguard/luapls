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
	"github.com/raiguard/luapls/util"
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
			e.AddFile(path, nil)
		}
		return nil
	})

	for _, root := range e.Files.Roots {
		e.log.Debug(root.Path)
	}
}

func (e *Environment) AddFile(path string, parent *File) *File {
	if existing := e.Files.Files[path]; existing != nil {
		if parent != nil && !slices.Contains(existing.Parents, parent) {
			existing.Parents = append(existing.Parents, parent)
			if i := slices.Index(e.Files.Roots, existing); i >= 0 {
				e.Files.Roots = slices.Delete(e.Files.Roots, i, i)
			}
		}
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
	if parent != nil {
		file.Parents = append(file.Parents, parent)
	} else if !slices.Contains(e.Files.Roots, file) {
		e.Files.Roots = append(e.Files.Roots, file)
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
		// Even though Lua differentiates the returned module based on the exact contents of the string, for the purposes of
		// linting, we want to deduplicate.
		// Remove quotes
		stringContents := pathNode.Token.Literal[1 : len(pathNode.Token.Literal)-1]
		e.log.Debugf("Found require path %s", stringContents)
		// TODO: Handle ..
		childPath := strings.ReplaceAll(stringContents, ".", "/")
		if !strings.HasSuffix(childPath, ".lua") {
			childPath += ".lua"
		}
		// TODO: This is hardcoded for use with Factorio and must be generalized in the future.
		if strings.HasPrefix(childPath, "__") {
			childPath = strings.ReplaceAll(childPath, "__", "")
		}

		relativePath := filepath.Join(filepath.Dir(file.Path), childPath)
		// e.log.Debugf("Trying relative path %s", relativePath)
		var child *File
		if util.FileExists(relativePath) {
			child = e.AddFile(relativePath, file)
		} else if util.FileExists(childPath) { // Root
			child = e.AddFile(childPath, file)
		}
		if child != nil {
			file.Children = append(file.Children, child)
			if !slices.Contains(child.Parents, file) {
				child.Parents = append(child.Parents, file)
			}
		} else {
			e.log.Errorf("Unable to find file to match require path %s", stringContents)
		}
		return false // No children to iterate
	})
	return file
}
