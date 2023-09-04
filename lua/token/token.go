package token

import "fmt"

type TokenType int

const (
	INVALID TokenType = iota
	EOF

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
	CARET
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
	structStart
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	LBRACE
	RBRACE
	structEnd

	// Grammar
	COLON
	COMMA
	DOT
	SEMICOLON
	SPREAD
)

func (t TokenType) String() string {
	return TokenStr[t]
}

type Token struct {
	Type    TokenType
	Literal string
	Range   Range
}

func (t Token) String() string {
	return fmt.Sprintf("[%s] %s %v", TokenStr[t.Type], t.Literal, t.Range)
}

type Range struct {
	StartCol, StartRow uint32
	EndCol, EndRow     uint32
}

var TokenStr = map[TokenType]string{
	IDENT: "identifier",
	LABEL: "label",

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
	CARET:   "^",
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
	SPREAD:    "...",
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
