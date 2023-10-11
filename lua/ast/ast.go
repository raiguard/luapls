package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	String() string
	Pos() token.Pos
	End() token.Pos
}

type File struct {
	Block      Block
	LineBreaks []int
	// TODO: Global exports, etc.
}

type Block struct {
	Stmts    []Statement
	StartPos token.Pos
}

func (b *Block) String() string {
	var out string
	for _, stmt := range b.Stmts {
		out += stmt.String() + "\n"
	}
	return strings.TrimSpace(out)
}
func (b *Block) Pos() token.Pos {
	return b.StartPos
}
func (b *Block) End() token.Pos {
	if len(b.Stmts) > 0 {
		return b.StartPos + b.Stmts[len(b.Stmts)-1].End()
	}
	return b.StartPos
}

type TableField struct {
	Key      Expression
	Value    Expression
	StartPos token.Pos // Needed in case of bracketed keys
}

func (tf *TableField) String() string {
	if tf.Key == nil {
		return tf.Value.String()
	}
	if ident, ok := tf.Key.(*Identifier); ok {
		return fmt.Sprintf("%s = %s", ident.String(), tf.Value.String())
	}
	return fmt.Sprintf("[%s] = %s", tf.Key.String(), tf.Value.String())
}
func (rf *TableField) Pos() token.Pos {
	return rf.StartPos
}
func (rf *TableField) End() token.Pos {
	return rf.Value.End()
}

func nodeListToString[T Node](nodes []T) string {
	items := []string{}
	for _, node := range nodes {
		items = append(items, node.String())
	}
	return strings.Join(items, ", ")
}
