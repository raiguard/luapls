package ast

import (
	"reflect"

	"github.com/raiguard/luapls/lua/token"
)

type Visitor func(node Node) bool

// WalkSemantic performs a depth-first traversal of the AST, calling the visitor for each node.
// If the visitor returns false, this node's children are not traversed.
func WalkSemantic(node Node, visitor Visitor) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}
	if !visitor(node) {
		return
	}

	for _, child := range node.GetSemanticChildren() {
		WalkSemantic(child, visitor)
	}
}

type NodePath struct {
	Node    Node
	Parents []Node
}

// GetSemanticNode returns the innermost node at the given position, and its parent nodes.
func GetSemanticNode(base Node, pos token.Pos) NodePath {
	var node Node
	parents := []Node{}
	WalkSemantic(base, func(n Node) bool {
		if n.Pos() <= pos && pos < n.End() {
			if node != nil {
				parents = append(parents, node)
			}
			node = n
			return true
		}
		return false
	})
	return NodePath{node, parents}
}
