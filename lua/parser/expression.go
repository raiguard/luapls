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
	case token.IDENT:
		left = p.parseIdentifier()
	case token.LBRACE:
		left = p.parseTableLiteral()
	case token.LEN, token.MINUS, token.NOT:
		left = p.parseUnaryExpression()
	case token.LPAREN:
		left = p.parseSurroundingExpression()
	case token.NUMBER:
		left = p.parseNumberLiteral()
	case token.STRING, token.RAWSTRING:
		left = p.parseStringLiteral()
	case token.TRUE, token.FALSE:
		left = p.parseBooleanLiteral()
	case token.VARARG:
		left = p.parseVararg()
	default:
		p.addError(fmt.Sprintf("Unable to parse unary expression for token: %s", p.tok.String()))
		p.next()
		return nil
	}

	for isSuffixOperator(p.tok.Type) {
		switch p.tok.Type {
		case token.LPAREN, token.LBRACE, token.STRING:
			if !allowCall {
				return left
			}
			left = p.parseFunctionCall(left)
		case token.LBRACK, token.DOT, token.COLON:
			left = p.parseIndexExpression(left)
		}
	}

	for isBinaryOperator(p.tok.Type) {
		tokPrecedence := p.tokPrecedence()
		if isRightAssociative(p.tok.Type) {
			tokPrecedence++
		}
		if precedence >= tokPrecedence {
			break
		}
		left = p.parseBinaryExpression(left)
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
		if !p.tokIs(token.VARARG) {
			list = append(list, p.parseIdentifier())
		}
	}

	return list
}

// Parses a namelist, then an optional vararg
func (p *Parser) parseParameterList() ([]*ast.Identifier, bool) {
	names := p.parseNameList()
	vararg := p.tokIs(token.VARARG)
	if vararg {
		p.next()
	}
	return names, vararg
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
	pos := p.tok.Pos
	p.expect(token.FUNCTION)
	p.expect(token.LPAREN)

	params, vararg := p.parseParameterList()

	p.expect(token.RPAREN)

	body := p.ParseBlock()

	end := p.tok.End()
	p.expect(token.END)

	return &ast.FunctionExpression{
		Params:   params,
		Vararg:   vararg,
		Body:     body,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) *ast.IndexExpression {
	isBrackets := p.tokIs(token.LBRACK)
	isColon := p.tokIs(token.COLON)
	p.next()

	var end token.Pos
	var inner ast.Expression
	if isBrackets {
		inner = p.parseExpression(LOWEST, true)
		end = p.tok.Pos
		p.expect(token.RBRACK)
	} else {
		inner = p.parseIdentifier()
		end = inner.End()
	}

	return &ast.IndexExpression{
		Left:       left,
		Inner:      inner,
		IsBrackets: isBrackets,
		IsColon:    isColon,
		EndPos:     end,
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
	pos := p.tok.Pos
	p.next()
	right := p.parseExpression(UNARY, true)
	return &ast.UnaryExpression{Operator: operator, Right: right, StartPos: pos}
}

func (p *Parser) parseBooleanLiteral() *ast.BooleanLiteral {
	value := p.tokIs(token.TRUE)
	pos := p.tok.Pos
	p.next()
	return &ast.BooleanLiteral{Value: value, StartPos: pos}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	ident := ast.Identifier{StartPos: p.tok.Pos}
	if p.tokIs(token.IDENT) {
		ident.Literal = p.tok.Literal
	} else {
		p.expectedTokenError(token.IDENT)
	}
	p.next()
	return &ident
}

func (p *Parser) parseNumberLiteral() *ast.NumberLiteral {
	lit := p.tok.Literal
	pos := p.tok.Pos

	// TODO: Parse all formats of number
	value, err := strconv.ParseFloat(p.tok.Literal, 64)
	if err != nil {
		// msg := fmt.Sprintf("could not parse %q", p.tok.Literal)
		// p.errors = append(p.errors, msg)
	}

	p.next()

	return &ast.NumberLiteral{Literal: lit, Value: float64(value), StartPos: pos}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	lit := &ast.StringLiteral{
		Value:    strings.Trim(p.tok.Literal, "\"'"),
		StartPos: p.tok.Pos,
	}
	p.next()
	return lit
}

func (p *Parser) parseTableLiteral() *ast.TableLiteral {
	pos := p.tok.Pos

	p.expect(token.LBRACE)

	fields := []*ast.TableField{}

	if p.tokIs(token.RBRACE) {
		end := p.tok.End()
		p.next()
		return &ast.TableLiteral{Fields: fields, StartPos: pos, EndPos: end}
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

	end := p.tok.End()
	p.expect(token.RBRACE)

	return &ast.TableLiteral{Fields: fields, StartPos: pos, EndPos: end}
}

func (p *Parser) parseTableField() *ast.TableField {
	var leftExp ast.Expression
	needClosingBracket := false
	pos := p.tok.Pos
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
		return &ast.TableField{Value: leftExp, StartPos: leftExp.Pos()}
	}
	p.expect(token.ASSIGN)
	rightExp := p.parseExpression(LOWEST, true)
	return &ast.TableField{
		Key:      leftExp,
		Value:    rightExp,
		StartPos: pos,
	}
}

func (p *Parser) parseVararg() *ast.Vararg {
	pos := p.tok.Pos
	p.expect(token.VARARG)
	return &ast.Vararg{StartPos: pos}
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
	token.LEN:     UNARY,
	token.POW:     POW,
}

var binaryOperators = map[token.TokenType]bool{
	token.AND:     true,
	token.CONCAT:  true,
	token.EQUAL:   true,
	token.GEQ:     true,
	token.GT:      true,
	token.LEQ:     true,
	token.LT:      true,
	token.MINUS:   true,
	token.NEQ:     true,
	token.OR:      true,
	token.PERCENT: true,
	token.PLUS:    true,
	token.POW:     true,
	token.SLASH:   true,
	token.STAR:    true,
}

func isBinaryOperator(tok token.TokenType) bool {
	return binaryOperators[tok]
}

func isRightAssociative(tok token.TokenType) bool {
	return tok == token.POW || tok == token.CONCAT
}

var suffixOperators = map[token.TokenType]bool{
	token.COLON:  true,
	token.DOT:    true,
	token.LBRACE: true,
	token.LBRACK: true,
	token.LPAREN: true,
	token.STRING: true,
}

func isSuffixOperator(tok token.TokenType) bool {
	return suffixOperators[tok]
}
