package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/raiguard/luapls/lua/token"
)

const digits string = "0123456789"
const hexDigits string = digits + "abcdefABCDEF"
const whitespace string = " \t\r\n"

type Lexer struct {
	input string // the string being scanned.
	start int    // start position of this token.
	pos   int    // current position in the input.
	width int    // width of last rune read.

	lineBreaks []int
}

func New(input string) *Lexer {
	return &Lexer{input: input, pos: 0, lineBreaks: []int{}}
}

func (l *Lexer) Next() token.Token {
	l.ignore()

	if l.acceptRun(whitespace) {
		return l.makeToken(token.WHITESPACE)
	}

	tok := token.INVALID

	switch r := l.read(); r {
	case 0: // EOF
		tok = token.EOF
	case '=':
		if l.accept("=") {
			tok = token.EQUAL
		} else {
			tok = token.ASSIGN
		}
	case '^':
		tok = token.POW
	case '>':
		if l.accept("=") {
			tok = token.GEQ
		} else {
			tok = token.GT
		}
	case '<':
		if l.accept("=") {
			tok = token.LEQ
		} else {
			tok = token.LT
		}
	case '#':
		tok = token.LEN
	case '-':
		if l.accept("-") {
			tok = token.COMMENT
			l.readComment()
		} else {
			tok = token.MINUS
		}
	case '%':
		tok = token.MOD
	case '+':
		tok = token.PLUS
	case '/':
		tok = token.SLASH
	case '*':
		tok = token.MUL
	case '~':
		if l.accept("=") {
			tok = token.NEQ
		}
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case '[':
		if l.readRawString() {
			tok = token.RAWSTRING
		} else {
			tok = token.LBRACK
		}
	case ']':
		tok = token.RBRACK
	case '{':
		tok = token.LBRACE
	case '}':
		tok = token.RBRACE
	case ':':
		if l.accept(":") {
			tok = token.LABEL
		} else {
			tok = token.COLON
		}
	case ',':
		tok = token.COMMA
	case '.':
		if l.accept(".") {
			if l.accept(".") {
				tok = token.VARARG
			} else {
				tok = token.CONCAT
			}
		} else if unicode.IsDigit(l.peek()) {
			l.backup()
			l.readNumber()
			tok = token.NUMBER
		} else {
			tok = token.DOT
		}
	case ';':
		tok = token.SEMICOLON
	case '\'', '"':
		if l.readString(r) {
			tok = token.STRING
		}
	default:
		if unicode.IsDigit(r) {
			l.backup()
			if l.readNumber() {
				tok = token.NUMBER
			} else {
				l.ignore()
			}
		} else if l.readIdentifier() {
			if reserved, ok := token.Reserved[l.input[l.start:l.pos]]; ok {
				tok = reserved
			} else {
				tok = token.IDENT
			}
		} else {
			l.read() // Always make progress
		}
	}

	return l.makeToken(tok)
}

func (l *Lexer) GetLineBreaks() []int {
	return l.lineBreaks
}

func (l *Lexer) Run() ([]token.Token, []int) {
	tokens := []token.Token{}
	var tok token.Token
	for tok.Type != token.EOF {
		tok = l.Next()
		tokens = append(tokens, tok)
	}
	return tokens, l.lineBreaks
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
	if r == '\n' && (len(l.lineBreaks) == 0 || l.lineBreaks[len(l.lineBreaks)-1] != l.pos) {
		l.lineBreaks = append(l.lineBreaks, l.pos)
	}
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

// TODO: Type annotations
func (l *Lexer) readComment() {
	if l.accept("[") && l.readRawString() {
		return
	}
	for l.acceptNot("\n") {
	}
	return
}

func (l *Lexer) readNumber() bool {
	useDigits := digits
	if l.accept("0") && l.accept("xX") {
		useDigits = hexDigits
	}
	l.acceptRun(useDigits)

	if l.accept(".") {
		l.acceptRun(useDigits)
	}

	if l.accept("eE") {
		l.acceptRun(digits)
		l.accept("+-")
		l.acceptRun(digits)
	} else if l.accept("pP") {
		l.acceptRun(hexDigits)
		l.accept("+-")
		l.acceptRun(hexDigits)
	}

	return true
}

func (l *Lexer) readIdentifier() bool {
	for isIdentifier(l.read()) {
	}
	l.backup()
	return true
}

func (l *Lexer) readString(quote rune) bool {
	for {
		if l.accept(string(quote)) {
			break
		}
		if l.accept("\n") {
			// TODO: Lexing error
			return false
		}
		if l.accept("\\") {
			if l.accept("z") {
				l.acceptRun(whitespace)
				continue
			}
			if l.accept(string(quote)) || l.accept("abfnrtvz\r\n") {
				continue
			}
		}
		l.read()
	}
	return true
}

func (l *Lexer) readRawString() bool {
	if !l.accept("[=") {
		return false
	}
	l.backup()
	level := 0
	for l.accept("=") {
		level++
	}
	if !l.accept("[") {
		return false
	}
	for {
		l.acceptNotRun("]")
		if !l.accept("]") {
			return false
		}
		thisLevel := 0
		for l.accept("=") {
			thisLevel++
		}
		if !l.accept("]") {
			continue
		}
		if thisLevel == level {
			break
		}
		l.backup()
	}
	return true
}

func isIdentifier(rune rune) bool {
	return unicode.IsLetter(rune) || unicode.IsDigit(rune) || rune == '_'
}
