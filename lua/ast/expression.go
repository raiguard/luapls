package ast

import (
	"fmt"

	"github.com/raiguard/luapls/lua/token"
)

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator.String(), be.Right.String())
}
func (be *BinaryExpression) Pos() token.Pos {
	return be.Left.Pos()
}
func (be *BinaryExpression) End() token.Pos {
	return be.Right.End()
}

type FunctionCall struct {
	Left   Expression
	Args   []Expression
	EndPos token.Pos
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", fc.Left.String(), nodeListToString(fc.Args))
}
func (fc *FunctionCall) Pos() token.Pos {
	return fc.Left.Pos()
}
func (fc *FunctionCall) End() token.Pos {
	return fc.EndPos
}

type FunctionExpression struct {
	Params   []*Identifier
	Vararg   bool
	Body     Block
	StartPos token.Pos
	EndPos   token.Pos
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) String() string {
	return fmt.Sprintf("function(%s)\n%s\nend", nodeListToString(fe.Params), fe.Body.String())
}
func (fe *FunctionExpression) Pos() token.Pos {
	return fe.StartPos
}
func (fe *FunctionExpression) End() token.Pos {
	return fe.EndPos
}

type IndexExpression struct {
	Left       Expression
	Inner      Expression
	IsBrackets bool
	IsColon    bool
	EndPos     token.Pos
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	if ie.IsBrackets {
		return fmt.Sprintf("%s[%s]", ie.Left.String(), ie.Inner.String())
	} else {
		return fmt.Sprintf("%s.%s", ie.Left.String(), ie.Inner.String())
	}
}
func (ie *IndexExpression) Pos() token.Pos {
	return ie.Left.Pos()
}
func (ie *IndexExpression) End() token.Pos {
	return ie.EndPos
}

type UnaryExpression struct {
	Operator token.TokenType
	Right    Expression
	StartPos token.Pos
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", ue.Operator, ue.Right.String())
}
func (ue *UnaryExpression) Pos() token.Pos {
	return ue.StartPos
}
func (ue *UnaryExpression) End() token.Pos {
	return ue.Right.End()
}

// Literals (also expressions)

type BooleanLiteral struct {
	Value    bool
	StartPos token.Pos
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}
func (bl *BooleanLiteral) Pos() token.Pos {
	return bl.StartPos
}
func (bl *BooleanLiteral) End() token.Pos {
	return bl.StartPos + len(bl.String())
}
func (bl *BooleanLiteral) leaf() {}

type Identifier struct {
	Literal  string
	StartPos token.Pos
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Literal }
func (i *Identifier) Pos() token.Pos {
	return i.StartPos
}
func (i *Identifier) End() token.Pos {
	return i.StartPos + len(i.String())
}
func (i *Identifier) leaf() {}

type NumberLiteral struct {
	Literal  string
	Value    float64
	StartPos token.Pos
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Literal }
func (nl *NumberLiteral) Pos() token.Pos {
	return nl.StartPos
}
func (nl *NumberLiteral) End() token.Pos {
	return nl.StartPos + len(nl.String())
}
func (nl *NumberLiteral) leaf() {}

type StringLiteral struct {
	Value    string
	StartPos token.Pos
	// TODO: Store type of quote
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }
func (sl *StringLiteral) Pos() token.Pos {
	return sl.StartPos
}
func (sl *StringLiteral) End() token.Pos {
	return sl.StartPos + len(sl.String())
}
func (sl *StringLiteral) leaf() {}

type TableLiteral struct {
	Fields   []*TableField
	StartPos token.Pos
	EndPos   token.Pos
}

func (tl *TableLiteral) expressionNode() {}
func (tl *TableLiteral) String() string  { return fmt.Sprintf("{ %s }", nodeListToString(tl.Fields)) }
func (tl *TableLiteral) Pos() token.Pos {
	return tl.StartPos
}
func (tl *TableLiteral) End() token.Pos {
	return tl.EndPos
}

type Vararg struct {
	StartPos token.Pos
}

func (va *Vararg) expressionNode() {}
func (va *Vararg) String() string  { return token.TokenStr[token.VARARG] }
func (va *Vararg) Pos() token.Pos {
	return va.StartPos
}
func (va *Vararg) End() token.Pos {
	return va.StartPos + len(token.VARARG.String())
}
func (va *Vararg) leaf() {}
