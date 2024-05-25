package lsp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Env struct {
	config  EnvConfig
	globals map[string]*ast.Identifier
	log     commonlog.Logger
	name    string
}

type EnvConfig struct {
	Libraries []string `json:"libraries"`
	Roots     []string `json:"roots"`
}

func newEnv(name string, config EnvConfig) *Env {
	return &Env{
		config:  config,
		globals: map[string]*ast.Identifier{},
		log:     commonlog.GetLogger("env." + name),
	}
}

func (e *Env) getFiles() []protocol.URI {
	e.log.Debug("Finding files")
	var files []protocol.URI
	for _, path := range e.config.Libraries {
		files = e.addPath(files, path)
	}
	for _, path := range e.config.Roots {
		files = e.addPath(files, path)
	}
	e.log.Debugf("Found files: %s", toJSON(files))
	return files
}

func (e *Env) addPath(files []protocol.URI, path string) []protocol.URI {
	stat, err := os.Stat(path)
	if err != nil {
		e.log.Errorf("Failed to initialize: %s", err)
		return files
	}
	if stat.IsDir() {
		filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
			if !strings.HasSuffix(path, ".lua") {
				return nil
			}
			files = e.addFile(files, path)
			return nil
		})
	} else {
		files = e.addFile(files, path)
	}
	return files
}

func (e *Env) addFile(files []protocol.URI, path string) []protocol.URI {
	uri, err := pathToURI(path)
	if err != nil {
		e.log.Errorf("%s", err)
		return nil
	}
	return append(files, uri)
}
