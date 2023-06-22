package lexer

import "luapls/lua/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
	col, row     int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, col: -1}
	l.nextChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			tok = l.newToken(token.EQUAL, "==")
			l.nextChar()
		} else {
			tok = l.newToken(token.ASSIGN, "=")
		}
	case '^':
		tok = l.newToken(token.CARET, "^")
	case '>':
		if l.peekChar() == '=' {
			tok = l.newToken(token.GEQ, ">=")
			l.nextChar()
		} else {
			tok = l.newToken(token.GT, ">")
		}
	case '#':
		tok = l.newToken(token.HASH, "#")
	case '<':
		if l.peekChar() == '=' {
			tok = l.newToken(token.LEQ, "<=")
			l.nextChar()
		} else {
			tok = l.newToken(token.LT, "<")
		}
	case '-':
		tok = l.newToken(token.MINUS, "-")
	case '%':
		tok = l.newToken(token.PERCENT, "%")
	case '+':
		tok = l.newToken(token.PLUS, "+")
	case '/':
		tok = l.newToken(token.SLASH, "/")
	case '*':
		tok = l.newToken(token.STAR, "*")
	case '(':
		tok = l.newToken(token.LPAREN, "(")
	case ')':
		tok = l.newToken(token.RPAREN, ")")
	case '[':
		// TODO: Raw string
		tok = l.newToken(token.LBRACK, "[")
	case ']':
		tok = l.newToken(token.RBRACK, "]")
	case '{':
		tok = l.newToken(token.LBRACE, "{")
	case '}':
		tok = l.newToken(token.RBRACE, "}")
	case ':':
		// TODO: Label
		tok = l.newToken(token.COLON, ":")
	case ',':
		tok = l.newToken(token.COMMA, ",")
	case '.':
		tok = l.newToken(token.DOT, ".")
	case ';':
		tok = l.newToken(token.SEMICOLON, ";")
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

func (l *Lexer) newToken(tokType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokType,
		Literal: literal,
		Range: token.Range{
			StartCol: l.col,
			StartRow: l.row,
			EndCol:   l.col + len(literal),
			EndRow:   l.row, // TODO: Raw strings
		},
	}
}
