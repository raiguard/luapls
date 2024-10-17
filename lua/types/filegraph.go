package types

import (
	"github.com/raiguard/luapls/lua/ast"
)

// FileGraph specifies a set of files and the relationships between them.
type FileGraph struct {
	Files map[string]*File
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
	LineBreaks []int

	Diagnostics []ast.Error
	Types       []*Type

	Parents  []*File
	Children []*File

	Path    string
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
