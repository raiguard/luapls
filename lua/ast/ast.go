package ast

import (
	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
	String() string
	GetLeadingTrivia() []token.Token
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

func (u *Unit) Range() token.Range {
	return token.Range{Start: u.Pos(), End: u.End()}
}

func (u *Unit) String() string {
	return u.Token.Literal
}

type Block = Punctuated[Statement]

type Invalid struct {
	Position token.Pos
}

func (i *Invalid) expressionNode() {}
func (i *Invalid) statementNode()  {}
func (i *Invalid) Pos() token.Pos {
	return i.Position
}
func (i *Invalid) End() token.Pos {
	return i.Position
}

// A Leaf node has no children and is interactable in the editor.
type LeafNode interface {
	leaf()
}
