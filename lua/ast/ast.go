package ast

import (
	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
	String() string
}

func Range(n Node) token.Range {
	return token.Range{Start: n.Pos(), End: n.End()}
}

type Unit struct {
	LeadingTrivia  []token.Token
	Token          token.Token
	TrailingTrivia []token.Token // Comments and whitespace up to the next newline
}

func (u *Unit) Type() token.TokenType {
	return u.Token.Type
}

func (u *Unit) Pos() token.Pos {
	return u.Token.Pos
}

func (u *Unit) End() token.Pos {
	return u.Token.End()
}

func (u *Unit) String() string {
	return u.Token.Literal
}

type Block = Punctuated[Statement]

type TableField struct {
	Key      Expression
	Value    Expression
	StartPos token.Pos `json:"-"` // Needed in case of bracketed keys
}

func (tf *TableField) Pos() token.Pos {
	return tf.StartPos
}
func (tf *TableField) End() token.Pos {
	return tf.Value.End()
}

type Invalid struct {
	Exps Punctuated[Expression] `json:",omitempty"`
	// OR
	Unit *Unit `json:",omitempty"`
}

func (i *Invalid) expressionNode() {}
func (i *Invalid) statementNode()  {}
func (i *Invalid) Pos() token.Pos {
	if i.Unit != nil {
		return i.Unit.Pos()
	}
	return i.Exps.Pairs[0].Pos()
}
func (i *Invalid) End() token.Pos {
	if i.Unit != nil {
		return i.Unit.End()
	}
	return i.Exps.Pairs[len(i.Exps.Pairs)-1].End()
}

// A Leaf node has no children and is interactable in the editor.
type LeafNode interface {
	leaf()
}
