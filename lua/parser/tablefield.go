package parser

import (
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/token"
)

func (p *Parser) parseTableFieldList() ast.Punctuated[ast.TableField] {
	exps := ast.Punctuated[ast.TableField]{StartPos: p.unit().Pos()}

	for {
		if p.tokIs(token.RBRACE) {
			break
		}
		pair := ast.Pair[ast.TableField]{
			Node: p.parseTableField(),
		}
		if tableSep[p.unit().Type()] {
			pair.Delimeter = p.unit()
			p.next()
		}
		exps.Pairs = append(exps.Pairs, pair)
		if pair.Delimeter == nil {
			break
		}
	}

	return exps
}

func (p *Parser) parseTableField() ast.TableField {
	if lbrack := p.accept(token.LBRACK); lbrack != nil {
		name := p.parseExpression(LOWEST, true)
		rbrack := p.expect(token.RBRACK)
		assignTok := p.expect(token.ASSIGN)
		expr := p.parseExpression(LOWEST, true)
		return &ast.TableExpressionKeyField{
			LeftBracket:  *lbrack,
			Name:         name,
			RightBracket: rbrack,
			AssignTok:    assignTok,
			Expr:         expr,
		}
	}

	expr := p.parseExpression(LOWEST, true)

	if !p.tokIs(token.ASSIGN) {
		return &ast.TableArrayField{Expr: expr}
	}

	assignTok := p.expect(token.ASSIGN)

	name, ok := expr.(*ast.Identifier)
	if !ok {
		panic("TODO: Table key was not an identifier!")
	}

	expr = p.parseExpression(LOWEST, true)

	return &ast.TableSimpleKeyField{
		Name:      *name,
		AssignTok: assignTok,
		Expr:      expr,
	}
}
