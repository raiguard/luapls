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
	tok      token.Token
	comments []token.Token
}

func New(input string) *Parser {
	p := &Parser{
		lexer:    lexer.New(input),
		errors:   []ParserError{},
		comments: []token.Token{},
	}

	p.next()

	return p
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) next() {
	p.tok = p.lexer.Next()
	// TODO: Parse type annotations
	for p.tokIs(token.COMMENT) {
		p.comments = append(p.comments, p.tok)
		p.tok = p.lexer.Next()
	}
}

func (p *Parser) ParseFile() File {
	return File{
		Block:      p.parseBlock(),
		Errors:     p.errors,
		LineBreaks: p.lexer.GetLineBreaks(),
	}
}

func (p *Parser) parseBlock() ast.Block {
	block := ast.Block{
		// TODO: Is this needed somewhere?
		// Comments: p.collectComments(),
		Stmts:    []ast.Statement{},
		StartPos: p.tok.Pos,
	}

	for !blockEnd[p.tok.Type] {
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

func (p *Parser) collectComments() ast.Comments {
	base := ast.Comments{CommentsBefore: p.comments}
	// for _, comment := range p.comments {
	// 	base.CommentsBefore = append(base.CommentsBefore, comment)
	// }
	p.comments = []token.Token{}
	return base
}
