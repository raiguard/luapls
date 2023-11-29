package lsp

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/tliron/commonlog"
)

type Env struct {
	config  EnvConfig
	globals map[string]*ast.Identifier
	log     commonlog.Logger
	name    string
}

type EnvConfig struct {
	Roots []string `json:"roots"`
}

func newEnv(name string, config EnvConfig) *Env {
	return &Env{
		config:  config,
		globals: map[string]*ast.Identifier{},
		log:     commonlog.GetLogger("env." + name),
	}
}

func (e *Env) getFiles() []string {
	e.log.Debug("Finding files")
	var files []string
	for _, path := range e.config.Roots {
		stat, err := os.Stat(path)
		if err != nil {
			e.log.Errorf("Failed to initialize: %s", err)
			continue
		}
		if stat.IsDir() {
			filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
				if !strings.HasSuffix(path, ".lua") {
					return nil
				}
				uri, err := pathToURI(path)
				if err != nil {
					e.log.Errorf("%s", err)
					return nil
				}
				files = append(files, uri)
				return nil
			})
		} else {
			uri, err := pathToURI(path)
			if err != nil {
				e.log.Errorf("%s", err)
				return nil
			}
			files = append(files, uri)
		}
	}
	e.log.Debugf("Found files: %s", toJSON(files))
	return files
}
