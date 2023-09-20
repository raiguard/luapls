package lexer

import (
	"github.com/raiguard/luapls/lua/token"
)

// TODO: Support unicode
type Lexer struct {
	input string
	pos   int
	char  byte

	lineBreaks []int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, pos: -1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	pos := l.pos
	tok := token.EOF

	switch l.char {
	case 0:
	case '=':
		l.readChar()
		if l.char == '=' {
			tok = token.EQUAL
			l.readChar()
		} else {
			tok = token.ASSIGN
		}
	case '^':
		tok = token.POW
		l.readChar()
	case '>':
		l.readChar()
		if l.char == '=' {
			tok = token.GEQ
			l.readChar()
		} else {
			tok = token.GT
		}
	case '<':
		l.readChar()
		if l.char == '=' {
			tok = token.LEQ
			l.readChar()
		} else {
			tok = token.LT
		}
	case '#':
		l.readChar()
		tok = token.HASH
	case '-':
		l.readChar()
		if l.char == '-' {
			tok = token.COMMENT
			l.readChar()
			l.readComment()
		} else {
			tok = token.MINUS
		}
	case '%':
		l.readChar()
		tok = token.PERCENT
	case '+':
		l.readChar()
		tok = token.PLUS
	case '/':
		l.readChar()
		tok = token.SLASH
	case '*':
		l.readChar()
		tok = token.STAR
	case '~':
		l.readChar()
		if l.char == '=' {
			l.readChar()
			tok = token.NEQ
		}
	case '(':
		l.readChar()
		tok = token.LPAREN
	case ')':
		l.readChar()
		tok = token.RPAREN
	case '[':
		if l.readRawString() {
			tok = token.RAWSTRING
		} else {
			tok = token.LBRACK
		}
	case ']':
		l.readChar()
		tok = token.RBRACK
	case '{':
		l.readChar()
		tok = token.LBRACE
	case '}':
		l.readChar()
		tok = token.RBRACE
	case ':':
		l.readChar()
		if l.char == ':' {
			l.readChar()
			tok = token.LABEL
		} else {
			tok = token.COLON
		}
	case ',':
		l.readChar()
		tok = token.COMMA
	case '.':
		l.readChar()
		if l.char == '.' {
			l.readChar()
			if l.char == '.' {
				l.readChar()
				tok = token.VARARG
			} else {
				tok = token.CONCAT
			}
		} else {
			tok = token.DOT
		}
	case ';':
		l.readChar()
		tok = token.SEMICOLON
	case '\'', '"':
		if l.readString() {
			tok = token.STRING
		}
	default:
		if l.readNumber() {
			tok = token.NUMBER
		} else if l.readIdentifier() {
			lit := l.input[pos:l.pos]
			if keyword, ok := token.Reserved[lit]; ok {
				tok = keyword
			} else {
				tok = token.IDENT
			}
		} else {
			l.readChar()
		}
	}

	return token.Token{
		Type:    tok,
		Literal: l.input[pos:l.pos],
		Pos:     pos,
	}
}

func (l *Lexer) readChar() {
	l.pos += 1
	if l.pos >= len(l.input) {
		l.char = 0
		return
	}
	l.char = l.input[l.pos]
	if l.char == '\n' {
		l.lineBreaks = append(l.lineBreaks, l.pos)
	}
}

// TODO: Type annotations
func (l *Lexer) readComment() {
	if l.char == '[' && l.readRawString() {
		return
	}
	for l.char != 0 && l.char != '\n' {
		l.readChar()
	}
	return
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() bool {
	if !isDigit(l.char) {
		return false
	}

	isZero := l.char == '0'

	inExponent := false
	inDecimal := false
	hexNum := false

	l.readChar()
	if l.char == 'x' || l.char == 'X' {
		if !isZero {
			return false
		}
		hexNum = true
		l.readChar()
	}

	for {
		if l.char == '.' {
			if inDecimal {
				return false
			}
			inDecimal = true
			l.readChar()
			continue
		}

		if isExponentLiteral(l.char, hexNum) {
			if inExponent {
				return false
			}
			inExponent = true
			l.readChar()
			if l.char == '+' || l.char == '-' {
				l.readChar()
			}
			continue
		}

		if !isNumberLiteral(l.char, hexNum) {
			break
		}
		l.readChar()
	}

	return true
}

func (l *Lexer) readIdentifier() bool {
	if !isIdentifier(l.char) {
		return false
	}
	for isIdentifier(l.char) {
		l.readChar()
	}
	return true
}

// TODO: Other escape sequences
func (l *Lexer) readString() bool {
	quote := l.char
	for {
		l.readChar()
		if l.char == '\n' {
			return false
		}
		if l.char == '\\' {
			l.readChar()
			if l.char == quote {
				continue
			}
		}
		if l.char == quote {
			break
		}
	}
	l.readChar()
	return true
}

func (l *Lexer) readRawString() bool {
	if l.char != '[' {
		return false
	}
	level := 0
	l.readChar()
	for l.char == '=' {
		level += 1
		l.readChar()
	}
	if l.char != '[' {
		return false
	}
	for {
		for l.char != ']' {
			if l.char == 0 {
				return false
			}
			l.readChar()
		}
		thisLevel := 0
		l.readChar()
		for l.char == '=' {
			thisLevel += 1
			l.readChar()
		}
		if l.char != ']' {
			continue
		}
		if thisLevel == level {
			break
		}
	}
	l.readChar()
	return true
}

func isDigit(lit byte) bool {
	return lit >= '0' && lit <= '9'
}

func isHex(lit byte) bool {
	return lit >= 'a' && lit <= 'f' || lit >= '0' && lit <= '9' || lit >= 'A' && lit <= 'F'
}

func isIdentifier(lit byte) bool {
	return (lit >= 'a' && lit <= 'z') || (lit >= 'A' && lit <= 'Z') || lit == '_' || isDigit(lit)
}

func isWhitespace(lit byte) bool {
	return lit == '\n' || lit == '\r' || lit == '\t' || lit == ' '
}

func isExponentLiteral(lit byte, hexNum bool) bool {
	if hexNum {
		return lit == 'p' || lit == 'P'
	} else {
		return lit == 'e' || lit == 'E'
	}
}

func isNumberLiteral(lit byte, hexNum bool) bool {
	if hexNum {
		return isHex(lit)
	} else {
		return isDigit(lit)
	}
}
