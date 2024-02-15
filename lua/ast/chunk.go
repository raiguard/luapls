package ast

import (
	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Chunk struct {
	Block      Block
	Errors     []Error
	LineBreaks []int
	// TODO: Global exports, etc.
}

func (f *Chunk) ToPos(position protocol.Position) token.Pos {
	line := int(position.Line)
	col := int(position.Character)
	if line >= len(f.LineBreaks) {
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

func (f *Chunk) ToProtocolPos(pos token.Pos) protocol.Position {
	if len(f.LineBreaks) == 0 {
		return protocol.Position{
			Line:      0,
			Character: uint32(pos),
		}
	}
	line := 0
	lineStart := 0
	lineEnd := -1
	for line < len(f.LineBreaks) {
		lineStart = lineEnd + 1
		lineEnd = f.LineBreaks[line]
		if lineStart <= pos && lineEnd >= pos {
			break
		}
		line++
	}
	return protocol.Position{
		Line:      uint32(line),
		Character: uint32(pos - lineStart),
	}
}

func (f *Chunk) ToProtocolRange(rng token.Range) protocol.Range {
	return protocol.Range{
		Start: f.ToProtocolPos(rng[0]),
		End:   f.ToProtocolPos(rng[1]),
	}
}
