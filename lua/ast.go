package lua

// A Node of the AST.
type Node interface {
	Pos() TokenPos
	End() TokenPos
}

// STATEMENTS

type Stmt interface {
	Node
	stmtNode()
}

type (
	AssignmentStmt struct {
		Explist ExprList
		Varlist ExprList
	}
	BreakStmt Identifier
	BlockStmt struct {
		Members []Stmt
		Start   TokenPos
	}
	DoStmt struct {
		Block      BlockStmt
		Do, EndTok TokenPos
	}
	ForStmt struct {
		Block           BlockStmt
		Init            Identifier
		InitExp         Expr
		LimitExp        Expr
		DeltaExp        *Expr // Might not exist
		For, Do, EndTok TokenPos
	}
	FunctionStmt struct {
		Body             BlockStmt
		Funcname         Expr
		Params           IdentifierList
		Function, EndTok TokenPos
	}
	ForInStmt struct {
		Block           BlockStmt
		Explist         ExprList
		Namelist        IdentifierList
		For, Do, EndTok TokenPos
	}
	GotoStmt struct {
		Label Identifier
		Goto  TokenPos
	}
	// TODO: Elseif / else
	IfStmt struct {
		Block            BlockStmt
		Exp              Expr
		If, Then, EndTok TokenPos
	}
	LabelStmt         Identifier
	LocalFunctionStmt struct {
		Body                    BlockStmt
		Funcname                Expr
		Params                  IdentifierList
		Local, Function, EndTok TokenPos
	}
	LocalStmt struct {
		Explist       ExprList
		Namelist      IdentifierList
		Local, Equals TokenPos
	}
	ReturnStmt struct {
		Explist ExprList
		Return  TokenPos
	}
	RepeatStmt struct {
		Block         BlockStmt
		Exp           Expr
		Repeat, Until TokenPos
	}
	WhileStmt struct {
		Block             BlockStmt
		Exp               Expr
		While, Do, EndTok TokenPos
	}
)

func (x *AssignmentStmt) Pos() TokenPos { return x.Varlist.Pos() }
func (x *AssignmentStmt) End() TokenPos { return x.Explist.End() }
func (x *AssignmentStmt) stmtNode()     {}

func (x *BreakStmt) Pos() TokenPos { return x.Pos() } // TODO: Identifier Pos()
func (x *BreakStmt) End() TokenPos { return x.End() } // TODO: Identifier End()
func (x *BreakStmt) stmtNode()     {}

func (x *BlockStmt) Pos() TokenPos { return x.Start }
func (x *BlockStmt) End() TokenPos {
	stmtsLen := len(x.Members)
	if stmtsLen == 0 {
		return x.Start
	}
	return x.Members[stmtsLen-1].End()
}
func (x *BlockStmt) stmtNode() {}

func (x *DoStmt) Pos() TokenPos { return x.Do }
func (x *DoStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *DoStmt) stmtNode() {}

func (x *ForStmt) Pos() TokenPos { return x.For }
func (x *ForStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *ForStmt) stmtNode() {}

func (x *FunctionStmt) Pos() TokenPos { return x.Function }
func (x *FunctionStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *FunctionStmt) stmtNode() {}

func (x *ForInStmt) Pos() TokenPos { return x.For }
func (x *ForInStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *ForInStmt) stmtNode() {}

func (x *GotoStmt) Pos() TokenPos { return x.Goto }
func (x *GotoStmt) End() TokenPos { return x.Label.End() }
func (x *GotoStmt) stmtNode()     {}

func (x *IfStmt) Pos() TokenPos { return x.If }
func (x *IfStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *IfStmt) stmtNode() {}

func (x *LabelStmt) Pos() TokenPos { return x.Pos() } // TODO: Identifier Pos()
func (x *LabelStmt) End() TokenPos { return x.End() } // TODO: Identifier End()
func (x *LabelStmt) stmtNode()     {}

func (x *LocalFunctionStmt) Pos() TokenPos { return x.Local }
func (x *LocalFunctionStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *LocalFunctionStmt) stmtNode() {}

func (x *LocalStmt) Pos() TokenPos { return x.Local }
func (x *LocalStmt) End() TokenPos { return x.Explist.End() }
func (x *LocalStmt) stmtNode()     {}

func (x *ReturnStmt) Pos() TokenPos { return x.Return }
func (x *ReturnStmt) End() TokenPos { return x.Explist.End() }
func (x *ReturnStmt) stmtNode()     {}

