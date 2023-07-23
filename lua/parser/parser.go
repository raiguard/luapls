package parser

import (
	"fmt"
	"strconv"

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
	// case token.LOCAL:
	// 	return p.parseLocalStatement()
	case token.IDENT:
		return p.parseAssignmentStatement()
	default:
		return nil
	}
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	ident := ast.Identifier(p.curToken)
	stmt := ast.AssignmentStatement{
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

	return &stmt
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
