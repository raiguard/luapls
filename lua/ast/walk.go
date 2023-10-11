package ast

import (
	"fmt"
)

type Visitor func(n Node) bool

// Walk performs a depth-first traversal of the AST, calling the visitor for each node.
// If the visitor returns false, this node's children are not traversed.
func Walk(n Node, v Visitor) {
	if n == nil || !v(n) {
		return
	}

	switch n.(type) {
	case *AssignmentStatement:
		WalkList(n.(*AssignmentStatement).Vars, v)
		WalkList(n.(*AssignmentStatement).Exps, v)

	case *BinaryExpression:
		Walk(n.(*BinaryExpression).Left, v)
		Walk(n.(*BinaryExpression).Right, v)

	case *Block:
		WalkList(n.(*Block).Stmts, v)

	case *BooleanLiteral:
		// Leaf

	case *BreakStatement:
		// Leaf

	case *DoStatement:
		Walk(&n.(*DoStatement).Body, v)

	case *ExpressionStatement:
		Walk(n.(*ExpressionStatement).Exp, v)

	case *File:
		Walk(&n.(*File).Block, v)

	case *ForInStatement:
		WalkList(n.(*ForInStatement).Names, v)
		WalkList(n.(*ForInStatement).Exps, v)
		Walk(&n.(*ForInStatement).Body, v)

	case *ForStatement:
		Walk(n.(*ForStatement).Start, v)
		Walk(n.(*ForStatement).Finish, v)
		Walk(n.(*ForStatement).Step, v)
		Walk(&n.(*ForStatement).Body, v)

	case *FunctionCall:
		Walk(n.(*FunctionCall).Left, v)
		WalkList(n.(*FunctionCall).Args, v)

	case *FunctionExpression:
		WalkList(n.(*FunctionExpression).Params, v)
		Walk(&n.(*FunctionExpression).Body, v)

	case *FunctionStatement:
		Walk(n.(*FunctionStatement).Left, v)
		WalkList(n.(*FunctionStatement).Params, v)
		Walk(&n.(*FunctionStatement).Body, v)

	case *GotoStatement:
		Walk(n.(*GotoStatement).Name, v)

	case *Identifier:
		// Leaf

	case *IfClause:
		Walk(n.(*IfClause).Condition, v)
		Walk(&n.(*IfClause).Body, v)

	case *IfStatement:
		WalkList(n.(*IfStatement).Clauses, v)

	case *IndexExpression:
		Walk(n.(*IndexExpression).Left, v)
		Walk(n.(*IndexExpression).Inner, v)

	case *LabelStatement:
		// Leaf

	case *LocalStatement:
		WalkList(n.(*LocalStatement).Names, v)
		WalkList(n.(*LocalStatement).Exps, v)

	case *NumberLiteral:
		// Leaf

	case *RepeatStatement:
		Walk(&n.(*RepeatStatement).Body, v)
		Walk(n.(*RepeatStatement).Condition, v)

	case *ReturnStatement:
		WalkList(n.(*ReturnStatement).Exps, v)

	case *StringLiteral:
		// Leaf

	case *TableField:
		Walk(n.(*TableField).Key, v)
		Walk(n.(*TableField).Value, v)

	case *TableLiteral:
		WalkList(n.(*TableLiteral).Fields, v)

	case *UnaryExpression:
		Walk(n.(*UnaryExpression).Right, v)

	case *Vararg:
		// Leaf

	case *WhileStatement:
		Walk(n.(*WhileStatement).Condition, v)
		Walk(&n.(*WhileStatement).Body, v)

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
