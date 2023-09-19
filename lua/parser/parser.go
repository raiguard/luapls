// Package Parser implements a recursive descent parser for Lua 5.2. It is
// heavily based on "Writing an Interpreter in Go" by Thorston Ball.
// https://interpreterbook.com/

package parser

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
)

type Parser struct {
	lexer *lexer.Lexer

	errors []string

	tok token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.next()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) next() {
	p.tok = p.lexer.NextToken()
}

func (p *Parser) ParseBlock() ast.Block {
	block := ast.Block{}

	for !blockEnd[p.tok.Type] {
		block = append(block, p.parseStatement())
	}

	return block
}

func (p *Parser) parseFunctionCall(left ast.Expression) *ast.FunctionCall {
	if p.tokIs(token.STRING) {
		args := []ast.Expression{p.parseStringLiteral()}
		p.next()
		return &ast.FunctionCall{Left: left, Args: args}
	}

	p.expect(token.LPAREN)

	args := []ast.Expression{}
	if !p.tokIs(token.RPAREN) {
		args = p.parseExpressionList()
	}

	p.expect(token.RPAREN)

	return &ast.FunctionCall{Left: left, Args: args}
}

func (p *Parser) expect(tokenType token.TokenType) {
	if !p.tokIs(tokenType) {
		p.invalidTokenError(tokenType)
	}
	p.next()
}

// Like expect, but doesn't advance
func (p *Parser) expect0(tokenType token.TokenType) {
	if !p.tokIs(tokenType) {
		p.invalidTokenError(tokenType)
	}
}

func (p *Parser) invalidTokenError(expected token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("Invalid token: expected %s, got %s",
			token.TokenStr[expected],
			p.tok.String()),
	)
}

func (p *Parser) tokIs(tokenType token.TokenType) bool {
	return p.tok.Type == tokenType
}

func (p *Parser) tokPrecedence() operatorPrecedence {
	if p, ok := precedences[p.tok.Type]; ok {
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
