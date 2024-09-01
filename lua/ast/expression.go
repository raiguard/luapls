package ast

import (
	"strconv"

	"github.com/raiguard/luapls/lua/token"
)

type Expression interface {
	Node
	expressionNode()
}

type FunctionCall struct {
	Name       Expression
	LeftParen  *Unit // Optional
	Args       []Expression
	RightParen *Unit // Optional
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) statementNode()  {}
func (fc *FunctionCall) Pos() token.Pos {
	return fc.Name.Pos()
}
func (fc *FunctionCall) End() token.Pos {
	if fc.RightParen != nil {
		return fc.RightParen.End()
	}
	if len(fc.Args) > 0 {
		return fc.Args[len(fc.Args)-1].End()
	}
	if fc.LeftParen != nil {
		return fc.LeftParen.End()
	}
	return fc.Name.End()
}

type FunctionExpression struct {
	Function   Unit
	LeftParen  Unit
	Params     []*Identifier
	Vararg     *Unit
	RightParen Unit
	Body       Block
	EndUnit    Unit
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) Pos() token.Pos {
	return fe.Function.Pos()
}
func (fe *FunctionExpression) End() token.Pos {
	return fe.EndUnit.End()
}

type IndexExpression struct {
	Left    Expression
	Indexer token.TokenType
	Inner   Expression
	EndPos  token.Pos `json:"-"`
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) Pos() token.Pos {
	return ie.Left.Pos()
}
func (ie *IndexExpression) End() token.Pos {
	return ie.EndPos
}

type InfixExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func (be *InfixExpression) expressionNode() {}
func (be *InfixExpression) Pos() token.Pos {
	return be.Left.Pos()
}
func (be *InfixExpression) End() token.Pos {
	return be.Right.End()
}

type PrefixExpression struct {
	Operator token.TokenType
	Right    Expression
	StartPos token.Pos `json:"-"`
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) Pos() token.Pos {
	return pe.StartPos
}
func (pe *PrefixExpression) End() token.Pos {
	return pe.Right.End()
}

// Literals (also expressions)

type BooleanLiteral struct {
	Value    bool
	StartPos token.Pos `json:"-"`
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) Pos() token.Pos {
	return bl.StartPos
}
func (bl *BooleanLiteral) End() token.Pos {
	return bl.StartPos + len(strconv.FormatBool(bl.Value))
}
func (bl *BooleanLiteral) leaf() {}

type Identifier struct {
	Literal  string
	StartPos token.Pos `json:"-"`
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Pos() token.Pos {
	return i.StartPos
}
func (i *Identifier) End() token.Pos {
	return i.StartPos + len(i.Literal)
}
func (i *Identifier) leaf() {}

type NilLiteral struct {
	StartPos token.Pos `json:"-"`
}

func (i *NilLiteral) expressionNode() {}
func (i *NilLiteral) Pos() token.Pos {
	return i.StartPos
}
func (i *NilLiteral) End() token.Pos {
	return i.StartPos + len("nil")
}
func (n *NilLiteral) leaf() {}

type NumberLiteral struct {
	Literal  string
	Value    float64
	StartPos token.Pos `json:"-"`
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) Pos() token.Pos {
	return nl.StartPos
}
func (nl *NumberLiteral) End() token.Pos {
	return nl.StartPos + len(nl.Literal)
}
func (nl *NumberLiteral) leaf() {}

type StringLiteral struct {
	Unit Unit
	// TODO: Store type of quote
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) Pos() token.Pos {
	return sl.Unit.Pos()
}
func (sl *StringLiteral) End() token.Pos {
	return sl.Unit.End()
}
func (sl *StringLiteral) leaf() {}

type TableLiteral struct {
	Fields   []*TableField
	StartPos token.Pos
	EndPos   token.Pos `json:"-"`
}

func (tl *TableLiteral) expressionNode() {}
func (tl *TableLiteral) Pos() token.Pos {
	return tl.StartPos
}
func (tl *TableLiteral) End() token.Pos {
	return tl.EndPos
}

type Vararg struct {
	StartPos token.Pos `json:"-"`
}

func (va *Vararg) expressionNode() {}
func (va *Vararg) Pos() token.Pos {
	return va.StartPos
}
func (va *Vararg) End() token.Pos {
	return va.StartPos + len(token.VARARG.String())
}
func (va *Vararg) leaf() {}
