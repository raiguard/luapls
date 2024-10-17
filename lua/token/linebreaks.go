package token

import protocol "github.com/tliron/glsp/protocol_3_16"

type LineBreaks []int

func (f LineBreaks) ToPos(position protocol.Position) Pos {
	line := int(position.Line)
	col := int(position.Character)
	if line >= len(f) {
		return InvalidPos
	}
	lineStart := 0
	if line > 0 {
		lineStart = f[line-1] + 1
	}
	lineEnd := f[line]
	if col > lineEnd-lineStart {
		return InvalidPos
	}
	return lineStart + col
}

func (f LineBreaks) ToProtocolPos(pos Pos) protocol.Position {
	if len(f) == 0 {
		return protocol.Position{
			Line:      0,
			Character: uint32(pos),
		}
	}
	line := 0
	lineStart := 0
	lineEnd := -1
	for line < len(f) {
		lineStart = lineEnd + 1
		lineEnd = f[line]
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

func (f LineBreaks) ToProtocolRange(rng Range) protocol.Range {
	return protocol.Range{
		Start: f.ToProtocolPos(rng.Start),
		End:   f.ToProtocolPos(rng.End),
	}
}
