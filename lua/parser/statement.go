package parser

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var stat ast.Statement
	switch p.curToken.Type {
	case token.BREAK:
		stat = p.parseBreakStatement()
	case token.GOTO:
		stat = p.parseGotoStatement()
	case token.IDENT:
		stat = p.parseAssignmentStatement()
	case token.IF:
		stat = p.parseIfStatement()
	case token.LOCAL:
		stat = p.parseLocalStatement()
	default:
		p.errors = append(p.errors, "Unexpected <exp>")
		return nil
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stat
}

func (p *Parser) parseAssignmentStatement() ast.Statement {
	ident := ast.Identifier(p.curToken)
	stmt := &ast.AssignmentStatement{
		Token: p.curToken,
		Name:  ident,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Exps = p.parseExpressionList()

	return stmt
}

func (p *Parser) parseBreakStatement() ast.Statement {
	stmt := ast.BreakStatement(p.curToken)
	return &stmt
}

func (p *Parser) parseGotoStatement() ast.Statement {
	stmt := ast.GotoStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Label = ast.Identifier(p.curToken)
	return &stmt
}

func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

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
