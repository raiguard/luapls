package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.unit().Type() {
	case token.BREAK:
		return p.parseBreakStatement()
	case token.DO:
		return p.parseDoStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement(p.unit().Token.Pos, false)
	case token.GOTO:
		return p.parseGotoStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LABEL:
		return p.parseLabelStatement()
	case token.LOCAL:
		tok := p.unit().Token
		p.next()
		switch p.unit().Type() {
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
	case token.SEMICOLON:
		return p.parseSemicolonStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	}

	exps := p.parseExpressionList()
	if p.tokIs(token.ASSIGN) {
		return p.parseAssignmentStatement(exps)
		// TODO: Can there ever be zero expressions?
	} else if fc, ok := exps.Pairs[0].Node.(*ast.FunctionCall); ok {
		return fc
	} else {
		stat := &ast.Invalid{Exps: exps}
		p.addErrorForNode(stat, "Invalid statement")
		return stat
	}
}

func (p *Parser) parseAssignmentStatement(vars ast.Punctuated[ast.Expression]) *ast.AssignmentStatement {
	assign := p.expect(token.ASSIGN)
	exps := p.parseExpressionList()

	return &ast.AssignmentStatement{
		Vars:   vars,
		Assign: assign,
		Exps:   exps,
	}
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	pos := p.unit().Token.Pos
	p.expect(token.BREAK)
	return &ast.BreakStatement{
		StartPos: pos,
	}
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	pos := p.unit().Token.Pos
	p.expect(token.DO)
	block := p.parseBlock()
	end := p.unit().Token.End()
	p.expect(token.END)
	return &ast.DoStatement{
		Body:     block,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	pos := p.unit().Token.Pos
	p.expect(token.FOR)
	names := p.parseNameList()

	bareLoop := len(names.Pairs) == 1 && p.tokIs(token.ASSIGN)
	if bareLoop {
		p.next()
	} else {
		p.expect(token.IN)
	}

	exps := p.parseExpressionList()

	p.expect(token.DO)

	body := p.parseBlock()

	end := p.unit().Token.End()
	p.expect(token.END)

	if bareLoop {
		var start, finish, step ast.Expression
		if len(exps.Pairs) < 1 || len(exps.Pairs) > 3 {
			p.addError("Expected 1 to 3 expressions")
		}
		start = exps.Pairs[0].Node
		if len(exps.Pairs) > 1 {
			finish = exps.Pairs[1].Node
		}
		if len(exps.Pairs) > 2 {
			step = exps.Pairs[2].Node
		}
		return &ast.ForStatement{
			Name:     names.Pairs[0].Node,
			Start:    start,
			Finish:   finish,
			Step:     step,
			Body:     body,
			StartPos: pos,
			EndPos:   end,
		}
	}

	return &ast.ForInStatement{
		Names:    names,
		Exps:     exps,
		Body:     body,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseFunctionStatement(pos token.Pos, isLocal bool) *ast.FunctionStatement {
	p.expect(token.FUNCTION)
	left := p.parseExpression(LOWEST, false)
	p.expect(token.LPAREN)
	params, vararg := p.parseParameterList()
	p.expect(token.RPAREN)
	body := p.parseBlock()
	end := p.unit().Token.End()
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
	pos := p.unit().Token.Pos
	p.expect(token.GOTO)
	name := p.parseIdentifier()
	return &ast.GotoStatement{
		Name:     name,
		StartPos: pos,
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	pos := p.unit().Token.Pos
	p.expect(token.IF)

	clauses := []*ast.IfClause{p.parseIfClause(pos)}

	for p.tokIs(token.ELSEIF) {
		clausePos := p.unit().Token.Pos
		p.next()
		clauses = append(clauses, p.parseIfClause(clausePos))
	}

	if p.tokIs(token.ELSE) {
		clausePos := p.unit().Token.Pos
		p.next()
		body := p.parseBlock()
		clauses = append(clauses, &ast.IfClause{Body: body, StartPos: clausePos, EndPos: body.End()})
	}

	end := p.unit().Token.End()
	p.expect(token.END)

	return &ast.IfStatement{
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
	pos := p.unit().Token.Pos
	p.expect(token.LABEL)
	name := p.parseIdentifier()
	end := p.unit().Token.End()
	p.expect(token.LABEL)
	return &ast.LabelStatement{
		Name:     name,
		StartPos: pos,
		EndPos:   end,
	}
}

func (p *Parser) parseLocalStatement(pos token.Pos) *ast.LocalStatement {
	names := p.parseNameList()
	if !p.tokIs(token.ASSIGN) {
		return &ast.LocalStatement{
			Names:    names,
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
	pos := p.unit().Token.Pos
	p.expect(token.REPEAT)
	body := p.parseBlock()
	p.expect(token.UNTIL)
	condition := p.parseExpression(LOWEST, true)
	return &ast.RepeatStatement{
		Body:      body,
		Condition: condition,
		StartPos:  pos,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	pos := p.unit().Token.Pos
	p.expect(token.RETURN)
	if !blockEnd[p.unit().Type()] {
		return &ast.ReturnStatement{
			Exps:     p.parseExpressionList(),
			StartPos: pos,
		}
	} else {
		return &ast.ReturnStatement{StartPos: pos}
	}
}

func (p *Parser) parseSemicolonStatement() *ast.SemicolonStatement {
	ss := &ast.SemicolonStatement{
		Unit: *p.unit(),
	}
	p.next()
	return ss
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	pos := p.unit().Token.Pos
	p.expect(token.WHILE)
	condition := p.parseExpression(LOWEST, true)
	p.expect(token.DO)
	body := p.parseBlock()
	end := p.unit().Token.End()
	p.expect(token.END)
	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
		StartPos:  pos,
		EndPos:    end,
	}
}
