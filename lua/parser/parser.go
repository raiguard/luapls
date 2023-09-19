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

	curToken token.Token
	// TODO: Remove this, it causes lots of smells
	peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseBlock() *ast.Block {
	block := ast.Block{}

	// TODO: Remove the need for this dumbness
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.END) && !p.curTokenIs(token.UNTIL) && !p.curTokenIs(token.ELSEIF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block = append(block, stmt)
		}
		p.nextToken()
	}

	return &block
}

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	fc := ast.FunctionCall{Name: *p.parseIdentifier()}

	if p.peekTokenIs(token.STRING) {
		p.nextToken()
		fc.Args = []ast.Expression{p.parseStringLiteral()}
		return &fc
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	fc.Args = p.parseExpressionList()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return &fc
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expect(tokenType token.TokenType) bool {
	if p.curTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	p.invalidTokenError(tokenType, p.curToken.Type)
	return false
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	p.invalidTokenError(tokenType, p.peekToken.Type)
	return false
}

func (p *Parser) invalidTokenError(expected token.TokenType, got token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("Invalid token: expected %s, got %s",
			token.TokenStr[expected],
			token.TokenStr[got]),
	)
}

func (p *Parser) curPrecedence() operatorPrecedence {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() operatorPrecedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
