package annotation

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Parse(src string) (Annotation, []ast.Diagnostic) {
	tokens, _ := lexer.Run(src)
	p := parser{tokens, -1, []ast.Diagnostic{}}
	return p.parse()
}

type parser struct {
	tokens      []token.Token
	pos         int
	diagnostics []ast.Diagnostic
}

func (p *parser) parse() (Annotation, []ast.Diagnostic) {
	tok := p.next()
	switch tok.Type {
	case token.DOC_CLASS:
		name := p.expect(token.IDENT)
		return &Class{Name: name.Literal}, p.diagnostics
	case token.INVALID:
		p.diagnostics = append(p.diagnostics, ast.Diagnostic{
			Message:  "Unknown annotation",
			Range:    tok.Range(),
			Severity: protocol.DiagnosticSeverityWarning,
		})
		return nil, p.diagnostics
	default:
		return nil, p.diagnostics
	}
}

func (p *parser) read() *token.Token {
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return &p.tokens[p.pos]
}

func (p *parser) next() *token.Token {
	tok := p.read()
	for tok.Type == token.WHITESPACE {
		tok = p.read()
	}
	return tok
}

func (p *parser) expect(typ token.TokenType) *token.Token {
	tok := p.next()
	if tok.Type != typ {
		p.diagnostics = append(p.diagnostics, ast.Diagnostic{
			Message:  fmt.Sprintf("Expected %s", token.TokenStr[typ]),
			Range:    tok.Range(),
			Severity: protocol.DiagnosticSeverityWarning,
		})
	}
	return tok
}
