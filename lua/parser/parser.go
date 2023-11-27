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

type Parser struct {
	lexer  *lexer.Lexer
	errors []ParserError
	tok    token.Token
}

func New(input string) *Parser {
	p := &Parser{
		lexer:  lexer.New(input),
		errors: []ParserError{},
	}

	p.next()

	return p
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) next() {
	p.tok = p.lexer.NextToken()
	// TODO: Parse type annotations
	for p.tokIs(token.COMMENT) {
		p.tok = p.lexer.NextToken()
	}
}

func (p *Parser) ParseFile() File {
	return File{
		Block:      p.ParseBlock(),
		Errors:     p.errors,
		LineBreaks: p.lexer.GetLineBreaks(),
	}
}

// TODO: Make private
func (p *Parser) ParseBlock() ast.Block {
	block := ast.Block{
		Stmts:    []ast.Statement{},
		StartPos: p.tok.Pos,
	}

	for !blockEnd[p.tok.Type] {
		block.Stmts = append(block.Stmts, p.parseStatement())
	}

	return block
}

func (p *Parser) parseFunctionCall(left ast.Expression) *ast.FunctionCall {
	if p.tokIs(token.STRING) {
		end := p.tok.End()
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

	end := p.tok.End()

	p.expect(token.RPAREN)

	return &ast.FunctionCall{Left: left, Args: args, EndPos: end}
}

func (p *Parser) expect(tokenType token.TokenType) token.Token {
	if !p.tokIs(tokenType) {
		p.expectedTokenError(tokenType)
	}
	current := p.tok
	p.next()
	return current
}

func (p *Parser) expectedTokenError(expected token.TokenType) {
	p.addError(
		fmt.Sprintf("Expected %s, got %s",
			token.TokenStr[expected],
			token.TokenStr[p.tok.Type]),
	)
}

func (p *Parser) invalidTokenError() {
	p.addError(fmt.Sprintf("Unexpected %s", token.TokenStr[p.tok.Type]))
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, ParserError{Range: p.tok.Range(), Message: message})
}

func (p *Parser) addErrorForNode(node ast.Node, message string) {
	p.errors = append(p.errors, ParserError{Range: ast.Range(node), Message: message})
}

func (p *Parser) tokIs(tokenType token.TokenType) bool {
	return p.tok.Type == tokenType
}

func (p *Parser) tokPrecedence() operatorPrecedence {
	if p, ok := precedences[p.tok.Type]; ok {
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
