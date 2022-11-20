// Lex converts Lua source code into a series of tokens.

package lua

import (
	"fmt"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type TokenType int

const (
	TokComment TokenType = iota
	TokIdentifier
	TokLabel
	TokKeyword
	TokNumber
	TokRawString
	TokString
	TokSymbol
	TokInvalid TokenType = -1
)

type Token struct {
	Raw   string
	Range protocol.Range
	Type  TokenType
}

// Converts a [[Token]] into a user-friendly string.
func StrToken(token *Token) string {
	typeStr := ""
	switch token.Type {
	case TokComment:
		typeStr = "Comment"
	case TokIdentifier:
		typeStr = "Identifier"
	case TokLabel:
		typeStr = "Label"
	case TokKeyword:
		typeStr = "Keyword"
	case TokNumber:
		typeStr = "Number"
	case TokRawString:
		typeStr = "RawString"
	case TokString:
		typeStr = "String"
	case TokSymbol:
		typeStr = "Symbol"
	case TokInvalid:
		typeStr = "Invalid"
	}

	return fmt.Sprintf(
		"%s %s | %d.%d:%d.%d",
		typeStr,
		token.Raw,
		token.Range.Start.Line,
		token.Range.Start.Character,
		token.Range.End.Line,
		token.Range.End.Character,
	)
}

var keywords = map[string]bool{
	"and": true, "break": true, "do": true, "else": true, "elseif": true, "end": true,
	"false": true, "for": true, "function": true, "goto": true, "if": true, "in": true,
	"local": true, "nil": true, "not": true, "or": true, "repeat": true, "return": true,
	"then": true, "true": true, "until": true, "while": true,
}

var symbols = []string{
	"+", "-", "*", "/", "%", "^", "#",
	"==", "~=", "<=", ">=", "<", ">", "=",
	"(", ")", "{", "}", "[", "]", "::",
	";", ":", ",", ".", "..", "...",
}

func Tokenize(lines []string) []Token {
	var out []Token
	for i, line := range lines {
		out = TokenizeLine(uint32(i), line, out)
	}
	return out
}

func TokenizeLine(line uint32, s string, out []Token) []Token {
	// TODO: Raw string
	var pos uint32 = 0
	for {
		for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
			pos++
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}
		tok, n := next(s)
		end := pos + uint32(n)
		out = append(
			out,
			Token{
				s[:n],
				protocol.Range{
					Start: protocol.Position{Line: line, Character: pos},
					End:   protocol.Position{Line: line, Character: end},
				},
				tok,
			},
		)
		pos = end
		s = s[n:]
	}

	return out
}

func next(s string) (TokenType, int) {
	if isDigit(s[0]) {
		i := 1
		hasDecimal := false
		hasExponent := false
		isHex := false
		if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
			isHex = true
			i++
		}
		for i < len(s) {
			if isDigit(s[i]) || (isHex && isHexAlpha(s[i])) {
				i++
			} else if !hasDecimal && s[i] == '.' {
				hasDecimal = true
				i++
			} else if !hasExponent && isExponent(s[i], isHex) {
				hasExponent = true
				i++
				if len(s) > i && (s[i] == '-' || s[i] == '+') {
					i++
				}
			} else {
				break
			}
		}
		return TokNumber, i
	}
	// TODO: Other kinds of numbers
	if strings.HasPrefix(s, "--") {
		// TODO: multiline comments / raw strings
		return TokComment, len(s)
	}
	if strings.HasPrefix(s, "::") {
		i := 2
		for i < len(s) && isIdentifier(s[i]) {
			i++
		}
		if s[i:i+2] == "::" {
			return TokLabel, i + 2
		}
		return TokInvalid, i
	}
	if s[0] == '"' {
		i := 1
		for {
			if i >= len(s) {
				return TokInvalid, len(s)
			}
			if s[i] == '"' {
				return TokString, i + 1
			}
			// TODO: Escapes
			i++
		}
	}
	if s[0] == '\'' {
		i := 1
		for {
			if i >= len(s) {
				return TokInvalid, len(s)
			}
			if s[i] == '\'' {
				return TokString, i + 1
			}
			// TODO: Escapes
			i++
		}
	}
	if isIdentifier(s[0]) {
		i := 1
		// Identifiers may have digits as long as the first is not a digit
		for i < len(s) && isAlnum(s[i]) {
			i++
		}
		if keywords[s[:i]] {
			return TokKeyword, i
		}
		return TokIdentifier, i
	}
	for _, sym := range symbols {
		if strings.HasPrefix(s, sym) {
			return TokSymbol, len(sym)
		}
	}

	return TokInvalid, len(s)
}

func isExponent(c byte, isHex bool) bool {
	if isHex {
		return c == 'p' || c == 'P'
	} else {
		return c == 'e' || c == 'E'
	}
}
func isHexAlpha(c byte) bool {
	return (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
func isIdentifier(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
func isAlnum(c byte) bool {
	return isIdentifier(c) || isDigit(c)
}
