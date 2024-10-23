package token

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TokenType int

const (
	INVALID TokenType = iota
	EOF
	COMMENT
	WHITESPACE

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
	LEN
	LEQ
	LT
	MINUS
	NEQ
	NOT
	OR
	MOD
	PLUS
	SLASH
	MUL

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

	// Annotation
	DOC_CLASS
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

var InvalidPos = Pos(-1)

func (t Token) End() int {
	return t.Pos + len(t.Literal)
}

func (t *Token) Range() Range {
	return Range{t.Pos, t.End()}
}

func (t *Token) String() string {
	return fmt.Sprintf("[%s] `%s` %v", TokenStr[t.Type], strings.NewReplacer("\n", "\\n", "\t", "\\t").Replace(t.Literal), t.Pos)
}

type Range struct {
	Start int
	End   int
}

func (r *Range) String() string {
	return fmt.Sprintf("%v:%v", r.Start, r.End)
}

func (r *Range) ContainsPos(pos Pos) bool {
	return r.Start <= pos && r.End > pos
}

func (r *Range) ContainsRange(rng Range) bool {
	return r.Start <= rng.Start && r.End >= rng.End
}

var TokenStr = map[TokenType]string{
	INVALID:    "invalid",
	EOF:        "eof",
	IDENT:      "identifier",
	LABEL:      "label",
	COMMENT:    "comment",
	WHITESPACE: "whitespace",

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
	AND:    "and",
	ASSIGN: "assign",
	POW:    "pow",
	EQUAL:  "equal",
	GEQ:    "geq",
	GT:     "gt",
	LEN:    "len",
	LEQ:    "leq",
	LT:     "lt",
	MINUS:  "minus",
	NEQ:    "neq",
	NOT:    "not",
	OR:     "or",
	MOD:    "mod",
	PLUS:   "plus",
	SLASH:  "slash",
	MUL:    "mul",

	// Structure
	LPAREN: "left paren",
	RPAREN: "right paren",
	LBRACK: "left bracket",
	RBRACK: "right bracket",
	LBRACE: "left brace",
	RBRACE: "right brace",

	// Grammar
	COLON:     "colon",
	COMMA:     "comma",
	CONCAT:    "concat",
	DOT:       "dot",
	SEMICOLON: "semicolon",
	VARARG:    "vararg",

	// Annotation
	DOC_CLASS: "@class",
}

var Reserved = map[string]TokenType{
	"and":      AND,
	"break":    BREAK,
	"do":       DO,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"end":      END,
	"false":    FALSE,
	"for":      FOR,
	"function": FUNCTION,
	"goto":     GOTO,
	"if":       IF,
	"in":       IN,
	"local":    LOCAL,
	"nil":      NIL,
	"not":      NOT,
	"or":       OR,
	"repeat":   REPEAT,
	"return":   RETURN,
	"then":     THEN,
	"true":     TRUE,
	"until":    UNTIL,
	"while":    WHILE,

	"@class": DOC_CLASS,
}
