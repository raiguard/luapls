package lua

// import protocol "github.com/tliron/glsp/protocol_3_16"

// Parser constructs an AST from Lua source code.
type Parser struct {
	Tree *BlockStmt

	filename string
	scanner  *scanner
	tok      Token
	tokRaw   string
}

// Initializes the parser.
func (p *Parser) Init(filename string, src []byte) {
	p.filename = filename
	p.Tree = &BlockStmt{}

	var s scanner
	s.init(src)
	p.scanner = &s
}

// Parses the file.
func (p *Parser) Parse() {
	p.parseBlock()
}

// Advance to the next token.
func (p *Parser) next() {
	p.tok, p.tokRaw = p.scanner.scan()
	// TODO: Handle scanner errors
}

func (p *Parser) parseBlock() {
	for {
		p.next() // Always make progress
		switch p.tok.Kind {
		case LOCAL:
			p.parseLocalStmt()
		default:
			// TODO:
		}
	}
}

// LocalStmt, LocalFunctionStmt
func (p *Parser) parseLocalStmt() {
	p.next()
	switch p.tok.Kind {
	case FUNCTION:
		// TODO:
	case IDENT:
		namelist := p.parseIdentList()
		explist := p.parseExprList()
	}
}

func (p *Parser) parseIdentList() IdentList {
	// TODO: Ensure that an ident comes first
	idents := []Ident{p.consumeIdent()}
	p.next()
	for p.tok.Kind == COMMA {
		idents = append(idents, p.consumeIdent())
	}

	return IdentList{
		Idents: idents,
		Start:  idents[0].Pos(),
		Stop:   idents[len(idents)-1].End(),
	}
}

func (p *Parser) consumeIdent() Ident {
	ident := Ident{Raw: p.tokRaw, TokenPos: p.tok.Pos}
	p.next()
	return ident
}

func (p *Parser) parseExprList() ExprList {
	panic("TODO")
}
