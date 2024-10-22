package types

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// FileGraph specifies a set of files and the relationships between them.
type FileGraph struct {
	Files map[protocol.URI]*File
	// Roots are all files with no parent files.
	Roots []*File
}

func (fg *FileGraph) Traverse(visitor func(fn *File) bool) {
	for _, root := range fg.Roots {
		root.Traverse(visitor)
	}
	for _, file := range fg.Files {
		file.Visited = false
	}
}

// TODO: Atomics to allow multithreading
type File struct {
	// AST is discarded after type checking is complete, unless the file is open in the editor.
	AST        *ast.Block
	LineBreaks token.LineBreaks

	Diagnostics []ast.Diagnostic

	Parents  []*File
	Children []*File

	URI     protocol.URI
	Visited bool
}

func (fn *File) Traverse(visitor func(fn *File) bool) {
	if fn == nil || fn.Visited {
		return
	}
	fn.Visited = true
	if visitor(fn) {
		for _, child := range fn.Children {
			child.Traverse(visitor)
		}
	}
}
