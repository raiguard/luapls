package lua

import protocol "github.com/tliron/glsp/protocol_3_16"

// Parser constructs an AST from Lua source code.
type Parser struct {
    // Tree *Block
    Tokens []TokenExt

    filename string
    scanner *Scanner
}

// Temporary type for representing tokens
type TokenExt struct {
    Lit string
    Range protocol.Range
    Token
}

// Initializes the parser.
func (p *Parser) Init(filename string, src []byte) {
    var s Scanner
    s.Init(src)
    p.filename = filename
    p.scanner = &s
}

// Parses the file.
func (p *Parser) Parse() {
	for {
		rng, tok, lit := p.scanner.Scan()
		if tok == EOF {
			break
		}
		p.Tokens = append(p.Tokens, TokenExt{lit, rng, tok})
	}
}
