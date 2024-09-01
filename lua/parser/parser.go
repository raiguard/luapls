// Package Parser implements a recursive descent parser for Lua 5.2. It is
// heavily based on "Writing an Interpreter in Go" by Thorston Ball.
// https://interpreterbook.com/

// Error recovery: https://supunsetunga.medium.com/writing-a-parser-syntax-error-handling-b71b67a8ac66

package parser

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
)

type ParserError struct {
	Message string
	Range   token.Range
}

func (pe *ParserError) String() string {
	return fmt.Sprintf("%s: %s", &pe.Range, pe.Message)
}

type Parser struct {
	lexer    *lexer.Lexer
	errors   []ParserError
	unit     ast.Unit
	tok      token.Token
	comments []token.Token
}

func New(input string) *Parser {
	p := &Parser{
		lexer:    lexer.New(input),
		errors:   []ParserError{},
		comments: []token.Token{},
	}

	p.tok = p.lexer.Next()
	p.next()

	return p
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) ParseFile() File {
	return File{
		Block:      p.parseBlock(),
		Errors:     p.errors,
		LineBreaks: p.lexer.GetLineBreaks(),
	}
}

func (p *Parser) next() {
	u := ast.Unit{
		LeadingTrivia:  []token.Token{},
		Token:          token.Token{},
		TrailingTrivia: []token.Token{},
	}
	// TODO: Parse type annotations
	for p.tok.Type == token.COMMENT || p.tok.Type == token.WHITESPACE {
		u.LeadingTrivia = append(u.LeadingTrivia, p.tok)
		p.tok = p.lexer.Next()
	}
	u.Token = p.tok
	p.tok = p.lexer.Next()
	for p.tok.Type == token.COMMENT || p.tok.Type == token.WHITESPACE {
		lit := p.tok.Literal
		u.TrailingTrivia = append(u.TrailingTrivia, p.tok)
		p.tok = p.lexer.Next()
		if lit == "\n" {
			break
		}
	}
	p.unit = u
}

func (p *Parser) parseBlock() ast.Block {
	block := ast.Block{
		Stmts:    []ast.Statement{},
		StartPos: p.unit.Token.Pos,
	}

	for !blockEnd[p.unit.Token.Type] {
		stat := p.parseStatement()
		block.Stmts = append(block.Stmts, stat)
		for p.tokIs(token.SEMICOLON) {
			p.next()
		}
	}

	return block
}

func (p *Parser) parseFunctionCall(left ast.Expression) *ast.FunctionCall {
	if p.tokIs(token.STRING) {
		end := p.unit.Token.End()
		args := []ast.Expression{p.parseStringLiteral()}
		return &ast.FunctionCall{Left: left, Args: args, EndPos: end}
	}

	if p.tokIs(token.LBRACE) {
		lit := p.parseTableLiteral()
		return &ast.FunctionCall{
			Left:   left,
			Args:   []ast.Expression{lit},
			EndPos: lit.End(),
		}
	}

	p.expect(token.LPAREN)

	args := []ast.Expression{}
	if !p.tokIs(token.RPAREN) {
		args = p.parseExpressionList()
	}

	end := p.unit.Token.End()

	p.expect(token.RPAREN)

	return &ast.FunctionCall{Left: left, Args: args, EndPos: end}
}

func (p *Parser) expect(tokenType token.TokenType) {
	if !p.tokIs(tokenType) {
		p.expectedTokenError(tokenType)
	}
	// TODO: Smart error recovery
	p.next()
}

func (p *Parser) expectedTokenError(expected token.TokenType) {
	p.addError(
		fmt.Sprintf("Expected %s, got %s",
			token.TokenStr[expected],
			token.TokenStr[p.unit.Token.Type]),
	)
}

func (p *Parser) invalidTokenError() {
	p.addError(fmt.Sprintf("Unexpected %s", token.TokenStr[p.unit.Token.Type]))
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, ParserError{Range: p.unit.Token.Range(), Message: message})
}

func (p *Parser) addErrorForNode(node ast.Node, message string) {
	p.errors = append(p.errors, ParserError{Range: ast.Range(node), Message: message})
}

func (p *Parser) tokIs(tokenType token.TokenType) bool {
	return p.unit.Token.Type == tokenType
}

func (p *Parser) tokPrecedence() operatorPrecedence {
	if p, ok := precedences[p.unit.Token.Type]; ok {
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
