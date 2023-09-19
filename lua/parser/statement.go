package parser

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	// TODO: SemicolonStatement so we can show errors
	for p.tokIs(token.SEMICOLON) {
		p.next()
	}
	var stat ast.Statement
	switch p.tok.Type {
	case token.BREAK:
		stat = p.parseBreakStatement()
	case token.DO:
		stat = p.parseDoStatement()
	case token.FOR:
		stat = p.parseForStatement()
	case token.FUNCTION:
		stat = p.parseFunctionStatement(false)
	case token.GOTO:
		stat = p.parseGotoStatement()
	case token.IF:
		stat = p.parseIfStatement()
	case token.LABEL:
		stat = p.parseLabelStatement()
	case token.LOCAL:
		p.next()
		switch p.tok.Type {
		case token.FUNCTION:
			stat = p.parseFunctionStatement(true)
		case token.IDENT:
			stat = p.parseLocalStatement()
		}
	case token.REPEAT:
		stat = p.parseRepeatStatement()
	case token.RETURN:
		stat = p.parseReturnStatement()
	case token.WHILE:
		stat = p.parseWhileStatement()
	default:
		exps := p.parseExpressionList()
		if p.tokIs(token.ASSIGN) {
			stat = p.parseAssignmentStatement(exps)
		} else if len(exps) == 1 {
			stat = &ast.ExpressionStatement{Exp: exps[0]}
		} else {
			p.errors = append(p.errors, fmt.Sprintf("Invalid token: %s", p.tok.Type))
		}
	}
	for p.tokIs(token.SEMICOLON) {
		p.next()
	}
	return stat
}

func (p *Parser) parseAssignmentStatement(vars []ast.Expression) *ast.AssignmentStatement {
	p.expect(token.ASSIGN)
	exps := p.parseExpressionList()

	return &ast.AssignmentStatement{
		Vars: vars,
		Exps: exps,
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	p.expect(token.BREAK)
	return &ast.BreakStatement{}
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	p.expect(token.DO)
	return &ast.DoStatement{Body: p.ParseBlock()}
}

func (p *Parser) parseForStatement() ast.Statement {
	p.expect(token.FOR)
	names := p.parseNameList()

	bareLoop := p.tokIs(token.ASSIGN)
	if bareLoop {
		if len(names) != 1 {
			p.errors = append(p.errors, "Expected 1 identifier")
		}
		p.next()
	} else {
		p.expect(token.IN)
	}

	exps := p.parseExpressionList()

	p.expect(token.DO)

	body := p.ParseBlock()

	p.expect(token.END)

	if bareLoop {
		var start, end, step ast.Expression
		if len(exps) < 1 || len(exps) > 3 {
			p.errors = append(p.errors, "Expected 1 to 3 expressions")
		}
		start = exps[0]
		if len(exps) > 1 {
			end = exps[1]
		}
		if len(exps) > 2 {
			step = exps[2]
		}
		return &ast.ForStatement{
			Name:  names[0],
			Start: start,
			End:   end,
			Step:  step,
			Body:  body,
		}
	} else {
		return &ast.ForInStatement{
			Names: names,
			Exps:  exps,
			Body:  body,
		}
	}
}

func (p *Parser) parseFunctionStatement(isLocal bool) *ast.FunctionStatement {
	p.expect(token.FUNCTION)
	left := p.parseExpression(LOWEST, false)
	p.expect(token.LPAREN)
	params := p.parseNameList()
	p.expect(token.RPAREN)
	body := p.ParseBlock()
	p.expect(token.END)

	return &ast.FunctionStatement{
		Left:    left,
		Params:  params,
		Body:    body,
		IsLocal: isLocal,
	}
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	p.expect(token.GOTO)
	name := p.parseIdentifier()
	return &ast.GotoStatement{Name: name}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	p.expect(token.IF)

	clauses := []*ast.IfClause{p.parseIfClause()}

	for p.tokIs(token.ELSEIF) {
		p.next()
		clauses = append(clauses, p.parseIfClause())
	}

	if p.tokIs(token.ELSE) {
		p.next()
		clauses = append(clauses, &ast.IfClause{Body: p.ParseBlock()})
	}

	p.expect(token.END)

	return &ast.IfStatement{Clauses: clauses}
}

func (p *Parser) parseIfClause() *ast.IfClause {
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.THEN)
	block := p.ParseBlock()
	return &ast.IfClause{
		Condition: condition,
		Body:      block,
	}
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	p.expect(token.LABEL)
	name := p.parseIdentifier()
	p.expect(token.LABEL)
	return &ast.LabelStatement{Name: name}
}

func (p *Parser) parseLocalStatement() *ast.LocalStatement {
	names := p.parseNameList()
	if !p.tokIs(token.ASSIGN) {
		return &ast.LocalStatement{
			Names: names,
			Exps:  []ast.Expression{},
		}
	}

	p.next()
	exps := p.parseExpressionList()

	return &ast.LocalStatement{Names: names, Exps: exps}
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	p.expect(token.REPEAT)
	body := p.ParseBlock()
	p.expect(token.UNTIL)
	condition := p.parseExpression(LOWEST, true)
	return &ast.RepeatStatement{Body: body, Condition: condition}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.expect(token.RETURN)
	return &ast.ReturnStatement{Exps: p.parseExpressionList()}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	p.expect(token.WHILE)
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.DO)
	body := p.ParseBlock()
	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}
