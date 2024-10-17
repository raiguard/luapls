package types

import (
	"github.com/raiguard/luapls/lua/ast"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type FileGraph struct {
	Roots []*FileNode
	Files map[protocol.URI]*FileNode
}

func (fg *FileGraph) Traverse(visitor func(fn *FileNode) bool) {
	for _, root := range fg.Roots {
		root.Traverse(visitor)
	}
	for _, file := range fg.Files {
		file.Visited = false
	}
}

// TODO: Atomics to allow multithreading
type FileNode struct {
	// AST is discarded after type checking is complete, unless the file is open in the editor.
	AST        *ast.Block
	LineBreaks []int

	Diagnostics []ast.Error
	Types       []*Type

	Parents  []*FileNode
	Children []*FileNode

	Path    string
	Visited bool
}

func (fn *FileNode) Traverse(visitor func(fn *FileNode) bool) {
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
