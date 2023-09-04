package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	String() string
}

type Block []Statement

func (b *Block) String() string {
	var out string
	for _, stmt := range *b {
		out += stmt.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func nodeListToString[T Node](nodes []T) string {
	items := []string{}
	for _, node := range nodes {
		items = append(items, node.String())
	}
	return strings.Join(items, ", ")
}

// Statements

type Statement interface {
	Node
	statementNode()
}

type AssignmentStatement struct {
	Vars []Identifier
	Exps []Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", nodeListToString(as.Vars), nodeListToString(as.Exps))
}

type BreakStatement struct{}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return "break"
}

type DoStatement struct {
	Body Block
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) String() string {
	return fmt.Sprintf("do\n%s\nend", ds.Body.String())
}

type ForStatement struct {
	Var   Identifier
	Start Expression
	End   Expression
	Step  *Expression // Optional
	Body  Block
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	if fs.Step != nil {
		return fmt.Sprintf(
			"for %s = %s, %s, %s do\n%s\nend",
			fs.Var.String(),
			fs.Start.String(),
			fs.End.String(),
			(*fs.Step).String(),
			fs.Body.String(),
		)
	} else {
		return fmt.Sprintf(
			"for %s = %s, %s do\n%s\nend",
			fs.Var.String(),
			fs.Start.String(),
			fs.End.String(),
			fs.Body.String(),
		)
	}
}

type ForInStatement struct {
	Vars []Identifier
	Exps []Expression
	Body Block
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", nodeListToString(fs.Vars), nodeListToString(fs.Exps), fs.Body.String())
}

type FunctionStatement struct {
	Name    Identifier
	Params  []Identifier
	Body    Block
	IsLocal bool
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) String() string {
	localStr := ""
	if fs.IsLocal {
		localStr = "local "
	}
	return fmt.Sprintf(
		"%sfunction %s(%s)\n%s\nend",
		localStr,
		fs.Name.String(),
		nodeListToString(fs.Params),
		fs.Body.String(),
	)
}

type GotoStatement struct {
	Label Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("goto %s", gs.Label.String())
}

type IfStatement struct {
	Clauses []IfClause
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("%send", nodeListToString(is.Clauses))
}

type IfClause struct {
	Condition Expression
	Body      Block
}

func (ic IfClause) statementNode() {}
func (ic IfClause) String() string {
	return fmt.Sprintf("if %s then\n%s\n", ic.Condition.String(), ic.Body.String())
}

type LabelStatement struct {
	Label Identifier
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return fmt.Sprintf("::%s::", ls.Label.String())
}

type LocalStatement struct {
	Names []Identifier
	Exps  []Expression
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("local %s = %s", nodeListToString(ls.Names), nodeListToString(ls.Exps))
}

type RepeatStatement struct {
	Body      Block
	Condition Expression
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) String() string {
	return fmt.Sprintf("repeat\n%s\nuntil %s", rs.Body.String(), rs.Condition.String())
}

type ReturnStatement struct {
	Exps []Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", nodeListToString(rs.Exps))
}

type WhileStatement struct {
	Condition Expression
	Body      Block
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s do\n%s\nend", ws.Condition.String(), ws.Body.String())
}

// Expressions

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func (ie *BinaryExpression) expressionNode() {}
func (ie *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator.String(), ie.Right.String())
}

type UnaryExpression struct {
	Operator token.TokenType
	Right    Expression
}

func (pe *UnaryExpression) expressionNode() {}
func (pe *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// Literals (also expressions)

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

type Identifier struct {
	Literal string
}

func (i Identifier) expressionNode() {}
func (i Identifier) String() string  { return i.Literal }

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return fmt.Sprintf("%.f", nl.Value) }

type StringLiteral struct {
	Value string // Without quotes
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Value }

// Other

// FunctionCall can be both a statement and an expression.
type FunctionCall struct {
	Name Identifier
	Args []Expression
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) statementNode()  {}
func (fc *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", fc.Name.String(), nodeListToString(fc.Args))
}
