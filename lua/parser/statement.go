package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var stat ast.Statement
	switch p.curToken.Type {
	case token.BREAK:
		stat = p.parseBreakStatement()
	case token.DO:
		stat = p.parseDoStatement()
	case token.FOR:
		p.nextToken()
		if p.peekTokenIs(token.ASSIGN) {
			stat = p.parseForStatement()
		} else if p.peekTokenIs(token.COMMA) {
			stat = p.parseForInStatement()
		}
	case token.GOTO:
		stat = p.parseGotoStatement()
	case token.IDENT:
		stat = p.parseAssignmentStatement()
	case token.IF:
		stat = p.parseIfStatement()
	case token.LABEL:
		stat = p.parseLabelStatement()
	case token.LOCAL:
		stat = p.parseLocalStatement()
	case token.REPEAT:
		stat = p.parseRepeatStatement()
	case token.WHILE:
		stat = p.parseWhileStatement()
	}
	if stat == nil {
		p.errors = append(p.errors, "Unexpected <exp>")
		return nil
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stat
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{
		Token: p.curToken,
	}
	stmt.Vars = parseNodeList(p, p.parseIdentifier)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Exps = p.parseExpressionList()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := ast.BreakStatement(p.curToken)
	return &stmt
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	stmt := ast.DoStatement{Token: p.curToken}
	p.nextToken()
	stmt.Block = *p.ParseBlock()
	return &stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	// Parser will have advanced to the next token, but tests will not
	if p.curTokenIs(token.FOR) {
		p.nextToken()
	}
	if !p.curTokenIs(token.IDENT) {
		p.invalidTokenError(token.IDENT, p.curToken.Type)
		return nil
	}
	stmt := ast.ForStatement{Var: *p.parseIdentifier()}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	start := p.parseExpression(LOWEST)
	if start == nil {
		return nil
	}
	stmt.Start = start
	if !p.expectPeek(token.COMMA) {
		return nil
	}
	p.nextToken()
	end := p.parseExpression(LOWEST)
	if end == nil {
		return nil
	}
	stmt.End = end
	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		step := p.parseExpression(LOWEST)
		if step == nil {
			return nil
		}
		stmt.Step = &step
	}
	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken()
	block := p.ParseBlock()
	if block == nil {
		return nil
	}
	stmt.Block = *block
	if !p.curTokenIs(token.END) {
		p.invalidTokenError(token.END, p.curToken.Type)
		return nil
	}

	return &stmt
}

func (p *Parser) parseForInStatement() *ast.ForInStatement {
	// Parser will have advanced to the next token, but tests will not
	if p.curTokenIs(token.FOR) {
		p.nextToken()
	}
	if !p.curTokenIs(token.IDENT) {
		p.invalidTokenError(token.IDENT, p.curToken.Type)
		return nil
	}
	stmt := ast.ForInStatement{}
	stmt.Vars = parseNodeList(p, p.parseIdentifier)
	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()
	stmt.Exps = p.parseExpressionList()
	if !p.expectPeek(token.DO) {
		return nil
	}
	p.nextToken()
	block := p.ParseBlock()
	if block == nil {
		return nil
	}
	stmt.Block = *block

	if !p.curTokenIs(token.END) {
		p.invalidTokenError(token.END, p.curToken.Type)
		return nil
	}

	return &stmt
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	stmt := ast.GotoStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Label = *p.parseIdentifier()
	return &stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{
		Token:   p.curToken,
		Clauses: []ast.IfClause{},
	}

	p.nextToken()

	stmt.Clauses = append(stmt.Clauses, *p.parseIfClause())

	for p.curTokenIs(token.ELSEIF) {
		p.nextToken()
		stmt.Clauses = append(stmt.Clauses, *p.parseIfClause())
	}

	if !p.curTokenIs(token.END) {
		p.invalidTokenError(token.END, p.curToken.Type)
		return nil
	}

	return stmt
}

func (p *Parser) parseIfClause() *ast.IfClause {
	condition := p.parseExpression(LOWEST)
	if condition == nil {
		return nil
	}

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()

	block := p.ParseBlock()

	if block == nil {
		p.errors = append(p.errors, "Failed to parse block")
		return nil
	}

	clause := ast.IfClause{
		Condition: condition,
		Block:     *block,
	}
	return &clause
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	stmt := ast.LabelStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Label = *p.parseIdentifier()
	if !p.expectPeek(token.LABEL) {
		return nil
	}
	return &stmt
}

func (p *Parser) parseLocalStatement() *ast.LocalStatement {
	stmt := &ast.LocalStatement{
		Token: p.curToken,
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Names = parseNodeList(p, p.parseIdentifier)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Exps = p.parseExpressionList()

	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	stmt := &ast.RepeatStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Block = *p.ParseBlock()
	if !p.curTokenIs(token.UNTIL) {
		return nil
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.DO) {
		return nil
	}
	p.nextToken()
	stmt.Block = *p.ParseBlock()
	return stmt
}

func parseNodeList[T ast.Node](p *Parser, parseFunc func() *T) []T {
	values := []T{}
	val := parseFunc()
	if val == nil {
		return values
	}
	values = append(values, *val)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		val = parseFunc()
		if val == nil {
			break
		}
		values = append(values, *val)
	}

	return values
}
