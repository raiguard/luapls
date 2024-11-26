package types

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/raiguard/luapls/lua/annotation"
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/util"
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Environment struct {
	Files    map[protocol.URI]*ast.File
	RootPath string

	Types map[string]Type

	log commonlog.Logger
}

func NewEnvironment() *Environment {
	return &Environment{
		Files: map[protocol.URI]*ast.File{},
		Types: map[string]Type{},
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
			e.AddFile(uri)
		}
		return nil
	})
	e.CheckPhase1()
	e.log.Debugf("Initialization took %s", time.Since(before).String())

	e.log.Debug("TYPES:")
	for name := range e.Types {
		e.log.Debug(name)
	}
}

func (e *Environment) AddFile(uri protocol.URI) *ast.File {
	if existing := e.Files[uri]; existing != nil {
		return existing
	}
	path, err := util.URIToPath(uri)
	if err != nil {
		e.log.Errorf("%s", err)
		return nil
	}
	src, err := os.ReadFile(path)
	if err != nil {
		e.log.Errorf("Failed to parse file %s: %s", path, err)
		return nil
	}
	timer := time.Now()
	file := util.Ptr(parser.New(string(src)).ParseFile())
	e.log.Debugf("Parsed file '%s' in %s", path, time.Since(timer).String())
	file.URI = uri
	e.Files[uri] = file

	return file
}

func (e *Environment) AddTransientFile(uri protocol.URI, content string) *ast.File {
	if existing := e.Files[uri]; existing != nil {
		return existing
	}
	path, err := util.URIToPath(uri)
	if err != nil {
		e.log.Errorf("%s", err)
		return nil
	}
	timer := time.Now()
	file := util.Ptr(parser.New(content).ParseFile())
	e.log.Debugf("Parsed file '%s' in %s", path, time.Since(timer).String())
	file.URI = uri
	e.Files[uri] = file

	return file
}

// CheckPhase1 executes the first phase of type checking.
// The first phase gathers a list of which types exist in the environment, but does not delve into details.
func (e *Environment) CheckPhase1() {
	for _, file := range e.Files {
		e.CheckFilePhase1(file)
	}
}

func (e *Environment) CheckFilePhase1(file *ast.File) {
	ast.WalkSemantic(file.Block, func(n ast.Node) bool {
		for _, trivia := range n.GetLeadingTrivia() {
			if trivia.Type != token.COMMENT {
				continue
			}
			content, ok := strings.CutPrefix(trivia.Literal, "---")
			if !ok {
				continue
			}
			a, diags := annotation.Parse(content)
			for _, diag := range diags {
				diag.Range.Start += trivia.Pos + 3
				diag.Range.End += trivia.Pos + 3
				file.Diagnostics = append(file.Diagnostics, diag)
			}
			if a == nil {
				continue
			}
			class, ok := a.(*annotation.Class)
			if !ok {
				continue
			}
			if e.Types[class.Name] != nil {
				continue
			}
			// TODO: Narrow location to actual name
			// TODO: Support multiple definition locations
			e.Types[class.Name] = &Named{Name: class.Name, Range: trivia.Range()}
		}
		return true
	})
}
