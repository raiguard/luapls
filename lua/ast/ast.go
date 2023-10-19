package ast

import (
	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

func Range(n Node) token.Range {
	return token.Range{n.Pos(), n.End()}
}

type Block struct {
	Stmts    []Statement
	StartPos token.Pos `json:"-"`
}

func (b *Block) Pos() token.Pos {
	return b.StartPos
}
func (b *Block) End() token.Pos {
	if len(b.Stmts) > 0 {
		return b.Stmts[len(b.Stmts)-1].End()
	}
	return b.StartPos
}

type TableField struct {
	Key      Expression
	Value    Expression
	StartPos token.Pos `json:"-"` // Needed in case of bracketed keys
}

func (rf *TableField) Pos() token.Pos {
	return rf.StartPos
}
func (rf *TableField) End() token.Pos {
	return rf.Value.End()
}

// A Leaf node has no children and is interactable in the editor.
type LeafNode interface {
	leaf()
}
