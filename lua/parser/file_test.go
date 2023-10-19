package parser

// import (
// 	"testing"

// 	"github.com/raiguard/luapls/lua/token"
// 	"github.com/stretchr/testify/assert"
// 	protocol "github.com/tliron/glsp/protocol_3_16"
// )

// func TestFilePositions(t *testing.T) {
// 	input := `i = 12345
// local function bar(baz)
//   print(baz)
// end
// bar(i)
// `
// 	p := New(input)
// 	file := p.ParseFile()
// 	checkParserErrors(t, p)
// 	assertSlicesEqual(t, []int{9, 33, 46, 50, 57}, file.LineBreaks)

// 	assertPositionsMatch(t, &file, protocol.Position{Line: 0, Character: 2}, 2)
// 	assertPositionsMatch(t, &file, protocol.Position{Line: 2, Character: 7}, 41)

// 	assertPositionsMatch(t, &file, protocol.Position{Line: 1, Character: 0}, file.Block.Stmts[1].Pos())
// }

// func assertSlicesEqual[T comparable](t *testing.T, expected []T, actual []T) {
// 	assert.Equal(t, len(expected), len(actual))

// 	for i, expected := range expected {
// 		assert.Equal(t, expected, actual[i])
// 	}
// }

// func assertPositionsMatch(t *testing.T, file *File, position protocol.Position, pos token.Pos) {
// 	protocolPos := file.ToProtocolPos(pos)
// 	assert.Equal(t, protocolPos.Line, position.Line)
// 	assert.Equal(t, protocolPos.Character, position.Character)

// 	astPos := file.ToPos(position)
// 	assert.Equal(t, pos, astPos)
// }
