package ast

func (node *AssignmentStatement) GetSemanticChildren() []Node {
	return []Node{&node.Vars, &node.Exps}
}

func (node *BooleanLiteral) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *BreakStatement) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *DoStatement) GetSemanticChildren() []Node {
	return []Node{&node.Body}
}

func (node *ForInStatement) GetSemanticChildren() []Node {
	return []Node{&node.Names, &node.Exps, &node.Body}
}

func (node *ForStatement) GetSemanticChildren() (children []Node) {
	children = append(children, &node.Start, &node.Finish)
	if node.Step != nil {
		children = append(children, node.Step)
	}
	children = append(children, &node.Body)
	return children
}

func (node *FunctionCall) GetSemanticChildren() []Node {
	return []Node{node.Name, &node.Args}
}

func (node *FunctionExpression) GetSemanticChildren() []Node {
	return []Node{&node.Params, &node.Body}
}

func (node *FunctionStatement) GetSemanticChildren() []Node {
	return []Node{node.Name, &node.Params, &node.Body}
}

func (node *GotoStatement) GetSemanticChildren() []Node {
	if node.Name != nil {
		return []Node{node.Name}
	}
	return []Node{}
}

func (node *Identifier) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *IfClause) GetSemanticChildren() []Node {
	return []Node{node.Condition, &node.Body}
}

func (node *IfStatement) GetSemanticChildren() (children []Node) {
	for _, clause := range node.Clauses {
		children = append(children, clause)
	}
	return children
}

func (node *IndexExpression) GetSemanticChildren() []Node {
	return []Node{node.Prefix, node.Inner}
}

func (node *InfixExpression) GetSemanticChildren() []Node {
	return []Node{node.Left, node.Right}
}

func (node *LabelStatement) GetSemanticChildren() []Node {
	return []Node{node.Name}
}

func (node *Invalid) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *LocalStatement) GetSemanticChildren() (children []Node) {
	children = append(children, &node.Names)
	if node.Exps != nil {
		children = append(children, node.Exps)
	}
	return children
}

func (node *NilLiteral) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *NumberLiteral) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *Pair[T]) GetSemanticChildren() []Node {
	return []Node{node.Node}
}

func (p *Punctuated[T]) GetSemanticChildren() (n []Node) {
	for i := 0; i < len(p.Pairs); i++ {
		n = append(n, &p.Pairs[i])
	}
	return n
}

func (node *PrefixExpression) GetSemanticChildren() []Node {
	return []Node{node.Right}
}

func (node *RepeatStatement) GetSemanticChildren() []Node {
	return []Node{&node.Body, node.Condition}
}

func (node *ReturnStatement) GetSemanticChildren() []Node {
	return []Node{node.Exps}
}

func (node *SemicolonStatement) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *StringLiteral) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *TableArrayField) GetSemanticChildren() []Node {
	return []Node{node.Expr}
}

func (node *TableSimpleKeyField) GetSemanticChildren() []Node {
	return []Node{&node.Name, node.Expr}
}

func (node *TableExpressionKeyField) GetSemanticChildren() []Node {
	return []Node{node.Name, node.Expr}
}

func (node *TableLiteral) GetSemanticChildren() []Node {
	return []Node{&node.Fields}
}

func (node *Vararg) GetSemanticChildren() []Node {
	return []Node{}
}

func (node *WhileStatement) GetSemanticChildren() []Node {
	return []Node{node.Condition, &node.Body}
}
