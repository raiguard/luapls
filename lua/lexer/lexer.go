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
	l.nextChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	l.savePosition()

	switch l.char {
	case '=':
		if l.expectPeek('=') {
			tok = l.newToken(token.EQUAL)
		} else {
			tok = l.newToken(token.ASSIGN)
		}
	case '^':
		tok = l.newToken(token.CARET)
	case '>':
		if l.expectPeek('=') {
			tok = l.newToken(token.GEQ)
		} else {
			tok = l.newToken(token.GT)
		}
	case '#':
		tok = l.newToken(token.HASH)
	case '<':
		if l.expectPeek('=') {
			tok = l.newToken(token.LEQ)
		} else {
			tok = l.newToken(token.LT)
		}
	case '-':
		if l.expectPeek('-') {
			l.skipComment()
		} else {
			tok = l.newToken(token.MINUS)
		}
	case '%':
		tok = l.newToken(token.PERCENT)
	case '+':
		tok = l.newToken(token.PLUS)
	case '/':
		tok = l.newToken(token.SLASH)
	case '*':
		tok = l.newToken(token.STAR)
	case '~':
		if l.expectPeek('=') {
			tok = l.newToken(token.NEQ)
		}
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '[':
		// TODO: Raw string
		tok = l.newToken(token.LBRACK)
	case ']':
		tok = l.newToken(token.RBRACK)
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case ':':
		// TODO: Label
		tok = l.newToken(token.COLON)
	case ',':
		tok = l.newToken(token.COMMA)
	case '.':
		if l.expectPeek('.') {
			if l.expectPeek('.') {
				tok = l.newToken(token.SPREAD)
			} else {
				tok = l.newToken(token.CONCAT)
			}
		} else {
			tok = l.newToken(token.DOT)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON)
	default:
		if isDigit(l.char) {
			if l.readNumber() {
				tok = l.newToken(token.NUMBER)
			}
			break
		}

		lit := l.readIdentifier()
		if keyword, ok := token.Keywords[lit]; ok {
			tok = l.newToken(keyword)
		} else {
			tok = l.newToken(token.IDENT)
		}
	}

	l.nextChar()

	return tok
}

func (l *Lexer) nextChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.char == '\n' {
		l.col = 0
		l.row++
	} else {
		l.col++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) expectPeek(expected byte) bool {
	if l.peekChar() == expected {
		l.nextChar()
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
			EndCol:   l.col + 1,
			EndRow:   l.row,
		},
	}
}

func (l *Lexer) curLiteral() string {
	if l.char == 0 {
		return l.input[l.savedPos:l.position]
	} else {
		return l.input[l.savedPos:l.readPosition]

	}
}

func (l *Lexer) skipComment() {
	// TODO: Raw string comments
	for l.char != '\n' {
		l.nextChar()
	}
	l.nextChar()
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.nextChar()
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
		l.nextChar()
		if !isHex(l.char) {
			return false
		}
	}

	for {
		l.nextChar()
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
				l.nextChar()
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
	for isIdentifier(l.peekChar()) {
		l.nextChar()
	}

	return l.curLiteral()
}

func isDigit(lit byte) bool {
	return lit >= '0' && lit <= '9'
}

func isHex(lit byte) bool {
	return lit >= '0' && lit <= '9' || lit >= 'A' && lit <= 'F'
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
