package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/util"
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
		return p.parseFunctionStatement(nil)
	case token.GOTO:
		return p.parseGotoStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LABEL:
		return p.parseLabelStatement()
	case token.LOCAL:
		tok := p.expect(token.LOCAL)
		switch p.unit().Type() {
		case token.FUNCTION:
			return p.parseFunctionStatement(&tok)
		case token.IDENT:
			return p.parseLocalStatement(tok)
		}
		stat := &ast.Invalid{Unit: &tok}
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
	node := util.Ptr(ast.BreakStatement(*p.unit()))
	p.expect(token.BREAK)
	return node
}

func (p *Parser) parseDoStatement() *ast.DoStatement {
	do := p.expect(token.DO)
	block := p.parseBlock()
	end := p.expect(token.END)
	return &ast.DoStatement{
		DoTok:  do,
		Body:   block,
		EndTok: end,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	forTok := p.expect(token.FOR)
	names := p.parseNameList()

	var sepTok ast.Unit

	bareLoop := len(names.Pairs) == 1 && p.tokIs(token.ASSIGN)
	if bareLoop {
		sepTok = p.expect(token.ASSIGN)
	} else {
		sepTok = p.expect(token.IN)
	}

	exps := p.parseExpressionList()
	doTok := p.expect(token.DO)
	body := p.parseBlock()
	endTok := p.expect(token.END)

	if bareLoop {
		var start, finish ast.Pair[ast.Expression]
		if len(exps.Pairs) < 2 || len(exps.Pairs) > 3 {
			p.addError("Expected 2 to 3 expressions")
		}
		start = exps.Pairs[0]
		finish = exps.Pairs[1]
		var step *ast.Pair[ast.Expression]
		if len(exps.Pairs) > 2 {
			step = &exps.Pairs[2]
		}
		return &ast.ForStatement{
			ForTok:    forTok,
			Name:      names.Pairs[0].Node,
			AssignTok: sepTok,
			Start:     start,
			Finish:    finish,
			Step:      step,
			DoTok:     doTok,
			Body:      body,
			EndTok:    endTok,
		}
	}

	return &ast.ForInStatement{
		ForTok: forTok,
		Names:  names,
		InTok:  sepTok,
		Exps:   exps,
		DoTok:  doTok,
		Body:   body,
		EndTok: endTok,
	}
}

func (p *Parser) parseFunctionStatement(localTok *ast.Unit) *ast.FunctionStatement {
	funcTok := p.expect(token.FUNCTION)
	name := p.parseExpression(LOWEST, false)
	lparen := p.expect(token.LPAREN)
	params, vararg := p.parseParameterList()
	rparen := p.expect(token.RPAREN)
	body := p.parseBlock()
	endTok := p.expect(token.END)

	return &ast.FunctionStatement{
		LocalTok:   localTok,
		FuncTok:    funcTok,
		Name:       name,
		LeftParen:  lparen,
		Params:     params,
		Vararg:     vararg,
		RightParen: rparen,
		Body:       body,
		EndTok:     endTok,
	}
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	gotoTok := p.expect(token.GOTO)
	name := p.parseIdentifier()
	return &ast.GotoStatement{
		GotoTok: gotoTok,
		Name:    name,
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	ifTok := p.expect(token.IF)

	clauses := []*ast.IfClause{p.parseIfClause(ifTok)}

	for p.tokIs(token.ELSEIF) {
		elseifTok := p.expect(token.ELSEIF)
		clauses = append(clauses, p.parseIfClause(elseifTok))
	}

	if p.tokIs(token.ELSE) {
		elseTok := p.expect(token.ELSE)
		body := p.parseBlock()
		clauses = append(clauses, &ast.IfClause{
			LeadingTok: elseTok,
			Condition:  nil,
			ThenTok:    nil,
			Body:       body,
		})
	}

	endTok := p.expect(token.END)

	return &ast.IfStatement{
		IfTok:   ifTok,
		Clauses: clauses,
		EndTok:  endTok,
	}
}

func (p *Parser) parseIfClause(leadingTok ast.Unit) *ast.IfClause {
	condition := p.parseExpression(LOWEST, true)
	thenTok := p.expect(token.THEN)
	block := p.parseBlock()
	return &ast.IfClause{
		LeadingTok: leadingTok,
		Condition:  condition,
		ThenTok:    &thenTok,
		Body:       block,
	}
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	leadingLabelTok := p.expect(token.LABEL)
	name := p.parseIdentifier()
	trailingLabelTok := p.expect(token.LABEL)
	return &ast.LabelStatement{
		LeadingLabelTok:  leadingLabelTok,
		Name:             name,
		TrailingLabelTok: trailingLabelTok,
	}
}

func (p *Parser) parseLocalStatement(localTok ast.Unit) *ast.LocalStatement {
	ls := &ast.LocalStatement{
		LocalTok: localTok,
		Names:    p.parseNameList(),
	}
	if assignTok := p.accept(token.ASSIGN); assignTok != nil {
		ls.AssignTok = assignTok
		ls.Exps = util.Ptr(p.parseExpressionList())
	}
	return ls
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	repeatTok := p.expect(token.REPEAT)
	body := p.parseBlock()
	untilTok := p.expect(token.UNTIL)
	condition := p.parseExpression(LOWEST, true)
	return &ast.RepeatStatement{
		RepeatTok: repeatTok,
		Body:      body,
		UntilTok:  untilTok,
		Condition: condition,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnTok := p.expect(token.RETURN)
	rs := &ast.ReturnStatement{ReturnTok: returnTok}
	if !blockEnd[p.unit().Type()] {
		rs.Exps = util.Ptr(p.parseExpressionList())
	}
	return rs
}

func (p *Parser) parseSemicolonStatement() *ast.SemicolonStatement {
	return util.Ptr(ast.SemicolonStatement(p.expect(token.SEMICOLON)))
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	whileTok := p.expect(token.WHILE)
	condition := p.parseExpression(LOWEST, true)
	doTok := p.expect(token.DO)
	body := p.parseBlock()
	endTok := p.expect(token.END)
	return &ast.WhileStatement{
		WhileTok:  whileTok,
		Condition: condition,
		DoTok:     doTok,
		Body:      body,
		EndTok:    endTok,
	}
}
