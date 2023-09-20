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
	Vars []Expression
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

type ExpressionStatement struct {
	Exp Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%s", es.Exp.String())
}

type ForStatement struct {
	Name  *Identifier
	Start Expression
	End   Expression
	Step  Expression
	Body  Block
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	if fs.Step != nil {
		return fmt.Sprintf(
			"for %s = %s, %s, %s do\n%s\nend",
			fs.Name.String(),
			fs.Start.String(),
			fs.End.String(),
			fs.Step.String(),
			fs.Body.String(),
		)
	} else {
		return fmt.Sprintf(
			"for %s = %s, %s do\n%s\nend",
			fs.Name.String(),
			fs.Start.String(),
			fs.End.String(),
			fs.Body.String(),
		)
	}
}

type ForInStatement struct {
	Names []*Identifier
	Exps  []Expression
	Body  Block
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", nodeListToString(fs.Names), nodeListToString(fs.Exps), fs.Body.String())
}

type FunctionStatement struct {
	Left    Expression
	Params  []*Identifier
	Vararg  bool
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
		fs.Left.String(),
		nodeListToString(fs.Params),
		fs.Body.String(),
	)
}

type GotoStatement struct {
	Name *Identifier
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("goto %s", gs.Name.String())
}

type IfStatement struct {
	Clauses []*IfClause
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("%send", nodeListToString(is.Clauses))
}

type IfClause struct {
	Condition Expression
	Body      Block
}

func (ic *IfClause) statementNode() {}
func (ic *IfClause) String() string {
	// TODO: elseif
	if ic.Condition != nil {
		return fmt.Sprintf("if %s then\n%s\n", ic.Condition.String(), ic.Body.String())
	} else {
		return fmt.Sprintf("else\n%s\n", ic.Body.String())
	}
}

type LabelStatement struct {
	Name *Identifier
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return fmt.Sprintf("::%s::", ls.Name.String())
}

type LocalStatement struct {
	Names []*Identifier
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

type FunctionCall struct {
	Left Expression
	Args []Expression
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", fc.Left.String(), nodeListToString(fc.Args))
}

type FunctionExpression struct {
	Params []*Identifier
	Vararg bool
	Body   Block
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) String() string {
	return fmt.Sprintf("function(%s)\n%s\nend", nodeListToString(fe.Params), fe.Body.String())
}

type IndexExpression struct {
	Left       Expression
	Inner      Expression
	IsBrackets bool
	IsColon    bool
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	if ie.IsBrackets {
		return fmt.Sprintf("%s[%s]", ie.Left.String(), ie.Inner.String())
	} else {
		return fmt.Sprintf("%s.%s", ie.Left.String(), ie.Inner.String())
	}
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

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Literal }

type NumberLiteral struct {
	Literal string
	Value   float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Literal }

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }

type TableLiteral struct {
	Fields []*TableField
}

func (tl *TableLiteral) expressionNode() {}
func (tl *TableLiteral) String() string  { return fmt.Sprintf("{ %s }", nodeListToString(tl.Fields)) }

type TableField struct {
	Key   Expression
	Value Expression
}

func (tf TableField) String() string {
	if tf.Key == nil {
		return tf.Value.String()
	}
	if ident, ok := tf.Key.(*Identifier); ok {
		return fmt.Sprintf("%s = %s", ident.String(), tf.Value.String())
	}
	return fmt.Sprintf("[%s] = %s", tf.Key.String(), tf.Value.String())
}

type Vararg struct{}

func (va *Vararg) expressionNode() {}
func (va *Vararg) String() string  { return token.TokenStr[token.VARARG] }
