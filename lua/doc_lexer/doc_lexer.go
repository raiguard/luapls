package doc_lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/raiguard/luapls/lua/token"
)

type Lexer struct {
	input string // the string being scanned.
	start int    // start position of this token.
	pos   int    // current position in the input.
	width int    // width of last rune read.
}

func Run(input string) []token.Token {
	l := Lexer{input: input, pos: 0}
	tokens := []token.Token{}
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) Next() token.Token {
	l.ignore()

	if l.acceptRun(" \t") {
		return l.makeToken(token.WHITESPACE)
	}

	tok := token.INVALID

	switch r := l.read(); r {
	case 0, '\r', '\n':
		tok = token.EOF
	case '@':
		l.read()
		l.readIdentifier()
		if reserved, ok := token.Reserved[l.input[l.start:l.pos]]; ok {
			tok = reserved
		}
	default:
		l.readIdentifier()
		tok = token.IDENT
	}

	return l.makeToken(tok)
}

func (l *Lexer) makeToken(typ token.TokenType) token.Token {
	return token.Token{
		Type:    typ,
		Literal: l.input[l.start:l.pos],
		Pos:     l.start,
	}
}

// read returns the next rune in the input.
func (l *Lexer) read() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	// if r == '\n' && (len(l.lineBreaks) == 0 || l.lineBreaks[len(l.lineBreaks)-1] != l.pos) {
	// 	l.lineBreaks = append(l.lineBreaks, l.pos)
	// }
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	rune := l.read()
	l.backup()
	return rune
}

// accept consumes the next rune if it is from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.read()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun accepts a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) bool {
	success := false
	for l.accept(valid) {
		success = true
	}
	return success
}

// acceptNot consumes the next rune if it is from the valid set.
func (l *Lexer) acceptNot(valid string) bool {
	r := l.read()
	if r != 0 && strings.IndexRune(valid, r) == -1 {
		return true
	}
	l.backup()
	return false
}

// acceptNotRun acceptNots a run of runes from the valid set.
func (l *Lexer) acceptNotRun(valid string) {
	for l.acceptNot(valid) {
	}
}

func (l *Lexer) readIdentifier() bool {
	for isIdentifier(l.read()) {
	}
	l.backup()
	return true
}

func isIdentifier(rune rune) bool {
	return unicode.IsLetter(rune) || unicode.IsDigit(rune) || rune == '_'
}
