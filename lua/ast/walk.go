package ast

import (
	"fmt"

	"github.com/raiguard/luapls/lua/token"
)

type Visitor func(n Node) bool

// Walk performs a depth-first traversal of the AST, calling the visitor for each node.
// If the visitor returns false, this node's children are not traversed.
func Walk(n Node, v Visitor) {
	if n == nil || !v(n) {
		return
	}

	switch n := n.(type) {
	case *AssignmentStatement:
		WalkList(n.Vars, v)
		WalkList(n.Exps, v)

	case *Block:
		WalkList(n.Stmts, v)

	case *BooleanLiteral:
		// Leaf

	case *BreakStatement:
		// Leaf

	case *DoStatement:
		Walk(&n.Body, v)

	case *ExpressionStatement:
		Walk(n.Exp, v)

	case *ForInStatement:
		WalkList(n.Names, v)
		WalkList(n.Exps, v)
		Walk(&n.Body, v)

	case *ForStatement:
		Walk(n.Start, v)
		Walk(n.Finish, v)
		Walk(n.Step, v)
		Walk(&n.Body, v)

	case *FunctionCall:
		Walk(n.Left, v)
		WalkList(n.Args, v)

	case *FunctionExpression:
		WalkList(n.Params, v)
		Walk(&n.Body, v)

	case *FunctionStatement:
		Walk(n.Left, v)
		WalkList(n.Params, v)
		Walk(&n.Body, v)

	case *GotoStatement:
		Walk(n.Name, v)

	case *Identifier:
		// Leaf

	case *IfClause:
		Walk(n.Condition, v)
		Walk(&n.Body, v)

	case *IfStatement:
		WalkList(n.Clauses, v)

	case *IndexExpression:
		Walk(n.Left, v)
		Walk(n.Inner, v)

	case *InfixExpression:
		Walk(n.Left, v)
		Walk(n.Right, v)

	case *Invalid:
		// Leaf

	case *LabelStatement:
		// Leaf

	case *LocalStatement:
		WalkList(n.Names, v)
		WalkList(n.Exps, v)

	case *NumberLiteral:
		// Leaf

	case *PrefixExpression:
		Walk(n.Right, v)

	case *RepeatStatement:
		Walk(&n.Body, v)
		Walk(n.Condition, v)

	case *ReturnStatement:
		WalkList(n.Exps, v)

	case *StringLiteral:
		// Leaf

	case *TableField:
		Walk(n.Key, v)
		Walk(n.Value, v)

	case *TableLiteral:
		WalkList(n.Fields, v)

	case *Vararg:
		// Leaf

	case *WhileStatement:
		Walk(n.Condition, v)
		Walk(&n.Body, v)

	default:
		panic(fmt.Sprintf("Walk unimplemented for %T", n))
	}
}

// WalkList traverses a slice of Nodes.
func WalkList[N Node](nodes []N, v Visitor) {
	for i := 0; i < len(nodes); i++ {
		Walk(nodes[i], v)
	}
}

// GetNode returns the innermost node at the given position, and its parent nodes.
func GetNode(base Node, pos token.Pos) (Node, []Node) {
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
	return node, parents
}
