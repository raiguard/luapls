package lua

import protocol "github.com/tliron/glsp/protocol_3_16"

// Parser constructs an AST from Lua source code.
type Parser struct {
    Tree *BlockStatement

    filename string
    scanner *scanner
}

// Temporary type for representing tokens
type TokenExt struct {
    Lit string
    Range protocol.Range
    TokenKind
}

// Initializes the parser.
func (p *Parser) Init(filename string, src []byte) {
    var s scanner
    s.init(src)
    p.filename = filename
    p.scanner = &s
}

// Parses the file.
func (p *Parser) Parse() {
	for {
		tok, _ := p.scanner.scan()
		if tok.Kind == EOF {
			break
		}
	}
}
