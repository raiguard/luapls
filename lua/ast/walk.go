package ast

import (
	"fmt"
	"reflect"

	"github.com/raiguard/luapls/lua/token"
)

type Visitor func(node Node) bool

// Walk performs a depth-first traversal of the AST, calling the visitor for each node.
// If the visitor returns false, this node's children are not traversed.
func Walk(node Node, visitor Visitor) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}
	if !visitor(node) {
		return
	}

	switch node := node.(type) {
	case *AssignmentStatement:
		Walk(&node.Vars, visitor)
		Walk(&node.Exps, visitor)

	case *BooleanLiteral:
		// Leaf

	case *BreakStatement:
		// Leaf

	case *DoStatement:
		Walk(&node.Body, visitor)

	case *ForInStatement:
		Walk(&node.Names, visitor)
		Walk(&node.Exps, visitor)
		Walk(&node.Body, visitor)

	case *ForStatement:
		Walk(node.Name, visitor)
		Walk(&node.Start, visitor)
		Walk(&node.Finish, visitor)
		Walk(node.Step, visitor)
		Walk(&node.Body, visitor)

	case *FunctionCall:
		Walk(node.Name, visitor)
		Walk(&node.Args, visitor)

	case *FunctionExpression:
		Walk(&node.Params, visitor)
		Walk(&node.Body, visitor)

	case *FunctionStatement:
		Walk(node.Name, visitor)
		Walk(&node.Params, visitor)
		Walk(&node.Body, visitor)

	case *GotoStatement:
		Walk(node.Name, visitor)

	case *Identifier:
		// Leaf

	case *IfClause:
		Walk(node.Condition, visitor)
		Walk(&node.Body, visitor)

	case *IfStatement:
		WalkList(node.Clauses, visitor)

	case *IndexExpression:
		Walk(node.Prefix, visitor)
		Walk(node.Inner, visitor)

	case *InfixExpression:
		Walk(node.Left, visitor)
		Walk(node.Right, visitor)

	case *Invalid:
		// Leaf

	case *LabelStatement:
		// Leaf

	case *LocalStatement:
		Walk(&node.Names, visitor)
		Walk(node.Exps, visitor)

	case *NilLiteral:
		// Leaf

	case *NumberLiteral:
		// Leaf

	case *Pair[Node]:
		Walk(node.Node, visitor)

	case *Pair[Statement]:
		Walk(node.Node, visitor)

	case *Pair[Expression]:
		Walk(node.Node, visitor)

	case *Pair[*Identifier]:
		Walk(node.Node, visitor)

	case *Pair[TableField]:
		Walk(node.Node, visitor)

	case *Pair[*TableSimpleKeyField]:
		Walk(node.Node, visitor)

	case *Pair[*TableArrayField]:
		Walk(node.Node, visitor)

	case *Pair[*TableExpressionKeyField]:
		Walk(node.Node, visitor)

	case *PrefixExpression:
		Walk(node.Right, visitor)

	case *Punctuated[Node]:
		for i := 0; i < len(node.Pairs); i++ {
			Walk(&node.Pairs[i], visitor)
		}

	case *Punctuated[Statement]:
		for i := 0; i < len(node.Pairs); i++ {
			Walk(&node.Pairs[i], visitor)
		}

	case *Punctuated[Expression]:
		for i := 0; i < len(node.Pairs); i++ {
			Walk(&node.Pairs[i], visitor)
		}

	case *Punctuated[*Identifier]:
		for i := 0; i < len(node.Pairs); i++ {
			Walk(&node.Pairs[i], visitor)
		}

	case *Punctuated[TableField]:
		for i := 0; i < len(node.Pairs); i++ {
			Walk(&node.Pairs[i], visitor)
		}

	case *RepeatStatement:
		Walk(&node.Body, visitor)
		Walk(node.Condition, visitor)

	case *ReturnStatement:
		Walk(node.Exps, visitor)

	case *SemicolonStatement:
		// Leaf

	case *StringLiteral:
		// Leaf

	case *TableArrayField:
		Walk(node.Expr, visitor)

	case *TableExpressionKeyField:
		Walk(node.Name, visitor)
		Walk(node.Expr, visitor)

	case *TableSimpleKeyField:
		Walk(&node.Name, visitor)
		Walk(node.Expr, visitor)

	case *TableLiteral:
		Walk(&node.Fields, visitor)

	case *Vararg:
		// Leaf

	case *WhileStatement:
		Walk(node.Condition, visitor)
		Walk(&node.Body, visitor)

	default:
		panic(fmt.Sprintf("Walk unimplemented for %T", node))
	}
}

// WalkList traverses a slice of Nodes.
func WalkList[N Node](nodes []N, v Visitor) {
	for i := 0; i < len(nodes); i++ {
		Walk(nodes[i], v)
	}
}

type NodePath struct {
	Node    Node
	Parents []Node
}

// GetNode returns the innermost node at the given position, and its parent nodes.
func GetNode(base Node, pos token.Pos) NodePath {
	var node Node
	parents := []Node{}
	Walk(base, func(n Node) bool {
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
