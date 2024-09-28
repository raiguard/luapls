// Package Parser implements a recursive descent parser for Lua 5.2. It is
// heavily based on "Writing an Interpreter in Go" by Thorston Ball.
// https://interpreterbook.com/

// Error recovery: https://supunsetunga.medium.com/writing-a-parser-syntax-error-handling-b71b67a8ac66

package parser

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/util"
)

type ParserError struct {
	Message string
	Range   token.Range
}

func (pe *ParserError) String() string {
	return fmt.Sprintf("%s: %s", &pe.Range, pe.Message)
}

type Parser struct {
	errors     []ParserError
	lineBreaks []int
	units      []ast.Unit
	pos        int
}

func New(input string) *Parser {
	units, lineBreaks := Run(input)
	p := &Parser{
		errors:     []ParserError{},
		lineBreaks: lineBreaks,
		units:      units,
	}

	return p
}

func Run(input string) ([]ast.Unit, []int) {
	// Consume all tokens and convert them into units
	tokens, lineBreaks := lexer.Run(input)
	units := []ast.Unit{}
	u := ast.Unit{
		LeadingTrivia:  []token.Token{},
		Token:          token.Token{},
		TrailingTrivia: []token.Token{},
	}
	state := "leading"
	newUnit := func() {
		units = append(units, u)
		u = ast.Unit{
			LeadingTrivia:  []token.Token{},
			Token:          token.Token{},
			TrailingTrivia: []token.Token{},
		}
	}
	for _, tok := range tokens {
		if tok.Type == token.COMMENT || tok.Type == token.WHITESPACE {
			if state == "leading" {
				u.LeadingTrivia = append(u.LeadingTrivia, tok)
			} else {
				u.TrailingTrivia = append(u.TrailingTrivia, tok)
				if strings.Contains(tok.Literal, "\n") {
					newUnit()
					state = "leading"
				}
			}
		} else {
			if state == "leading" {
				state = "trailing"
				u.Token = tok
			} else {
				newUnit()
				u.Token = tok
			}
		}
	}
	newUnit()

	return units, lineBreaks
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) ParseFile() File {
	return File{
		Block:      p.parseBlock(),
		Errors:     p.errors,
		LineBreaks: p.lineBreaks,
	}
}

func (p *Parser) unit() *ast.Unit {
	return &p.units[p.pos]
}

func (p *Parser) next() ast.Unit {
	if p.pos < len(p.units)-1 {
		p.pos++
	}
	return p.units[p.pos]
}

func (p *Parser) parseBlock() ast.Block {
	block := ast.Block{}

	for {
		pair := ast.Pair[ast.Statement]{Node: p.parseStatement()}
		if p.tokIs(token.SEMICOLON) {
			pair.Delimeter = p.unit()
			p.next()
		}
		block.Pairs = append(block.Pairs, pair)
		if blockEnd[p.unit().Type()] {
			break
		}
	}

	return block
}

func (p *Parser) parseFunctionCall(name ast.Expression) *ast.FunctionCall {
	fc := &ast.FunctionCall{Name: name}
	if p.tokIs(token.STRING) {
		fc.Args = ast.SimplePunctuated[ast.Expression](p.parseTableLiteral())
		return fc
	}

	if p.tokIs(token.LBRACE) {
		fc.Args = ast.SimplePunctuated[ast.Expression](p.parseTableLiteral())
		return fc
	}

	fc.LeftParen = util.Ptr(p.expect(token.LPAREN))

	if rparen := p.accept(token.RPAREN); rparen != nil {
		fc.RightParen = rparen
		return fc
	}

	fc.Args = p.parseExpressionList()

	fc.RightParen = util.Ptr(p.expect(token.RPAREN))

	return fc
}

func (p *Parser) accept(tokenType token.TokenType) *ast.Unit {
	if !p.tokIs(tokenType) {
		return nil
	}
	return util.Ptr(p.next())
}

func (p *Parser) expect(tokenType token.TokenType) ast.Unit {
	if !p.tokIs(tokenType) {
		p.expectedTokenError(tokenType)
	}
	// TODO: Smart error recovery
	p.next()
	return p.units[p.pos-1]
}

func (p *Parser) expectedTokenError(expected token.TokenType) {
	p.addError(
		fmt.Sprintf("Expected %s, got %s",
			token.TokenStr[expected],
			token.TokenStr[p.unit().Type()]),
	)
}

func (p *Parser) invalidTokenError() {
	p.addError(fmt.Sprintf("Unexpected %s", token.TokenStr[p.unit().Type()]))
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, ParserError{Range: p.unit().Token.Range(), Message: message})
}

func (p *Parser) addErrorForNode(node ast.Node, message string) {
	p.errors = append(p.errors, ParserError{Range: ast.Range(node), Message: message})
}

func (p *Parser) tokIs(tokenType token.TokenType) bool {
	return p.unit().Type() == tokenType
}

func (p *Parser) tokPrecedence() operatorPrecedence {
	if p, ok := precedences[p.unit().Type()]; ok {
		return p
	}
	return LOWEST
}

var blockEnd = map[token.TokenType]bool{
	token.ELSEIF: true,
	token.ELSE:   true,
	token.END:    true,
	token.EOF:    true,
	token.UNTIL:  true,
}
