// Parse converts a set of tokens into a Node on the AST.

package lua

import (
	"errors"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Parse(tokens []Token) *Block {
	root := &Block{}
	parseBlock(root, tokens)
	return root
}

func parseBlock(block *Block, tokens []Token) []Token {
	i := 0
	for len(tokens) > 0 {
		offset, err := parseStatement(block, tokens[i:])
		if err != nil {
			panic("Not yet handling errors")
		}
		i += offset
	}
	return tokens
}

func parseStatement(block *Block, tokens []Token) (int, error) {
	var toAdd Node = nil

	first := tokens[0]
	switch first.Type {
	case TokSymbol:
		if first.Raw == ";" {
			toAdd = &EmptyStatement{Range: first.Range}
		} else {
			return 0, errors.New("Unexpected symbol '" + first.Raw + "'")
		}
	case TokIdentifier:
		// TODO: varlist '=' explist
		// TODO: functioncall
	case TokLabel:
		toAdd = &Label{Range: first.Range, Raw: first.Raw}
	case TokKeyword:
		switch first.Raw {
		case "break":
			toAdd = &BreakStatement{Range: first.Range}
		case "goto":
			second := tokens[1]
			if second.Type == TokIdentifier {
				toAdd = &GotoStatement{
					Label: Identifier{second.Range, second.Raw},
					Range: expandRange(first.Range, second.Range),
				}
			} else {
				return 0, errors.New("Unexpected token '" + second.Raw + "'")
			}
		case "do":
			innerBlock := &Block{}
			i := 1
			for i < len(tokens) {
				offset, err := parseStatement(innerBlock, tokens[i:])
				if err != nil {
					return 0, err
				}
				i += offset
				if len(tokens) > i+1 && checkToken(tokens[i+1], TokKeyword, "end") {
					i++
					break
				}
			}
			toAdd = &DoStatement{Block: innerBlock, Range: expandRange(first.Range, tokens[i].Range)}
		case "while":
			// TODO: whilestatement
		case "repeat":
			// TODO: repeatstatement
		case "if":
			// TODO: ifstatement
		case "for":
			// TODO: forstatement
			// TODO: forinstatement
		case "function":
			// TODO: functiondeclarationstatement
		case "local":
			// TODO: localfunctiondeclarationstatement
			// TODO: localassignmentstatement
		}
	}

	if toAdd != nil {
		block.Range.End = toAdd.GetRange().End
		block.Stmts = append(block.Stmts, toAdd)
	}
	return i, nil
}

func checkToken(token Token, tokenType TokenType, raw string) bool {
	return token.Type == tokenType && token.Raw == raw
}
func expandRange(base protocol.Range, extend protocol.Range) protocol.Range {
	return protocol.Range{
		Start: protocol.Position{
			Line:      base.Start.Line,
			Character: base.Start.Character,
		},
		End: protocol.Position{
			Line:      extend.End.Line,
			Character: extend.End.Character,
		},
	}
}
