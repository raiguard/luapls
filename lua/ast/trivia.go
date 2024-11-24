package ast

import "github.com/raiguard/luapls/lua/token"

type Trivia interface {
	triviaNode()
}

type SimpleTrivia struct{ Tok token.Token }

func (st *SimpleTrivia) triviaNode()    {}
func (st *SimpleTrivia) Pos() token.Pos { return st.Tok.Pos }
func (st *SimpleTrivia) End() token.Pos { return st.Tok.End() }

func (node *AssignmentStatement) GetLeadingTrivia() []Trivia {
	return node.Vars.GetLeadingTrivia()
}

func (node *BooleanLiteral) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *BreakStatement) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *DoStatement) GetLeadingTrivia() []Trivia {
	return node.DoTok.LeadingTrivia
}

func (node *ForInStatement) GetLeadingTrivia() []Trivia {
	return node.ForTok.LeadingTrivia
}

func (node *ForStatement) GetLeadingTrivia() []Trivia {
	return node.ForTok.LeadingTrivia
}

func (node *FunctionCall) GetLeadingTrivia() []Trivia {
	return node.Name.GetLeadingTrivia()
}

func (node *FunctionExpression) GetLeadingTrivia() []Trivia {
	return node.FuncTok.LeadingTrivia
}

func (node *FunctionStatement) GetLeadingTrivia() []Trivia {
	if node.LocalTok != nil {
		return node.LocalTok.LeadingTrivia
	}
	return node.FuncTok.LeadingTrivia
}

func (node *GotoStatement) GetLeadingTrivia() []Trivia {
	return node.GotoTok.LeadingTrivia
}

func (node *Identifier) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *IfClause) GetLeadingTrivia() []Trivia {
	return node.LeadingTok.LeadingTrivia
}

func (node *IfStatement) GetLeadingTrivia() []Trivia {
	return node.IfTok.LeadingTrivia
}

func (node *IndexExpression) GetLeadingTrivia() []Trivia {
	return node.Prefix.GetLeadingTrivia()
}

func (node *InfixExpression) GetLeadingTrivia() []Trivia {
	return node.Left.GetLeadingTrivia()
}

func (node *LabelStatement) GetLeadingTrivia() []Trivia {
	return node.LeadingLabelTok.LeadingTrivia
}

func (node *Invalid) GetLeadingTrivia() []Trivia {
	return []Trivia{} // TODO:
}

func (node *LocalStatement) GetLeadingTrivia() []Trivia {
	return node.LocalTok.LeadingTrivia
}

func (node *NilLiteral) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *NumberLiteral) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *Pair[T]) GetLeadingTrivia() []Trivia {
	return node.Node.GetLeadingTrivia()
}

func (node *Punctuated[T]) GetLeadingTrivia() []Trivia {
	if len(node.Pairs) == 0 {
		return []Trivia{}
	}
	return node.Pairs[0].GetLeadingTrivia()
}

func (node *PrefixExpression) GetLeadingTrivia() []Trivia {
	return node.Operator.LeadingTrivia
}

func (node *RepeatStatement) GetLeadingTrivia() []Trivia {
	return node.RepeatTok.LeadingTrivia
}

func (node *ReturnStatement) GetLeadingTrivia() []Trivia {
	return node.ReturnTok.LeadingTrivia
}

func (node *SemicolonStatement) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *SimpleTrivia) GetLeadingTrivia() []Trivia {
	return []Trivia{}
}

func (node *StringLiteral) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *TableArrayField) GetLeadingTrivia() []Trivia {
	return node.Expr.GetLeadingTrivia()
}

func (node *TableSimpleKeyField) GetLeadingTrivia() []Trivia {
	return node.Name.LeadingTrivia
}

func (node *TableExpressionKeyField) GetLeadingTrivia() []Trivia {
	return node.LeftBracket.LeadingTrivia
}

func (node *TableLiteral) GetLeadingTrivia() []Trivia {
	return node.LeftBrace.LeadingTrivia
}

func (node *Vararg) GetLeadingTrivia() []Trivia {
	return node.LeadingTrivia
}

func (node *WhileStatement) GetLeadingTrivia() []Trivia {
	return node.WhileTok.LeadingTrivia
}
