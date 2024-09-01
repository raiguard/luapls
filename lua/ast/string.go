package ast

import "fmt"

func (node *AssignmentStatement) String() string {
	return fmt.Sprintf("%s()@%v", "AssignmentStatement", node.Pos())
}

func (node *Block) String() string {
	return fmt.Sprintf("%s()@%v", "Block", node.Pos())
}

func (node *BooleanLiteral) String() string {
	return fmt.Sprintf("%s(%v)@%v", "BooleanLiteral", node.Value, node.Pos())
}

func (node *BreakStatement) String() string {
	return fmt.Sprintf("%s()@%v", "BreakStatement", node.Pos())
}

func (node *DoStatement) String() string {
	return fmt.Sprintf("%s()@%v", "DoStatement", node.Pos())
}

func (node *ForInStatement) String() string {
	return fmt.Sprintf("%s()@%v", "ForInStatement", node.Pos())
}

func (node *ForStatement) String() string {
	return fmt.Sprintf("%s()@%v", "ForStatement", node.Pos())
}

func (node *FunctionCall) String() string {
	return fmt.Sprintf("%s()@%v", "FunctionCall", node.Pos())
}

func (node *FunctionExpression) String() string {
	return fmt.Sprintf("%s()@%v", "FunctionExpression", node.Pos())
}

func (node *FunctionStatement) String() string {
	return fmt.Sprintf("%s()@%v", "FunctionStatement", node.Pos())
}

func (node *GotoStatement) String() string {
	return fmt.Sprintf("%s()@%v", "GotoStatement", node.Pos())
}

func (node *Identifier) String() string {
	return fmt.Sprintf("%s(%s)@%v", "Identifier", node.Literal, node.Pos())
}

func (node *IfClause) String() string {
	return fmt.Sprintf("%s()@%v", "IfClause", node.Pos())
}

func (node *IfStatement) String() string {
	return fmt.Sprintf("%s()@%v", "IfStatement", node.Pos())
}

func (node *IndexExpression) String() string {
	return fmt.Sprintf("%s()@%v", "IndexExpression", node.Pos())
}

func (node *InfixExpression) String() string {
	return fmt.Sprintf("%s()@%v", "InfixExpression", node.Pos())
}

func (node *LabelStatement) String() string {
	return fmt.Sprintf("%s(%s)@%v", "LabelStatement", node.Name, node.Pos())
}

func (node *Invalid) String() string {
	return fmt.Sprintf("%s()@%v", "Invalid", node.Pos())
}

func (node *LocalStatement) String() string {
	return fmt.Sprintf("%s()@%v", "LocalStatement", node.Pos())
}

func (node *NilLiteral) String() string {
	return fmt.Sprintf("%s()@%v", "NilLiteral", node.Pos())
}

func (node *NumberLiteral) String() string {
	return fmt.Sprintf("%s(%s)@%v", "NumberLiteral", node.Literal, node.Pos())
}

func (node *PrefixExpression) String() string {
	return fmt.Sprintf("%s()@%v", "PrefixExpression", node.Pos())
}

func (node *RepeatStatement) String() string {
	return fmt.Sprintf("%s()@%v", "RepeatStatement", node.Pos())
}

func (node *ReturnStatement) String() string {
	return fmt.Sprintf("%s()@%v", "ReturnStatement", node.Pos())
}

func (node *StringLiteral) String() string {
	return fmt.Sprintf("%s(%s)@%v", "StringLiteral", node.Unit.Token.Literal, node.Pos())
}

func (node *TableField) String() string {
	return fmt.Sprintf("%s()@%v", "TableField", node.Pos())
}

func (node *TableLiteral) String() string {
	return fmt.Sprintf("%s()@%v", "TableLiteral", node.Pos())
}

func (node *Vararg) String() string {
	return fmt.Sprintf("%s()@%v", "Vararg", node.Pos())
}

func (node *WhileStatement) String() string {
	return fmt.Sprintf("%s()@%v", "WhileStatement", node.Pos())
}
