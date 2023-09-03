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
	Token   token.Token
	Name    *Identifier
	Value   Expression
	isLocal bool
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Name.Literal, as.Value.String())
}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence Block
}

func (ls *IfStatement) statementNode() {}
func (ls *IfStatement) String() string {
	return fmt.Sprintf("%s %s then\n%s\nend", ls.Token.Literal, ls.Condition.String(), ls.Consequence.String())
}

type LocalStatement struct {
	Token     token.Token
	Statement Statement
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("%s %s", ls.Token.Literal, ls.Statement.String())
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier token.Token

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Literal }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string // Without quotes
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Token.Literal }
