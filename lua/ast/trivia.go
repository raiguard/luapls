package ast

import "github.com/raiguard/luapls/lua/token"

func (node *AssignmentStatement) GetLeadingTrivia() []token.Token {
	return node.Vars.GetLeadingTrivia()
}

func (node *BooleanLiteral) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *BreakStatement) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *DoStatement) GetLeadingTrivia() []token.Token {
	return node.DoTok.LeadingTrivia
}

func (node *ForInStatement) GetLeadingTrivia() []token.Token {
	return node.ForTok.LeadingTrivia
}

func (node *ForStatement) GetLeadingTrivia() []token.Token {
	return node.ForTok.LeadingTrivia
}

func (node *FunctionCall) GetLeadingTrivia() []token.Token {
	return node.Name.GetLeadingTrivia()
}

func (node *FunctionExpression) GetLeadingTrivia() []token.Token {
	return node.FuncTok.LeadingTrivia
}

func (node *FunctionStatement) GetLeadingTrivia() []token.Token {
	if node.LocalTok != nil {
		return node.LocalTok.LeadingTrivia
	}
	return node.FuncTok.LeadingTrivia
}

func (node *GotoStatement) GetLeadingTrivia() []token.Token {
	return node.GotoTok.LeadingTrivia
}

func (node *Identifier) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *IfClause) GetLeadingTrivia() []token.Token {
	return node.LeadingTok.LeadingTrivia
}

func (node *IfStatement) GetLeadingTrivia() []token.Token {
	return node.IfTok.LeadingTrivia
}

func (node *IndexExpression) GetLeadingTrivia() []token.Token {
	return node.Prefix.GetLeadingTrivia()
}

func (node *InfixExpression) GetLeadingTrivia() []token.Token {
	return node.Left.GetLeadingTrivia()
}

func (node *LabelStatement) GetLeadingTrivia() []token.Token {
	return node.LeadingLabelTok.LeadingTrivia
}

func (node *Invalid) GetLeadingTrivia() []token.Token {
	return []token.Token{} // TODO:
}

func (node *LocalStatement) GetLeadingTrivia() []token.Token {
	return node.LocalTok.LeadingTrivia
}

func (node *NilLiteral) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *NumberLiteral) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *Pair[T]) GetLeadingTrivia() []token.Token {
	return node.Node.GetLeadingTrivia()
}

func (node *Punctuated[T]) GetLeadingTrivia() []token.Token {
	if len(node.Pairs) == 0 {
		return []token.Token{}
	}
	return node.Pairs[0].GetLeadingTrivia()
}

func (node *PrefixExpression) GetLeadingTrivia() []token.Token {
	return node.Operator.LeadingTrivia
}

func (node *RepeatStatement) GetLeadingTrivia() []token.Token {
	return node.RepeatTok.LeadingTrivia
}

func (node *ReturnStatement) GetLeadingTrivia() []token.Token {
	return node.ReturnTok.LeadingTrivia
}

func (node *SemicolonStatement) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *StringLiteral) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *TableArrayField) GetLeadingTrivia() []token.Token {
	return node.Expr.GetLeadingTrivia()
}

func (node *TableSimpleKeyField) GetLeadingTrivia() []token.Token {
	return node.Name.LeadingTrivia
}

func (node *TableExpressionKeyField) GetLeadingTrivia() []token.Token {
	return node.LeftBracket.LeadingTrivia
}

func (node *TableLiteral) GetLeadingTrivia() []token.Token {
	return node.LeftBrace.LeadingTrivia
}

func (node *Vararg) GetLeadingTrivia() []token.Token {
	return node.LeadingTrivia
}

func (node *WhileStatement) GetLeadingTrivia() []token.Token {
	return node.WhileTok.LeadingTrivia
}
