package lexer

import (
	"luapls/lua/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
	col, row     int

	savedCol, savedRow, savedPos int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, col: -1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

tryAgain:
	l.skipWhitespace()
	l.savePosition()

	switch l.char {
	case 0:
		return token.Token{Type: token.EOF}
	case '=':
		if l.expectPeek('=') {
			tok = l.readNewToken(token.EQUAL)
		} else {
			tok = l.readNewToken(token.ASSIGN)
		}
	case '^':
		tok = l.readNewToken(token.CARET)
	case '>':
		if l.expectPeek('=') {
			tok = l.readNewToken(token.GEQ)
		} else {
			tok = l.readNewToken(token.GT)
		}
	case '#':
		tok = l.readNewToken(token.HASH)
	case '<':
		if l.expectPeek('=') {
			tok = l.readNewToken(token.LEQ)
		} else {
			tok = l.readNewToken(token.LT)
		}
	case '-':
		if l.expectPeek('-') {
			if l.skipComment() {
				// FIXME: This sucks
				goto tryAgain
			}
		} else {
			tok = l.readNewToken(token.MINUS)
		}
	case '%':
		tok = l.readNewToken(token.PERCENT)
	case '+':
		tok = l.readNewToken(token.PLUS)
	case '/':
		tok = l.readNewToken(token.SLASH)
	case '*':
		tok = l.readNewToken(token.STAR)
	case '~':
		if l.expectPeek('=') {
			tok = l.readNewToken(token.NEQ)
		}
	case '(':
		tok = l.readNewToken(token.LPAREN)
	case ')':
		tok = l.readNewToken(token.RPAREN)
	case '[':
		if l.peekChar() == '[' || l.peekChar() == '=' {
			if l.readRawString() {
				tok = l.readNewToken(token.RAWSTRING)
			}
		} else {
			tok = l.readNewToken(token.LBRACK)
		}
	case ']':
		tok = l.readNewToken(token.RBRACK)
	case '{':
		tok = l.readNewToken(token.LBRACE)
	case '}':
		tok = l.readNewToken(token.RBRACE)
	case ':':
		if l.expectPeek(':') {
			tok = l.readNewToken(token.LABEL)
		} else {
			tok = l.readNewToken(token.COLON)
		}
	case ',':
		tok = l.readNewToken(token.COMMA)
	case '.':
		if l.expectPeek('.') {
			if l.expectPeek('.') {
				tok = l.readNewToken(token.SPREAD)
			} else {
				tok = l.readNewToken(token.CONCAT)
			}
		} else {
			tok = l.readNewToken(token.DOT)
		}
	case ';':
		tok = l.readNewToken(token.SEMICOLON)
	case '\'', '"':
		if l.readString(l.char) {
			tok = l.readNewToken(token.STRING)
		}
	default:
		if isDigit(l.char) {
			if l.readNumber() {
				tok = l.newToken(token.NUMBER)
			}
			tok.Literal = l.curLiteral()
			break
		} else if isIdentifier(l.char) {
			lit := l.readIdentifier()
			if keyword, ok := token.Keywords[lit]; ok {
				tok = l.newToken(keyword)
			} else {
				tok = l.newToken(token.IDENT)
			}
		} else {
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readChar() {
	if l.char == '\n' {
		l.col = 0
		l.row++
	} else {
		l.col++
	}
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) expectPeek(expected byte) bool {
	if l.peekChar() == expected {
		l.readChar()
		return true
	}
	return false
}

func (l *Lexer) savePosition() {
	l.savedCol = l.col
	l.savedRow = l.row
	l.savedPos = l.position
}

func (l *Lexer) newToken(tokType token.TokenType) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: l.curLiteral(),
		Range: token.Range{
			StartCol: l.savedCol,
			StartRow: l.savedRow,
			EndCol:   l.col,
			EndRow:   l.row,
		},
	}
}

func (l *Lexer) readNewToken(tokType token.TokenType) token.Token {
	l.readChar()
	return l.newToken(tokType)
}

func (l *Lexer) curLiteral() string {
	return l.input[l.savedPos:l.position]
}

func (l *Lexer) skipComment() bool {
	if l.expectPeek('[') {
		return l.readRawString()
	}
	for l.char != '\n' {
		l.readChar()
	}
	return true
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() bool {
	inExponent := false
	inDecimal := false
	hexNum := false

	if l.peekChar() == 'x' || l.peekChar() == 'X' {
		if l.char != '0' {
			return false
		}
		hexNum = true
		l.readChar()
		if !isHex(l.peekChar()) {
			return false
		}
	}

	for {
		l.readChar()
		if l.char == '.' {
			if inDecimal {
				return false
			}
			inDecimal = true
			continue
		}

		if isExponentLiteral(l.char, hexNum) {
			if inExponent {
				return false
			}
			inExponent = true
			if l.peekChar() == '+' || l.peekChar() == '-' {
				l.readChar()
			}
			continue
		}

		if !isNumberLiteral(l.char, hexNum) {
			break
		}
	}

	return true
}

func (l *Lexer) readIdentifier() string {
	for isIdentifier(l.char) {
		l.readChar()
	}

	return l.curLiteral()
}

func (l *Lexer) readString(quote byte) bool {
	l.readChar()
	for l.char != '\n' && l.char != quote {
		l.readChar()
	}
	return l.char == quote
}

func (l *Lexer) readRawString() bool {
	level := 0
	for l.expectPeek('=') {
		level += 1
	}
	l.readChar()
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
		for l.expectPeek('=') {
			thisLevel += 1
		}
		if thisLevel == level && l.expectPeek(']') {
			break
		}
		l.readChar()
	}

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
