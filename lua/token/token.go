package token

import (
	"encoding/json"
	"fmt"
)

type TokenType int

const (
	INVALID TokenType = iota
	EOF
	COMMENT

	// Variable
	IDENT
	LABEL

	// Keywords
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
	REPEAT
	RETURN
	THEN
	UNTIL
	WHILE

	// Literals
	FALSE
	NIL
	NUMBER
	RAWSTRING
	STRING
	TRUE

	// Operators
	AND
	ASSIGN
	POW
	CONCAT
	EQUAL
	GEQ
	GT
	HASH
	LEQ
	LT
	MINUS
	NEQ
	NOT
	OR
	PERCENT
	PLUS
	SLASH
	STAR

	// Structure
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	LBRACE
	RBRACE

	// Grammar
	COLON
	COMMA
	DOT
	SEMICOLON
	VARARG
)

func (t TokenType) String() string {
	return TokenStr[t]
}

func (t *TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Pos
}

type Pos = int

func (t Token) End() int {
	return t.Pos + len(t.Literal)
}

func (t Token) String() string {
	return fmt.Sprintf("[%s] %s %v", TokenStr[t.Type], t.Literal, t.Pos)
}

var TokenStr = map[TokenType]string{
	INVALID: "invalid",
	EOF:     "eof",
	IDENT:   "identifier",
	LABEL:   "label",
	COMMENT: "comment",

	// Keywords
	BREAK:    "break",
	DO:       "do",
	ELSE:     "else",
	ELSEIF:   "elseif",
	END:      "end",
	FOR:      "for",
	FUNCTION: "function",
	GOTO:     "goto",
	IF:       "if",
	IN:       "in",
	LOCAL:    "local",
	REPEAT:   "repeat",
	RETURN:   "return",
	THEN:     "then",
	UNTIL:    "until",
	WHILE:    "while",

	// Literals
	FALSE:     "false",
	NIL:       "nil",
	NUMBER:    "number",
	RAWSTRING: "rawstring",
	STRING:    "string",
	TRUE:      "true",

	// Operators
	AND:     "and",
	ASSIGN:  "=",
	POW:     "^",
	EQUAL:   "==",
	GEQ:     ">=",
	GT:      ">",
	HASH:    "#",
	LEQ:     "<=",
	LT:      "<",
	MINUS:   "-",
	NEQ:     "~=",
	NOT:     "not",
	OR:      "or",
	PERCENT: "%",
	PLUS:    "+",
	SLASH:   "/",
	STAR:    "*",

	// Structure
	LPAREN: "(",
	RPAREN: ")",
	LBRACK: "[",
	RBRACK: "]",
	LBRACE: "{",
	RBRACE: "}",

	// Grammar
	COLON:     ":",
	COMMA:     ",",
	CONCAT:    "..",
	DOT:       ".",
	SEMICOLON: ";",
	VARARG:    "...",
}

var Reserved = map[string]TokenType{
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
}
