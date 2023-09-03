// Package Parser implements a recursive descent parser for Lua 5.2. It is
// heavily based on "Writing an Interpreter in Go" by Thorston Ball.
// https://interpreterbook.com/

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
)

type Parser struct {
	lexer *lexer.Lexer

	errors []string

	curToken  token.Token
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

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseBlock() *ast.Block {
	block := &ast.Block{}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.END) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		} else {
			p.nextToken() // Always make some progress
		}
	}

	return block
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors,
		fmt.Sprintf("Invalid token: expected %s, got %s",
			token.TokenStr[p.peekToken.Type],
			token.TokenStr[tokenType]),
	)
	// TODO: Error
	return false
}

func (p *Parser) peekPrecedence() operatorPrecedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", token.TokenStr[t])
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		return p.parseAssignmentStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LOCAL:
		return p.parseLocalStatement()
	default:
		return nil
	}
}

func (p *Parser) parseAssignmentStatement() ast.Statement {
	ident := ast.Identifier(p.curToken)
	stmt := &ast.AssignmentStatement{
		Token: p.curToken,
		Name:  &ident,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST) // TODO: We don't have required semicolons in Lua

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()

	consequence := p.ParseBlock()

	if consequence == nil {
		p.errors = append(p.errors, "Failed to parse block")
		return nil
	}

	stmt.Consequence = *consequence

	if !p.curTokenIs(token.END) {
		p.errors = append(p.errors, fmt.Sprintf("Expected 'end', got %s", p.curToken.Literal))
		return nil
	}

	p.nextToken()

	return stmt
}

func (p *Parser) parseLocalStatement() ast.Statement {
	stmt := &ast.LocalStatement{
		Token: p.curToken,
	}
	p.nextToken()
	switch p.curToken.Type {
	case token.IDENT:
		stmt.Statement = p.parseAssignmentStatement()
	default:
		p.errors = append(p.errors, fmt.Sprintf("Invalid token in local statement: %s", token.TokenStr[p.curToken.Type]))
	}
	return stmt
}

func (p *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	prefix := p.getPrefixParser()
	if prefix == nil {
		p.noPrefixParseFnError(p.peekToken.Type)
		return nil
	}
	leftExp := prefix()

	for token.IsOperator(p.peekToken.Type) && precedence < p.peekPrecedence() {
		infix := p.getInfixParser()
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) getPrefixParser() prefixParseFn {
	switch p.curToken.Type {
	case token.IDENT:
		return p.parseIdentifier
	case token.NUMBER:
		return p.parseNumberLiteral
	case token.STRING:
		return p.parseStringLiteral
	}
	return nil
}

func (p *Parser) getInfixParser() infixParseFn {
	return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := ast.Identifier(p.curToken)
	return &ident
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}

	// TODO: Handle all kinds of number
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = float64(value)

	p.nextToken()

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{
		Token: p.curToken,
		Value: strings.Trim(p.curToken.Literal, "\"'"),
	}
	p.nextToken()
	return lit
}

type infixParseFn func(ast.Expression) ast.Expression
type prefixParseFn func() ast.Expression

type operatorPrecedence int

const (
	_ operatorPrecedence = iota
	LOWEST
	ASSIGN
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]operatorPrecedence{
	token.ASSIGN: ASSIGN,
	token.NEQ:    ASSIGN,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.STAR:   PRODUCT,
	token.LPAREN: CALL,
}
