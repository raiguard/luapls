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
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Environment struct {
	Files    FileGraph
	RootPath string

	log commonlog.Logger
}

func NewEnvironment() *Environment {
	return &Environment{
		Files: FileGraph{Files: map[protocol.URI]*File{}, Roots: []*File{}},
		log:   commonlog.GetLogger("luapls.environment"),
	}
}

// Init parses all Lua files in the root directory and builds the type graph.
func (e *Environment) Init() {
	before := time.Now()
	filepath.WalkDir(e.RootPath, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".lua") {
			uri, err := util.PathToURI(path)
			if err != nil {
				return err
			}
			e.AddFile(uri, nil)
		}
		return nil
	})
	e.log.Debugf("Initialization took %s", time.Since(before).String())

	e.log.Debugf("ROOTS:")
	for _, root := range e.Files.Roots {
		e.log.Debug(root.URI)
	}

	e.log.Debugf("FILES:")
	for uri := range e.Files.Files {
		e.log.Debug(uri)
	}
}

func (e *Environment) AddFile(uri protocol.URI, parent *File) *File {
	if existing := e.Files.Files[uri]; existing != nil {
		e.log.Debugf("FOUND EXISTING %s", uri)
		if parent != nil && !slices.Contains(existing.Parents, parent) {
			existing.Parents = append(existing.Parents, parent)
			e.Files.Roots = slices.DeleteFunc(e.Files.Roots, func(file *File) bool { return file == existing })
		}
		return existing
	}
	e.log.Debugf("Parsing uri %s", uri)
	path, err := util.URIToPath(uri)
	if err != nil {
		e.log.Errorf("%s", err)
		return nil
	}
	// e.log.Debugf("Parsing %s", path)
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
		URI:         uri,
		Visited:     false,
	}
	if parent != nil {
		file.Parents = append(file.Parents, parent)
	} else if !slices.Contains(e.Files.Roots, file) {
		e.Files.Roots = append(e.Files.Roots, file)
	}
	e.Files.Files[uri] = file

	ast.Walk(&astFile.Block, func(n ast.Node) bool {
		fc, ok := n.(*ast.FunctionCall)
		if !ok {
			return true
		}
		ident, ok := fc.Name.(*ast.Identifier)
		// TODO: Don't hardcode the name!
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

		var pathToUse string
		relativePath := filepath.Join(filepath.Dir(file.URI), childPath)
		e.log.Debugf("Relative: %s | Child: %s", relativePath, childPath)
		e.log.Debugf("Trying relative path %s", relativePath)
		if util.FileExists(relativePath) {
			pathToUse = relativePath
		} else if util.FileExists(childPath) { // Root
			pathToUse = childPath
		}
		if pathToUse == "" {
			e.log.Errorf("Unable to match %s", childPath)
			return false
		}

		uri, err := util.PathToURI(pathToUse)
		if err != nil {
			e.log.Errorf("%s", err)
			return false
		}
		// e.log.Debugf("Produced URI %s", uri)

		if child := e.AddFile(uri, file); child != nil {
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
