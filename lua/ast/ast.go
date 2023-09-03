package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	String() string
}

type Block struct {
	Node
	Statements []Statement
}

func (b *Block) String() string {
	var out string
	for _, stmt := range b.Statements {
		out += stmt.String() + "\n"
	}
	return strings.TrimSpace(out)
}

type Statement interface {
	Node
	statementNode()
}

type AssignmentStatement struct {
	Token token.Token
	Vars  []Identifier
	Exps  []Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", nodeListToString(as.Vars), nodeListToString(as.Exps))
}

type BreakStatement token.Token

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return bs.Literal
}

type DoStatement struct {
	Token token.Token
	Block Block
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) String() string {
	return fmt.Sprintf("%s\n%s\nend", ds.Token.Literal, ds.Block.String())
}

type GotoStatement struct {
	Token token.Token
	Label Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("%s %s", gs.Token.Literal, gs.Label.String())
}

type IfStatement struct {
	Token     token.Token
	Condition Expression
	Block     Block
}

func (ls *IfStatement) statementNode() {}
func (ls *IfStatement) String() string {
	return fmt.Sprintf("%s %s then\n%s\nend", ls.Token.Literal, ls.Condition.String(), ls.Block.String())
}

type LocalStatement struct {
	Token token.Token
	Names []Identifier
	Exps  []Expression
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("%s %s = %s", ls.Token.Literal, nodeListToString(ls.Names), nodeListToString(ls.Exps))
}

type RepeatStatement struct {
	Token     token.Token
	Block     Block
	Condition Expression
}

func (ws *RepeatStatement) statementNode() {}
func (ws *RepeatStatement) String() string {
	return fmt.Sprintf("%s\n%s\nuntil %s", ws.Token.Literal, ws.Block.String(), ws.Condition.String())
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Block     Block
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("%s %s do\n%s\nend", ws.Token.Literal, ws.Condition.String(), ws.Block.String())
}

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ie *BinaryExpression) expressionNode() {}
func (ie *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Token.Literal, ie.Right.String())
}

type Identifier token.Token

func (i Identifier) expressionNode() {}
func (i Identifier) String() string  { return i.Literal }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Token.Literal }

type UnaryExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *UnaryExpression) expressionNode() {}
func (pe *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type StringLiteral struct {
	Token token.Token
	Value string // Without quotes
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Token.Literal }

func nodeListToString[T Node](nodes []T) string {
	items := []string{}
	for _, node := range nodes {
		items = append(items, node.String())
	}
	return strings.Join(items, ", ")
}
