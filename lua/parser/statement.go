package parser

import (
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
		stat = p.parseFunctionStatement(p.tok.Pos, false)
	case token.GOTO:
		stat = p.parseGotoStatement()
	case token.IF:
		stat = p.parseIfStatement()
	case token.LABEL:
		stat = p.parseLabelStatement()
	case token.LOCAL:
		tok := p.tok
		p.next()
		switch p.tok.Type {
		case token.FUNCTION:
			stat = p.parseFunctionStatement(tok.Pos, true)
		case token.IDENT:
			stat = p.parseLocalStatement(tok.Pos)
		default:
			stat = &ast.Invalid{Token: &tok}
			p.addErrorForNode(stat, "Invalid statement")
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
			// TODO: Can there ever be zero expressions?
		} else if _, ok := exps[0].(*ast.FunctionCall); ok {
			stat = exps[0].(*ast.FunctionCall)
		} else {
			stat = &ast.Invalid{Exps: exps}
			p.addErrorForNode(stat, "Invalid statement")
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
	pos := p.tok.Pos
	p.expect(token.BREAK)
	return &ast.BreakStatement{
		StartPos: pos,
	}
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	pos := p.tok.Pos
	p.expect(token.DO)
	block := p.ParseBlock()
	end := p.tok.End()
	p.expect(token.END)
	return &ast.DoStatement{
		Body:     block,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	pos := p.tok.Pos
	p.expect(token.FOR)
	names := p.parseNameList()

	bareLoop := len(names) == 1 && p.tokIs(token.ASSIGN)
	if bareLoop {
		p.next()
	} else {
		p.expect(token.IN)
	}

	exps := p.parseExpressionList()

	p.expect(token.DO)

	body := p.ParseBlock()

	end := p.tok.End()
	p.expect(token.END)

	if bareLoop {
		var start, finish, step ast.Expression
		if len(exps) < 1 || len(exps) > 3 {
			p.addError("Expected 1 to 3 expressions")
		}
		start = exps[0]
		if len(exps) > 1 {
			finish = exps[1]
		}
		if len(exps) > 2 {
			step = exps[2]
		}
		return &ast.ForStatement{
			Name:     names[0],
			Start:    start,
			Finish:   finish,
			Step:     step,
			Body:     body,
			StartPos: pos,
			EndPos:   end,
		}
	} else {
		return &ast.ForInStatement{
			Names:    names,
			Exps:     exps,
			Body:     body,
			StartPos: pos,
			EndPos:   end,
		}
	}
}

func (p *Parser) parseFunctionStatement(pos token.Pos, isLocal bool) *ast.FunctionStatement {
	p.expect(token.FUNCTION)
	left := p.parseExpression(LOWEST, false)
	p.expect(token.LPAREN)
	params, vararg := p.parseParameterList()
	p.expect(token.RPAREN)
	body := p.ParseBlock()
	end := p.tok.End()
	p.expect(token.END)

	return &ast.FunctionStatement{
		Left:     left,
		Params:   params,
		Vararg:   vararg,
		Body:     body,
		IsLocal:  isLocal,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	pos := p.tok.Pos
	p.expect(token.GOTO)
	name := p.parseIdentifier()
	return &ast.GotoStatement{
		Name:     name,
		StartPos: pos,
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	pos := p.tok.Pos
	p.expect(token.IF)

	clauses := []*ast.IfClause{p.parseIfClause(pos)}

	for p.tokIs(token.ELSEIF) {
		clausePos := p.tok.Pos
		p.next()
		clauses = append(clauses, p.parseIfClause(clausePos))
	}

	if p.tokIs(token.ELSE) {
		clausePos := p.tok.Pos
		p.next()
		body := p.ParseBlock()
		clauses = append(clauses, &ast.IfClause{Body: body, StartPos: clausePos, EndPos: body.End()})
	}

	end := p.tok.End()
	p.expect(token.END)

	return &ast.IfStatement{Clauses: clauses, StartPos: pos, EndPos: end}
}

func (p *Parser) parseIfClause(pos token.Pos) *ast.IfClause {
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.THEN)
	block := p.ParseBlock()
	return &ast.IfClause{
		Condition: condition,
		Body:      block,
		StartPos:  pos,
		EndPos:    block.End(),
	}
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	pos := p.tok.Pos
	p.expect(token.LABEL)
	name := p.parseIdentifier()
	end := p.tok.End()
	p.expect(token.LABEL)
	return &ast.LabelStatement{Name: name, StartPos: pos, EndPos: end}
}

func (p *Parser) parseLocalStatement(pos token.Pos) *ast.LocalStatement {
	names := p.parseNameList()
	if !p.tokIs(token.ASSIGN) {
		return &ast.LocalStatement{
			Names:    names,
			Exps:     []ast.Expression{},
			StartPos: pos,
		}
	}

	p.next()
	exps := p.parseExpressionList()

	return &ast.LocalStatement{
		Names:    names,
		Exps:     exps,
		StartPos: pos,
	}
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	pos := p.tok.Pos
	p.expect(token.REPEAT)
	body := p.ParseBlock()
	p.expect(token.UNTIL)
	condition := p.parseExpression(LOWEST, true)
	return &ast.RepeatStatement{
		Body:      body,
		Condition: condition,
		StartPos:  pos,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	pos := p.tok.Pos
	p.expect(token.RETURN)
	if !blockEnd[p.tok.Type] {
		return &ast.ReturnStatement{
			Exps:     p.parseExpressionList(),
			StartPos: pos,
		}
	} else {
		return &ast.ReturnStatement{
			Exps:     []ast.Expression{},
			StartPos: pos,
		}
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	pos := p.tok.Pos
	p.expect(token.WHILE)
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.DO)
	body := p.ParseBlock()
	end := p.tok.End()
	p.expect(token.END)
	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
		StartPos:  pos,
		EndPos:    end,
	}
}
