package lua

import (
	"strings"
)

// A representation of all possible tokens in the Lua language
type Token uint8

const (
	// Special
	INVALID Token = iota
	EOF
	COMMENT // Can be single-line or multi-line (raw string)
	SPACE   // Tabs or spaces
	IDENTIFIER
	// Keywords
	keyword_beg
	AND
	BREAK
	DO
	ELSE
	ELSEIF
	END
	FOR
	FUNCTION
	GOTO
	IF
	IN
	LOCAL
	NOT
	OR
	REPEAT
	RETURN
	THEN
	UNTIL
	WHILE
	keyword_end
	// Literals
	literal_beg
	FALSE
	LABEL // ::foo::
	NIL
	NUMBER
	RAWSTRING
	STRING // "foo" or 'foo'
	TRUE
	literal_end
	// Symbols
	symbol_beg
	ADD    // +
	SUB    // -
	MUL    // *
	DIV    // /
	MOD    // %
	POW    // ^
	LEN    // #
	EQL    // ==
	NEQ    // ~=
	LEQ    // <=
	GEQ    // >=
	LSS    // <
	GTR    // >
	ASSIGN // =
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	LBRACK // [
	RBRACK // ]
	// DCOLON // ::
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	CONCAT    // ..
	VARARG    // ...
	symbol_end
)

var tokens = map[string]Token{
	// Keywords
	"and":      AND,
	"break":    BREAK,
	"do":       DO,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"end":      END,
	"for":      FOR,
	"function": FUNCTION,
	"goto":     GOTO,
	"if":       IF,
	"in":       IN,
	"local":    LOCAL,
	"not":      NOT,
	"or":       OR,
	"repeat":   REPEAT,
	"return":   RETURN,
	"then":     THEN,
	"until":    UNTIL,
	"while":    WHILE,
	// Literals
	"false": FALSE,
	"nil":   NIL,
	"true":  TRUE,
	// Symbols
	"+":   ADD,
	"-":   SUB,
	"*":   MUL,
	"/":   DIV,
	"%":   MOD,
	"^":   POW,
	"#":   LEN,
	"==":  EQL,
	"~=":  NEQ,
	"<=":  LEQ,
	">=":  GEQ,
	"<":   LSS,
	">":   GTR,
	"=":   ASSIGN,
	"(":   LPAREN,
	")":   RPAREN,
	"{":   LBRACE,
	"}":   RBRACE,
	"[":   LBRACK,
	"]":   RBRACK,
	",":   COMMA,
	";":   SEMICOLON,
	":":   COLON,
	".":   DOT,
	"..":  CONCAT,
	"...": VARARG,
}

func isIdentifier(s string) bool {
	if s == "" || isKeyword(s) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !isLetter(c) || c != '_' || (i == 0 || !isDigit(c)) {
			return false
		}
	}
	return true
}

func isKeyword(s string) bool {
	tok, ok := tokens[s]
	return ok && keyword_beg < tok && tok < keyword_end
}

func isLabel(s string) bool {
	return strings.HasPrefix(s, "::") && strings.HasSuffix(s, "::")
}