func (x *RepeatStmt) Pos() TokenPos { return x.Repeat }
func (x *RepeatStmt) End() TokenPos { return x.Exp.End() }
func (x *RepeatStmt) stmtNode()     {}

func (x *WhileStmt) Pos() TokenPos { return x.While }
func (x *WhileStmt) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *WhileStmt) stmtNode() {}

// EXPRESSIONS

type Expr interface {
	Node
	exprNode()
}

type (
	BinopExpr struct {
		Left  Expr
		Op    Token
		Right Expr
	}
	ClassMemberExpr struct { // Index with :
		Base Expr
		Key  Identifier
	}
	FunctionExpr struct {
		Body   BlockStmt
		Params IdentifierList

		Function, EndTok TokenPos
	}
	FunctionCallExpr struct {
		Base Expr
		Args ExprList
		// TODO: Parenthesis are optional :(
		Lparen, Rparen TokenPos
	}
	IndexExpr struct { // Index with []
		Base Expr
		Key  Expr
	}
	MemberExpr struct { // Index with .
		Base Expr
		Key  Identifier
	}
	TableConstructorExpr struct {
		Fields      ExprList
		Open, Close TokenPos
	}
	UnopExpr struct {
		Expr Expr
		Op   Token
	}
)

func (x *BinopExpr) Pos() TokenPos { return x.Left.Pos() }
func (x *BinopExpr) End() TokenPos { return x.Right.End() }
func (x *BinopExpr) exprNode()     {}

func (x *ClassMemberExpr) Pos() TokenPos { return x.Base.Pos() }
func (x *ClassMemberExpr) End() TokenPos { return x.Key.End() }
func (x *ClassMemberExpr) exprNode()     {}

func (x *FunctionExpr) Pos() TokenPos { return x.Function }
func (x *FunctionExpr) End() TokenPos {
	return TokenPos{Col: x.EndTok.Col + 2, Line: x.EndTok.Line}
}
func (x *FunctionExpr) exprNode() {}

func (x *FunctionCallExpr) Pos() TokenPos { return x.Base.Pos() }
func (x *FunctionCallExpr) End() TokenPos {
	return TokenPos{Col: x.Rparen.Col + 1, Line: x.Rparen.Line}
}
func (x *FunctionCallExpr) exprNode() {}

func (x *IndexExpr) Pos() TokenPos { return x.Base.Pos() }
func (x *IndexExpr) End() TokenPos { return x.Key.End() }
func (x *IndexExpr) exprNode()     {}

func (x *MemberExpr) Pos() TokenPos { return x.Base.Pos() }
func (x *MemberExpr) End() TokenPos { return x.Key.End() }
func (x *MemberExpr) exprNode()     {}

func (x *TableConstructorExpr) Pos() TokenPos { return x.Open }
func (x *TableConstructorExpr) End() TokenPos {
	return TokenPos{Col: x.Close.Col + 1, Line: x.Close.Line}
}
func (x *TableConstructorExpr) exprNode() {}

func (x *UnopExpr) Pos() TokenPos { return x.Op.Pos }
func (x *UnopExpr) End() TokenPos { return x.Expr.End() }
func (x *UnopExpr) exprNode()     {}

// A list of expressions.
type ExprList struct {
	Exprs []Expr
	Start TokenPos
}

func (l *ExprList) Pos() TokenPos { return l.Start }
func (l *ExprList) End() TokenPos {
	exprsLen := len(l.Exprs)
	if exprsLen == 0 {
		return l.Start
	}
	return l.Exprs[exprsLen-1].End()
}

// IDENTIFIERS

type Identifier struct {
	Raw string
	TokenPos
}

func (x *Identifier) Pos() TokenPos { return x.TokenPos }
func (x *Identifier) End() TokenPos {
	return TokenPos{Col: x.TokenPos.Col + len(x.Raw), Line: x.TokenPos.Line}
}

// A list of identifiers.
type IdentifierList struct {
	Node
	Identifiers []Identifier
	Start, Stop TokenPos
}

func (l *IdentifierList) Pos() TokenPos { return l.Start }
func (l *IdentifierList) End() TokenPos {
	identsLen := len(l.Identifiers)
	if identsLen == 0 {
		return l.Start
	}
	return l.Identifiers[identsLen-1].End()
}
