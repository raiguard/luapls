package types

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

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
			// TODO: Proper parser for type annotations
			content = strings.TrimSpace(content)
			content, ok = strings.CutPrefix(content, "@class")
			if !ok {
				continue
			}
			content = strings.TrimSpace(content)
			if len(content) == 0 {
				file.Diagnostics = append(file.Diagnostics, ast.Diagnostic{
					Message:  "Missing class name",
					Range:    trivia.Range(),
					Severity: protocol.DiagnosticSeverityWarning,
				})
				continue
			}
			parts := strings.Split(content, " ")
			e.log.Debugf("%d", len(parts))
			if len(parts) == 0 {
				file.Diagnostics = append(file.Diagnostics, ast.Diagnostic{
					Message:  "Missing class name",
					Range:    trivia.Range(),
					Severity: protocol.DiagnosticSeverityWarning,
				})
				continue
			}
			name := parts[0]
			if e.Types[name] != nil {
				continue
			}
			// TODO: Narrow location to actual name
			// TODO: Support multiple definition locations
			e.Types[name] = &Named{Name: name, Range: trivia.Range()}
		}
		return true
	})
}
