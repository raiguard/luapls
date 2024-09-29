package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/util"
)

func (p *Parser) parseExpression(precedence operatorPrecedence, allowCall bool) ast.Expression {
	var left ast.Expression
	switch p.unit().Type() {
	case token.FUNCTION:
		left = p.parseFunctionExpression()
	case token.IDENT:
		left = p.parseIdentifier()
	case token.LBRACE:
		left = p.parseTableLiteral()
	case token.LEN, token.MINUS, token.NOT:
		left = p.parsePrefixExpression()
	case token.LPAREN:
		left = p.parseSurroundingExpression()
	case token.NIL:
		left = p.parseNilLiteral()
	case token.NUMBER:
		left = p.parseNumberLiteral()
	case token.STRING, token.RAWSTRING:
		left = p.parseStringLiteral()
	case token.TRUE, token.FALSE:
		left = p.parseBooleanLiteral()
	case token.VARARG:
		left = p.parseVararg()
	default:
		invalid := ast.Invalid{Position: p.unit().Pos()}
		p.addError("Expected expression")
		p.next()
		return &invalid
	}

	for isSuffixOperator(p.unit().Type()) {
		switch p.unit().Type() {
		case token.LPAREN, token.LBRACE, token.STRING:
			if !allowCall {
				return left
			}
			left = p.parseFunctionCall(left)
		case token.LBRACK, token.DOT, token.COLON:
			left = p.parseIndexExpression(left)
		}
	}

	for isInfixOperator(p.unit().Type()) {
		tokPrecedence := p.tokPrecedence()
		if isRightAssociative(p.unit().Type()) {
			tokPrecedence++
		}
		if precedence >= tokPrecedence {
			break
		}
		left = p.parseInfixExpression(left)
	}

	return left
}

func (p *Parser) parseExpressionList() ast.Punctuated[ast.Expression] {
	exps := ast.Punctuated[ast.Expression]{StartPos: p.unit().Pos()}

	for {
		pair := ast.Pair[ast.Expression]{
			Node: p.parseExpression(LOWEST, true),
		}
		if p.tokIs(token.COMMA) {
			pair.Delimeter = p.unit()
			p.next()
		}
		exps.Pairs = append(exps.Pairs, pair)
		if pair.Delimeter == nil {
			break
		}
	}

	return exps
}

// Identical to parseExpressionList, but only for identifiers
func (p *Parser) parseNameList() ast.Punctuated[*ast.Identifier] {
	list := ast.Punctuated[*ast.Identifier]{StartPos: p.unit().Pos()}

	for {
		if !p.tokIs(token.IDENT) {
			break
		}
		pair := ast.Pair[*ast.Identifier]{Node: p.parseIdentifier()}
		if p.tokIs(token.COMMA) {
			pair.Delimeter = p.unit()
			p.next()
		}
		list.Pairs = append(list.Pairs, pair)
		if pair.Delimeter == nil {
			break
		}
	}

	return list
}

// Parses a namelist, then an optional vararg
func (p *Parser) parseParameterList() (ast.Punctuated[*ast.Identifier], *ast.Unit) {
	return p.parseNameList(), p.accept(token.VARARG)
}

func (p *Parser) parseFunctionExpression() *ast.FunctionExpression {
	function := p.expect(token.FUNCTION)
	lparen := p.expect(token.LPAREN)

	params, vararg := p.parseParameterList()

	rparen := p.expect(token.RPAREN)

	body := p.parseBlock()

	end := p.expect(token.END)

	return &ast.FunctionExpression{
		Function:   function,
		LeftParen:  lparen,
		Params:     params,
		Vararg:     vararg,
		RightParen: rparen,
		Body:       body,
		EndUnit:    end,
	}
}

