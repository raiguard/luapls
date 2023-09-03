package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	prefix := p.getPrefixParser()
	if prefix == nil {
		p.noPrefixParseFnError(p.peekToken.Type)
		return nil
	}
	leftExp := prefix()

	for isBinaryOperator(p.peekToken.Type) {
		peekPrecedence := p.peekPrecedence()
		if isRightAssociative(p.peekToken.Type) {
			peekPrecedence += 1
		}
		if precedence >= peekPrecedence {
			break
		}
		p.nextToken()
		leftExp = p.parseBinaryExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) getPrefixParser() prefixParseFn {
	switch p.curToken.Type {
	case token.HASH:
		return p.parseUnaryExpression
	case token.IDENT:
		return p.parseIdentifier
	case token.LPAREN:
		return p.parseParenthesizedExpression
	case token.MINUS:
		return p.parseUnaryExpression
	case token.NOT:
		return p.parseUnaryExpression
	case token.NUMBER:
		return p.parseNumberLiteral
	case token.STRING:
		return p.parseStringLiteral
	}
	return nil
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
		Right:    nil,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseParenthesizedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	exp := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

// Basic types

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

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{
		Token: p.curToken,
		Value: strings.Trim(p.curToken.Literal, "\"'"),
	}
	return lit
}

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

var binaryOperators = map[token.TokenType]bool{
	token.AND:     true,
	token.CARET:   true,
	token.CONCAT:  true,
	token.EQUAL:   true,
	token.GEQ:     true,
	token.GT:      true,
	token.LEQ:     true,
	token.LT:      true,
	token.MINUS:   true, // Also a prefix operator
	token.NEQ:     true,
	token.OR:      true,
	token.PERCENT: true,
	token.PLUS:    true,
	token.SLASH:   true,
	token.STAR:    true,
}

func isBinaryOperator(tok token.TokenType) bool {
	return binaryOperators[tok]
}

func isRightAssociative(tok token.TokenType) bool {
	return tok == token.CARET || tok == token.CONCAT
}
