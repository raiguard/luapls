package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	// TODO: Some of these should probably be separate AST types
	var leftExp ast.Expression
	switch p.curToken.Type {
	case token.HASH:
		leftExp = p.parseUnaryExpression()
	case token.IDENT:
		leftExp = p.parseIdentifier()
	case token.LPAREN:
		leftExp = p.parseSurroundingExpression(token.RPAREN)
	case token.MINUS:
		leftExp = p.parseUnaryExpression()
	case token.NOT:
		leftExp = p.parseUnaryExpression()
	case token.NUMBER:
		leftExp = p.parseNumberLiteral()
	case token.STRING:
		leftExp = p.parseStringLiteral()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unable to parse unary expression for token: %s", p.curToken.String()))
		return nil
	}

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

func (p *Parser) parseExpressionList() []ast.Expression {
	exps := []ast.Expression{}
	exp := p.parseExpression(LOWEST)
	if exp == nil {
		return nil
	}
	exps = append(exps, exp)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exp = p.parseExpression(LOWEST)
		if exp == nil {
			// TODO: Error
			return nil
		}
		exps = append(exps, exp)
	}

	return exps
}

func (p *Parser) parseBinaryExpression(left ast.Expression) *ast.BinaryExpression {
	expression := &ast.BinaryExpression{
		Left:     left,
		Operator: p.curToken.Type,
		Right:    nil,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseSurroundingExpression(end token.TokenType) ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(end) {
		return nil
	}
	return exp
}

func (p *Parser) parseUnaryExpression() *ast.UnaryExpression {
	exp := &ast.UnaryExpression{
		Operator: p.curToken.Type,
	}
	p.nextToken()
	exp.Right = p.parseExpression(UNARY)
	return exp
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	ident := ast.Identifier{p.curToken.Literal}
	return &ident
}

func (p *Parser) parseNumberLiteral() *ast.NumberLiteral {
	lit := &ast.NumberLiteral{}

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

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	lit := &ast.StringLiteral{
		Value: strings.Trim(p.curToken.Literal, "\"'"),
	}
	return lit
}

type unaryParseFn func() ast.Expression

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
	UNARY
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
	token.NOT:     UNARY,
	token.HASH:    UNARY,
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
	token.MINUS:   true, // Also a unary operator
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
