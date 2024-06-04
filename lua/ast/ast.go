package ast

import (
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
	String() string

	GetComments() string
	// CommentsAfter() string
}

func Range(n Node) token.Range {
	return token.Range{Start: n.Pos(), End: n.End()}
}

type Comments struct {
	CommentsBefore []token.Token `json:",omitempty"`
}

// TODO: Remove this
func (c *Comments) GetComments() string {
	if c == nil || c.CommentsBefore == nil {
		return ""
	}

	var output string
	for _, comment := range c.CommentsBefore {
		// TODO: Preserve indentation
		output += strings.TrimSpace(strings.TrimPrefix(comment.Literal, "---")) + "  \n"
	}
	return output
}

type Block struct {
	Comments
	Stmts    []Statement
	StartPos token.Pos `json:"-"`
}

func (b *Block) Pos() token.Pos {
	return b.StartPos
}
func (b *Block) End() token.Pos {
	if len(b.Stmts) == 0 {
		return b.StartPos
	}
	return b.Stmts[len(b.Stmts)-1].End()
}

type TableField struct {
	Comments
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
	Comments
	Exps []Expression `json:",omitempty"`
	// OR
	Token *token.Token `json:",omitempty"`
}

func (i *Invalid) expressionNode() {}
func (i *Invalid) statementNode()  {}
func (i *Invalid) Pos() token.Pos {
	if i.Token != nil {
		return i.Token.Pos
	}
	return i.Exps[0].Pos()
}
func (i *Invalid) End() token.Pos {
	if i.Token != nil {
		return i.Token.End()
	}
	return i.Exps[len(i.Exps)-1].End()
}

// A Leaf node has no children and is interactable in the editor.
type LeafNode interface {
	leaf()
}
