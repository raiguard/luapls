package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseExpression(precedence operatorPrecedence, allowCall bool) ast.Expression {
	var left ast.Expression
	switch p.tok.Type {
	case token.FUNCTION:
		left = p.parseFunctionExpression()
	case token.TRUE, token.FALSE:
		left = p.parseBooleanLiteral()
	case token.HASH:
		left = p.parseUnaryExpression()
	case token.IDENT:
		left = p.parseIdentifier()
	case token.LPAREN:
		left = p.parseSurroundingExpression()
	case token.MINUS:
		left = p.parseUnaryExpression()
	case token.NOT:
		left = p.parseUnaryExpression()
	case token.NUMBER:
		left = p.parseNumberLiteral()
	case token.STRING, token.RAWSTRING:
		left = p.parseStringLiteral()
	case token.LBRACE:
		left = p.parseTableLiteral()
	case token.VARARG:
		left = p.parseVararg()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unable to parse unary expression for token: %s", p.tok.String()))
		p.next()
		return nil
	}

	for isBinaryOperator(p.tok.Type) {
		tokPrecedence := p.tokPrecedence()
		if isRightAssociative(p.tok.Type) {
			tokPrecedence += 1
		}
		if precedence >= tokPrecedence {
			break
		}
		if p.tokIs(token.LPAREN) && !allowCall {
			break
		}
		switch p.tok.Type {
		case token.LPAREN, token.LBRACE:
			left = p.parseFunctionCall(left)
		case token.LBRACK, token.DOT, token.COLON:
			left = p.parseIndexExpression(left)
		default:
			left = p.parseBinaryExpression(left)
		}
	}

	return left
}

func (p *Parser) parseExpressionList() []ast.Expression {
	exps := []ast.Expression{p.parseExpression(LOWEST, true)}

	for p.tokIs(token.COMMA) {
		p.next()
		exps = append(exps, p.parseExpression(LOWEST, true))
	}

	return exps
}

// Identical to parseExpressionList, but only for identifiers
func (p *Parser) parseNameList() []*ast.Identifier {
	list := []*ast.Identifier{}

	if !p.tokIs(token.IDENT) {
		return list
	}

	list = append(list, p.parseIdentifier())

	for p.tokIs(token.COMMA) {
		p.next()
		list = append(list, p.parseIdentifier())
	}

	return list
}

func (p *Parser) parseBinaryExpression(left ast.Expression) *ast.BinaryExpression {
	expression := &ast.BinaryExpression{
		Left:     left,
		Operator: p.tok.Type,
		Right:    nil,
	}

	precedence := p.tokPrecedence()
	p.next()
	expression.Right = p.parseExpression(precedence, true)

	return expression
}

func (p *Parser) parseFunctionExpression() *ast.FunctionExpression {
	p.expect(token.FUNCTION)
	p.expect(token.LPAREN)

	params := p.parseNameList()

	p.expect(token.RPAREN)

	body := p.ParseBlock()

	p.expect(token.END)

	return &ast.FunctionExpression{
		Params: params,
		Body:   body,
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) *ast.IndexExpression {
	isBrackets := p.tokIs(token.LBRACK)
	isColon := p.tokIs(token.COLON)
	p.next()

	var inner ast.Expression
	if isBrackets {
		inner = p.parseExpression(LOWEST, true)
		p.expect(token.RBRACK)
	} else {
		inner = p.parseIdentifier()
	}

	return &ast.IndexExpression{
		Left:       left,
		Inner:      inner,
		IsBrackets: isBrackets,
		IsColon:    isColon,
	}
}

func (p *Parser) parseSurroundingExpression() ast.Expression {
	p.expect(token.LPAREN)
	exp := p.parseExpression(LOWEST, true)
	p.expect(token.RPAREN)
	return exp
}

func (p *Parser) parseUnaryExpression() *ast.UnaryExpression {
	operator := p.tok.Type
	p.next()
	right := p.parseExpression(UNARY, true)
	return &ast.UnaryExpression{Operator: operator, Right: right}
}

func (p *Parser) parseBooleanLiteral() *ast.BooleanLiteral {
	// If this returns an error, something has gone VERY wrong, so just explode
	value, err := strconv.ParseBool(p.tok.Literal)
	if err != nil {
		panic(err)
	}
	p.next()
	return &ast.BooleanLiteral{Value: value}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	var ident *ast.Identifier
	if p.tokIs(token.IDENT) {
		ident = &ast.Identifier{Literal: p.tok.Literal}
	} else {
		p.invalidTokenError(token.IDENT)
	}
	p.next()
	return ident
}

func (p *Parser) parseNumberLiteral() *ast.NumberLiteral {
	lit := p.tok.Literal

	// TODO: Parse all formats of number
	value, err := strconv.ParseFloat(p.tok.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q", p.tok.Literal)
		p.errors = append(p.errors, msg)
	}

	p.next()

	return &ast.NumberLiteral{Literal: lit, Value: float64(value)}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	lit := &ast.StringLiteral{
		Value: strings.Trim(p.tok.Literal, "\"'"),
	}
	p.next()
	return lit
}

func (p *Parser) parseTableLiteral() *ast.TableLiteral {
	p.expect(token.LBRACE)

	fields := []*ast.TableField{}

	if p.tokIs(token.RBRACE) {
		p.next()
		return &ast.TableLiteral{Fields: fields}
	}

	fields = append(fields, p.parseTableField())

	for tableSep[p.tok.Type] {
		p.next()
		// Trailing separator
		if p.tokIs(token.RBRACE) {
			break
		}
		fields = append(fields, p.parseTableField())
	}

	p.expect(token.RBRACE)

	return &ast.TableLiteral{Fields: fields}
}

func (p *Parser) parseTableField() *ast.TableField {
	var leftExp ast.Expression
	needClosingBracket := false
	if p.tokIs(token.LBRACK) {
		needClosingBracket = true
		p.next()
	}
	leftExp = p.parseExpression(LOWEST, true)
	if needClosingBracket {
		p.expect(token.RBRACK)
	}
	// If key is in brackets, value is required
	if !needClosingBracket && (tableSep[p.tok.Type] || p.tokIs(token.RBRACE)) {
		return &ast.TableField{Value: leftExp}
	}
	p.expect(token.ASSIGN)
	rightExp := p.parseExpression(LOWEST, true)
	return &ast.TableField{
		Key:   leftExp,
		Value: rightExp,
	}
}

func (p *Parser) parseVararg() *ast.Vararg {
	p.expect(token.VARARG)
	return &ast.Vararg{}
}

var tableSep = map[token.TokenType]bool{
	token.COMMA:     true,
	token.SEMICOLON: true,
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
	INDEX
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
	token.COLON:   INDEX,
	token.DOT:     INDEX,
	token.LBRACK:  INDEX,
	token.LPAREN:  INDEX,
	token.LBRACE:  INDEX,
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
	token.COLON:   true,
	token.DOT:     true,
	token.LBRACK:  true,
	token.LPAREN:  true,
	token.LBRACE:  true,
}

func isBinaryOperator(tok token.TokenType) bool {
	return binaryOperators[tok]
}

func isRightAssociative(tok token.TokenType) bool {
	return tok == token.CARET || tok == token.CONCAT
}
