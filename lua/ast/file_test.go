package ast

// func TestFilePositions(t *testing.T) {
// 	input := `i = 12345
// local function bar(baz)
//   print(baz)
// end
// bar(i)
// `
// 	p := parser.New(input)
// 	file := p.ParseFile()
// 	assert.ElementsMatch(t, []int{9, 33, 46, 50, 57}, file.LineBreaks)

// 	assertPositionsMatch(t, &file, protocol.Position{Line: 0, Character: 2}, 2)
// 	assertPositionsMatch(t, &file, protocol.Position{Line: 2, Character: 7}, 41)

// 	assertPositionsMatch(t, &file, protocol.Position{Line: 1, Character: 0}, file.Block.Pairs[1].Pos())
// }

// func assertPositionsMatch(t *testing.T, file *File, position protocol.Position, pos token.Pos) {
// 	protocolPos := file.ToProtocolPos(pos)
// 	assert.Equal(t, protocolPos.Line, position.Line)
// 	assert.Equal(t, protocolPos.Character, position.Character)

// 	astPos := file.ToPos(position)
// 	assert.Equal(t, pos, astPos)
// }
