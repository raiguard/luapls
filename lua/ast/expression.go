package ast

import (
	"github.com/raiguard/luapls/lua/token"
)

type Expression interface {
	Node
	expressionNode()
}

type FunctionCall struct {
	Name       Expression
	LeftParen  *Unit // Optional
	Args       Punctuated[Expression]
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
	if len(fc.Args.Pairs) > 0 {
		return fc.Args.Pairs[len(fc.Args.Pairs)-1].End()
	}
	if fc.LeftParen != nil {
		return fc.LeftParen.End()
	}
	return fc.Name.End()
}

type FunctionExpression struct {
	Function   Unit
	LeftParen  Unit
	Params     Punctuated[*Identifier]
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
	Prefix       Expression
	LeftIndexer  Unit
	Inner        Expression
	RightIndexer *Unit
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) Pos() token.Pos {
	return ie.Prefix.Pos()
}
func (ie *IndexExpression) End() token.Pos {
	if ie.RightIndexer != nil {
		return ie.RightIndexer.End()
	}
	return ie.Inner.End()
}

type InfixExpression struct {
	Left     Expression
	Operator Unit
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
	Operator Unit
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) Pos() token.Pos {
	return pe.Operator.Pos()
}
func (pe *PrefixExpression) End() token.Pos {
	return pe.Right.End()
}

// Literals (also expressions)

type BooleanLiteral Unit

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) Pos() token.Pos {
	return bl.Token.Pos
}
func (bl *BooleanLiteral) End() token.Pos {
	return bl.Token.End()
}
func (bl *BooleanLiteral) leaf() {}

type Identifier Unit

func (i *Identifier) expressionNode() {}
func (i *Identifier) Pos() token.Pos {
	return i.Token.Pos
}
func (i *Identifier) End() token.Pos {
	return i.Token.End()
}
func (i *Identifier) leaf() {}

type NilLiteral Unit

func (i *NilLiteral) expressionNode() {}
func (i *NilLiteral) Pos() token.Pos {
	return i.Token.Pos
}
func (i *NilLiteral) End() token.Pos {
	return i.Token.End()
}
func (n *NilLiteral) leaf() {}

type NumberLiteral Unit

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) Pos() token.Pos {
	return nl.Token.Pos
}
func (nl *NumberLiteral) End() token.Pos {
	return nl.Token.End()
}
func (nl *NumberLiteral) leaf() {}

type StringLiteral Unit

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) Pos() token.Pos {
	return sl.Token.Pos
}
func (sl *StringLiteral) End() token.Pos {
	return sl.Token.End()
}
func (sl *StringLiteral) leaf() {}

type TableLiteral struct {
	LeftBrace  Unit
	Fields     Punctuated[TableField]
	RightBrace Unit
}

func (tl *TableLiteral) expressionNode() {}
func (tl *TableLiteral) Pos() token.Pos {
	return tl.LeftBrace.Pos()
}
func (tl *TableLiteral) End() token.Pos {
	return tl.RightBrace.End()
}

type Vararg Unit

func (va *Vararg) expressionNode() {}
func (va *Vararg) Pos() token.Pos {
	return va.Token.Pos
}
func (va *Vararg) End() token.Pos {
	return va.Token.End()
}
func (va *Vararg) leaf() {}
