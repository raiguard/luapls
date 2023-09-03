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

func (p *Parser) Errors() []string {
	return p.errors
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
		}
		p.nextToken()
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
			token.TokenStr[tokenType],
			token.TokenStr[p.peekToken.Type]),
	)
	// TODO: Error
	return false
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
		p.errors = append(p.errors, "Unexpected <exp>")
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

	for token.IsInfixOperator(p.peekToken.Type) {
		peekPrecedence := p.peekPrecedence()
		if token.IsRightAssociative(p.peekToken.Type) {
			peekPrecedence += 1
		}
		if precedence >= peekPrecedence {
			break
		}
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) getPrefixParser() prefixParseFn {
	switch p.curToken.Type {
	case token.HASH:
		return p.parsePrefixExpression
	case token.IDENT:
		return p.parseIdentifier
	case token.LPAREN:
		return p.parseParenthesizedExpression
	case token.MINUS:
		return p.parsePrefixExpression
	case token.NOT:
		return p.parsePrefixExpression
	case token.NUMBER:
		return p.parseNumberLiteral
	case token.STRING:
		return p.parseStringLiteral
	}
	return nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
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

	return lit
}

func (p *Parser) parseParenthesizedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{
		Token: p.curToken,
		Value: strings.Trim(p.curToken.Literal, "\"'"),
	}
	return lit
}

type infixParseFn func(ast.Expression) ast.Expression
type prefixParseFn func() ast.Expression

type operatorPrecedence int

const (
	_ operatorPrecedence = iota
	LOWEST
	OR
	AND
	CMP
	CONCAT
	SUM
	PRODUCT
	PREFIX
	POW
	CALL
)

var precedences = map[token.TokenType]operatorPrecedence{
	token.OR:      OR,
	token.AND:     AND,
	token.LT:      CMP,
	token.GT:      CMP,
	token.LEQ:     CMP,
	token.GEQ:     CMP,
	token.NEQ:     CMP,
	token.EQUAL:   CMP,
	token.CONCAT:  CONCAT,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.STAR:    PRODUCT,
	token.SLASH:   PRODUCT,
	token.PERCENT: PRODUCT,
	token.NOT:     PREFIX,
	token.HASH:    PREFIX,
	token.CARET:   POW,
	token.LPAREN:  CALL,
	token.LBRACK:  CALL,
}