func (p *Parser) parseIndexExpression(prefix ast.Expression) *ast.IndexExpression {
	ie := &ast.IndexExpression{
		Prefix:      prefix,
		LeftIndexer: *p.unit(),
	}
	p.next()

	if ie.LeftIndexer.Type() == token.LBRACK {
		ie.Inner = p.parseExpression(LOWEST, true)
		ie.RightIndexer = util.Ptr(p.expect(token.RBRACK))
	} else {
		ie.Inner = p.parseIdentifier()
	}

	return ie
}

func (p *Parser) parseInfixExpression(left ast.Expression) *ast.InfixExpression {
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: *p.unit(),
		Right:    nil,
	}

	precedence := p.tokPrecedence()
	p.next()
	expression.Right = p.parseExpression(precedence, true)

	return expression
}

func (p *Parser) parseSurroundingExpression() ast.Expression {
	p.expect(token.LPAREN)
	exp := p.parseExpression(LOWEST, true)
	p.expect(token.RPAREN)
	return exp
}

func (p *Parser) parsePrefixExpression() *ast.PrefixExpression {
	operator := *p.unit()
	p.next()
	right := p.parseExpression(PREFIX, true)
	return &ast.PrefixExpression{Operator: operator, Right: right}
}

func (p *Parser) parseBooleanLiteral() *ast.BooleanLiteral {
	bl := ast.BooleanLiteral(*p.unit())
	p.next()
	return &bl
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return util.Ptr(ast.Identifier(p.expect(token.IDENT)))
}

func (p *Parser) parseNilLiteral() *ast.NilLiteral {
	return util.Ptr(ast.NilLiteral(p.expect(token.NIL)))
}

func (p *Parser) parseNumberLiteral() *ast.NumberLiteral {
	return util.Ptr(ast.NumberLiteral(p.expect(token.NUMBER)))
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	return util.Ptr(ast.StringLiteral(p.expect(token.STRING)))
}

func (p *Parser) parseTableLiteral() *ast.TableLiteral {
	tl := &ast.TableLiteral{LeftBrace: p.expect(token.LBRACE)}

	if rbrace := p.accept(token.RBRACE); rbrace != nil {
		tl.RightBrace = *rbrace
		return tl
	}

	tl.Fields = p.parseTableFieldList()
	tl.RightBrace = p.expect(token.RBRACE)

	return tl
}

func isType[T any](thing any) bool {
	_, ok := thing.(*T)
	return ok
}

func (p *Parser) parseVararg() *ast.Vararg {
	return util.Ptr(ast.Vararg(p.expect(token.VARARG)))
}

var tableSep = map[token.TokenType]bool{
	token.COMMA:     true,
	token.SEMICOLON: true,
}

type prefixParseFn func() ast.Expression

type operatorPrecedence int

const (
	LOWEST operatorPrecedence = iota
	OR
	AND
	CMP
	CONCAT
	SUM
	PRODUCT
	PREFIX
	POW
)

var precedences = map[token.TokenType]operatorPrecedence{
	token.OR:     OR,
	token.AND:    AND,
	token.LT:     CMP,
	token.GT:     CMP,
	token.LEQ:    CMP,
	token.GEQ:    CMP,
	token.NEQ:    CMP,
	token.EQUAL:  CMP,
	token.CONCAT: CONCAT,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.MUL:    PRODUCT,
	token.SLASH:  PRODUCT,
	token.MOD:    PRODUCT,
	token.NOT:    PREFIX,
	token.LEN:    PREFIX,
	token.POW:    POW,
}

var infixOperators = map[token.TokenType]bool{
	token.AND:    true,
	token.CONCAT: true,
	token.EQUAL:  true,
	token.GEQ:    true,
	token.GT:     true,
	token.LEQ:    true,
	token.LT:     true,
	token.MINUS:  true,
	token.NEQ:    true,
	token.OR:     true,
	token.MOD:    true,
	token.PLUS:   true,
	token.POW:    true,
	token.SLASH:  true,
	token.MUL:    true,
}

func isInfixOperator(tok token.TokenType) bool {
	return infixOperators[tok]
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
