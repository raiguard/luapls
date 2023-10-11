package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
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

func (f *File) String() string {
	return f.Block.String()
}
func (f *File) Pos() token.Pos {
	return f.Block.Pos()
}
func (f *File) End() token.Pos {
	return f.Block.End()
}

func (f *File) ToPos(position protocol.Position) token.Pos {
	line := int(position.Line)
	col := int(position.Character)
	if line > len(f.LineBreaks) {
		return token.InvalidPos
	}
	lineStart := 0
	if line > 0 {
		lineStart = f.LineBreaks[line-1] + 1
	}
	lineEnd := f.LineBreaks[line]
	if col > lineEnd-lineStart {
		return token.InvalidPos
	}
	return lineStart + col
}

func (f *File) ToProtocolPos(pos token.Pos) protocol.Position {
	if len(f.LineBreaks) == 0 {
		return protocol.Position{
			Line:      0,
			Character: uint32(pos),
		}
	}
	line := 0
	lineStart := 0
	lineEnd := -1
	for i := 0; i < len(f.LineBreaks); i++ {
		lineStart = lineEnd + 1
		lineEnd = f.LineBreaks[i]
		if lineStart <= pos && lineEnd >= pos {
			line = i
			break
		}
	}
	return protocol.Position{
		Line:      uint32(line),
		Character: uint32(pos - lineStart),
	}
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
