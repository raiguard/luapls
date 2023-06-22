package token

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
	CONCAT
	DOT
	SEMICOLON
	SPREAD
)

type Token struct {
	Type    TokenType
	Literal string
	Range   Range
}

type Range struct {
	StartCol, StartRow int
	EndCol, EndRow     int
}

var TokenStr = map[TokenType]string{
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
	FALSE: "false",
	NIL:   "nil",
	TRUE:  "true",

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
