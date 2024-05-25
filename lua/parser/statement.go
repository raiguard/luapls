package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.tok.Type {
	case token.BREAK:
		return p.parseBreakStatement()
	case token.DO:
		return p.parseDoStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement(p.tok.Pos, false)
	case token.GOTO:
		return p.parseGotoStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LABEL:
		return p.parseLabelStatement()
	case token.LOCAL:
		tok := p.tok
		p.next()
		switch p.tok.Type {
		case token.FUNCTION:
			return p.parseFunctionStatement(tok.Pos, true)
		case token.IDENT:
			return p.parseLocalStatement(tok.Pos)
		}
		stat := &ast.Invalid{Token: &tok}
		p.addErrorForNode(stat, "Invalid statement")
		return stat
	case token.REPEAT:
		return p.parseRepeatStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	}

	exps := p.parseExpressionList()
	if p.tokIs(token.ASSIGN) {
		return p.parseAssignmentStatement(exps)
		// TODO: Can there ever be zero expressions?
	} else if _, ok := exps[0].(*ast.FunctionCall); ok {
		return exps[0].(*ast.FunctionCall)
	} else {
		stat := &ast.Invalid{Exps: exps}
		p.addErrorForNode(stat, "Invalid statement")
		return stat
	}
}

func (p *Parser) parseAssignmentStatement(vars []ast.Expression) *ast.AssignmentStatement {
	p.expect(token.ASSIGN)
	exps := p.parseExpressionList()

	return &ast.AssignmentStatement{
		Comments: p.collectComments(),
		Vars:     vars,
		Exps:     exps,
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	pos := p.tok.Pos
	p.expect(token.BREAK)
	return &ast.BreakStatement{
		Comments: p.collectComments(),
		StartPos: pos,
	}
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	comments := p.collectComments()
	pos := p.tok.Pos
	p.expect(token.DO)
	block := p.parseBlock()
	end := p.tok.End()
	p.expect(token.END)
	return &ast.DoStatement{
		Comments: comments,
		Body:     block,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	comments := p.collectComments()
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

	body := p.parseBlock()

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
			Comments: comments,
			Name:     names[0],
			Start:    start,
			Finish:   finish,
			Step:     step,
			Body:     body,
			StartPos: pos,
			EndPos:   end,
		}
	}

	return &ast.ForInStatement{
		Comments: comments,
		Names:    names,
		Exps:     exps,
		Body:     body,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseFunctionStatement(pos token.Pos, isLocal bool) *ast.FunctionStatement {
	comments := p.collectComments()
	p.expect(token.FUNCTION)
	left := p.parseExpression(LOWEST, false)
	p.expect(token.LPAREN)
	params, vararg := p.parseParameterList()
	p.expect(token.RPAREN)
	body := p.parseBlock()
	end := p.tok.End()
	p.expect(token.END)

	return &ast.FunctionStatement{
		Comments: comments,
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
		Comments: p.collectComments(),
		Name:     name,
		StartPos: pos,
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	comments := p.collectComments()
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
		body := p.parseBlock()
		clauses = append(clauses, &ast.IfClause{Body: body, StartPos: clausePos, EndPos: body.End()})
	}

	end := p.tok.End()
	p.expect(token.END)

	return &ast.IfStatement{
		Comments: comments,
		Clauses:  clauses,
		StartPos: pos,
		EndPos:   end,
	}
}

// TODO: Comments
func (p *Parser) parseIfClause(pos token.Pos) *ast.IfClause {
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.THEN)
	block := p.parseBlock()
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
	return &ast.LabelStatement{
		Comments: p.collectComments(),
		Name:     name,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseLocalStatement(pos token.Pos) *ast.LocalStatement {
	comments := p.collectComments()
	names := p.parseNameList()
	if !p.tokIs(token.ASSIGN) {
		return &ast.LocalStatement{
			Comments: comments,
			Names:    names,
			Exps:     []ast.Expression{},
			StartPos: pos,
		}
	}

	p.next()
	exps := p.parseExpressionList()

	return &ast.LocalStatement{
		Comments: comments,
		Names:    names,
		Exps:     exps,
		StartPos: pos,
	}
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	comments := p.collectComments()
	pos := p.tok.Pos
	p.expect(token.REPEAT)
	body := p.parseBlock()
	p.expect(token.UNTIL)
	condition := p.parseExpression(LOWEST, true)
	return &ast.RepeatStatement{
		Comments:  comments,
		Body:      body,
		Condition: condition,
		StartPos:  pos,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	comments := p.collectComments()
	pos := p.tok.Pos
	p.expect(token.RETURN)
	if !blockEnd[p.tok.Type] {
		return &ast.ReturnStatement{
			Comments: comments,
			Exps:     p.parseExpressionList(),
			StartPos: pos,
		}
	} else {
		return &ast.ReturnStatement{
			Comments: comments,
			Exps:     []ast.Expression{},
			StartPos: pos,
		}
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	comments := p.collectComments()
	pos := p.tok.Pos
	p.expect(token.WHILE)
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.DO)
	body := p.parseBlock()
	end := p.tok.End()
	p.expect(token.END)
	return &ast.WhileStatement{
		Comments:  comments,
		Condition: condition,
		Body:      body,
		StartPos:  pos,
		EndPos:    end,
	}
}
