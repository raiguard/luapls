package ast

import (
	"fmt"
	"strings"

	"github.com/raiguard/luapls/lua/token"
)

type Node interface {
	String() string
	Pos() token.Pos
	End() token.Pos
}

type Block struct {
	Stmts []Statement
	pos   token.Pos
}

func (b *Block) String() string {
	var out string
	for _, stmt := range b.Stmts {
		out += stmt.String() + "\n"
	}
	return strings.TrimSpace(out)
}
func (b *Block) Pos() token.Pos {
	return b.pos
}
func (b *Block) End() token.Pos {
	if len(b.Stmts) > 0 {
		return b.pos + b.Stmts[len(b.Stmts)-1].End()
	}
	return b.pos
}

func nodeListToString[T Node](nodes []T) string {
	items := []string{}
	for _, node := range nodes {
		items = append(items, node.String())
	}
	return strings.Join(items, ", ")
}

// Statements

type Statement interface {
	Node
	statementNode()
}

type AssignmentStatement struct {
	Vars []Expression
	Exps []Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", nodeListToString(as.Vars), nodeListToString(as.Exps))
}
func (as *AssignmentStatement) Pos() token.Pos {
	return as.Vars[0].Pos()
}
func (as *AssignmentStatement) End() token.Pos {
	return as.Exps[len(as.Exps)-1].Pos()
}

type BreakStatement struct {
	pos token.Pos
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return token.BREAK.String()
}
func (bs *BreakStatement) Pos() token.Pos {
	return bs.pos
}
func (bs *BreakStatement) End() token.Pos {
	return bs.pos + len(bs.String())
}

type DoStatement struct {
	Body Block
	do   token.Pos
	end  token.Pos
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) String() string {
	return fmt.Sprintf("do\n%s\nend", ds.Body.String())
}
func (ds *DoStatement) Pos() token.Pos {
	return ds.do
}
func (ds *DoStatement) End() token.Pos {
	return ds.end
}

type ExpressionStatement struct {
	Exp Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%s", es.Exp.String())
}
func (es *ExpressionStatement) Pos() token.Pos {
	return es.Exp.Pos()
}
func (es *ExpressionStatement) End() token.Pos {
	return es.Exp.End()
}

type ForStatement struct {
	Name   *Identifier
	Start  Expression
	Finish Expression
	Step   Expression
	Body   Block
	pos    token.Pos
	end    token.Pos
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	if fs.Step != nil {
		return fmt.Sprintf(
			"for %s = %s, %s, %s do\n%s\nend",
			fs.Name.String(),
			fs.Start.String(),
			fs.Finish.String(),
			fs.Step.String(),
			fs.Body.String(),
		)
	} else {
		return fmt.Sprintf(
			"for %s = %s, %s do\n%s\nend",
			fs.Name.String(),
			fs.Start.String(),
			fs.Finish.String(),
			fs.Body.String(),
		)
	}
}
func (fs *ForStatement) Pos() token.Pos {
	return fs.pos
}
func (fs *ForStatement) End() token.Pos {
	return fs.end
}

type ForInStatement struct {
	Names []*Identifier
	Exps  []Expression
	Body  Block
	pos   token.Pos
	end   token.Pos
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", nodeListToString(fs.Names), nodeListToString(fs.Exps), fs.Body.String())
}
func (fs *ForInStatement) Pos() token.Pos {
	return fs.pos
}
func (fs *ForInStatement) End() token.Pos {
	return fs.end
}

type FunctionStatement struct {
	Left    Expression
	Params  []*Identifier
	Vararg  bool
	Body    Block
	IsLocal bool
	pos     token.Pos
	end     token.Pos
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) String() string {
	localStr := ""
	if fs.IsLocal {
		localStr = "local "
	}
	return fmt.Sprintf(
		"%sfunction %s(%s)\n%s\nend",
		localStr,
		fs.Left.String(),
		nodeListToString(fs.Params),
		fs.Body.String(),
	)
}
func (fs *FunctionStatement) Pos() token.Pos {
	return fs.pos
}
func (fs *FunctionStatement) End() token.Pos {
	return fs.end
}

type GotoStatement struct {
	Name *Identifier
	pos  token.Pos
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("goto %s", gs.Name.String())
}
func (gs *GotoStatement) Pos() token.Pos {
	return gs.pos
}
func (gs *GotoStatement) End() token.Pos {
	return gs.Name.End()
}

type IfStatement struct {
	Clauses []*IfClause
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("%send", nodeListToString(is.Clauses))
}
func (is *IfStatement) Pos() token.Pos {
	return is.Clauses[0].Pos()
}
func (is *IfStatement) End() token.Pos {
	return is.Clauses[len(is.Clauses)-1].End()
}

type IfClause struct {
	Condition Expression
	Body      Block
	pos       token.Pos
	end       token.Pos
}

func (ic *IfClause) statementNode() {}
func (ic *IfClause) String() string {
	// TODO: elseif
	if ic.Condition != nil {
		return fmt.Sprintf("if %s then\n%s\n", ic.Condition.String(), ic.Body.String())
	} else {
		return fmt.Sprintf("else\n%s\n", ic.Body.String())
	}
}
func (ic *IfClause) Pos() token.Pos {
	return ic.pos
}
func (ic *IfClause) End() token.Pos {
	return ic.end
}

type LabelStatement struct {
	Name *Identifier
	pos  token.Pos
	end  token.Pos
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return fmt.Sprintf("::%s::", ls.Name.String())
}
func (ls *LabelStatement) Pos() token.Pos {
	return ls.pos
}
func (ls *LabelStatement) End() token.Pos {
	return ls.end
}

type LocalStatement struct {
	Names []*Identifier
	Exps  []Expression
	pos   token.Pos
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("local %s = %s", nodeListToString(ls.Names), nodeListToString(ls.Exps))
}
func (ls *LocalStatement) Pos() token.Pos {
	return ls.pos
}
func (ls *LocalStatement) End() token.Pos {
	if len(ls.Exps) > 0 {
		return ls.Exps[len(ls.Exps)-1].End()
	}
	return ls.Names[len(ls.Names)-1].End()
}

type RepeatStatement struct {
	Body      Block
	Condition Expression
	pos       token.Pos
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) String() string {
	return fmt.Sprintf("repeat\n%s\nuntil %s", rs.Body.String(), rs.Condition.String())
}
func (rs *RepeatStatement) Pos() token.Pos {
	return rs.pos
}
func (rs *RepeatStatement) End() token.Pos {
	return rs.Condition.End()
}

type ReturnStatement struct {
	Exps []Expression
	pos  token.Pos
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", nodeListToString(rs.Exps))
}
func (rs *ReturnStatement) Pos() token.Pos {
	return rs.pos
}
func (rs *ReturnStatement) End() token.Pos {
	if len(rs.Exps) > 0 {
		return rs.Exps[len(rs.Exps)-1].End()
	}
	return rs.pos + len(token.RETURN.String())
}

type WhileStatement struct {
	Condition Expression
	Body      Block
	pos       token.Pos
	end       token.Pos
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s do\n%s\nend", ws.Condition.String(), ws.Body.String())
}
func (ws *WhileStatement) Pos() token.Pos {
	return ws.pos
}
func (ws *WhileStatement) End() token.Pos {
	return ws.end
}

// Expressions

type Expression interface {
	Node
	expressionNode()
}

type BinaryExpression struct {
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator.String(), be.Right.String())
}
func (be *BinaryExpression) Pos() token.Pos {
	return be.Left.Pos()
}
func (be *BinaryExpression) End() token.Pos {
	return be.Right.End()
}

type FunctionCall struct {
	Left Expression
	Args []Expression
	end  token.Pos
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) String() string {
	return fmt.Sprintf("%s(%s)", fc.Left.String(), nodeListToString(fc.Args))
}
func (fc *FunctionCall) Pos() token.Pos {
	return fc.Left.Pos()
}
func (fc *FunctionCall) End() token.Pos {
	return fc.end
}

type FunctionExpression struct {
	Params []*Identifier
	Vararg bool
	Body   Block
	pos    token.Pos
	end    token.Pos
}

func (fe *FunctionExpression) expressionNode() {}
func (fe *FunctionExpression) String() string {
	return fmt.Sprintf("function(%s)\n%s\nend", nodeListToString(fe.Params), fe.Body.String())
}
func (fe *FunctionExpression) Pos() token.Pos {
	return fe.pos
}
func (fe *FunctionExpression) End() token.Pos {
	return fe.end
}

type IndexExpression struct {
	Left       Expression
	Inner      Expression
	IsBrackets bool
	IsColon    bool
	end        token.Pos
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	if ie.IsBrackets {
		return fmt.Sprintf("%s[%s]", ie.Left.String(), ie.Inner.String())
	} else {
		return fmt.Sprintf("%s.%s", ie.Left.String(), ie.Inner.String())
	}
}
func (ie *IndexExpression) Pos() token.Pos {
	return ie.Left.Pos()
}
func (ie *IndexExpression) End() token.Pos {
	return ie.end
}

type UnaryExpression struct {
	Operator token.TokenType
	Right    Expression
	pos      token.Pos
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", ue.Operator, ue.Right.String())
}
func (ue *UnaryExpression) Pos() token.Pos {
	return ue.pos
}
func (ue *UnaryExpression) End() token.Pos {
	return ue.Right.End()
}

// Literals (also expressions)

type BooleanLiteral struct {
	Value bool
	pos   token.Pos
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}
func (bl *BooleanLiteral) Pos() token.Pos {
	return bl.pos
}
func (bl *BooleanLiteral) End() token.Pos {
	return bl.pos + len(bl.String())
}

type Identifier struct {
	Literal string
	pos     token.Pos
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Literal }
func (i *Identifier) Pos() token.Pos {
	return i.pos
}
func (i *Identifier) End() token.Pos {
	return i.pos + len(i.String())
}

type NumberLiteral struct {
	Literal string
	Value   float64
	pos     token.Pos
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Literal }
func (nl *NumberLiteral) Pos() token.Pos {
	return nl.pos
}
func (nl *NumberLiteral) End() token.Pos {
	return nl.pos + len(nl.String())
}

type StringLiteral struct {
	Value string
	pos   token.Pos
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }
func (sl *StringLiteral) Pos() token.Pos {
	return sl.pos
}
func (sl *StringLiteral) End() token.Pos {
	return sl.pos + len(sl.String())
}

type TableLiteral struct {
	Fields []*TableField
	pos    token.Pos
	end    token.Pos
}

func (tl *TableLiteral) expressionNode() {}
func (tl *TableLiteral) String() string  { return fmt.Sprintf("{ %s }", nodeListToString(tl.Fields)) }
func (tl *TableLiteral) Pos() token.Pos {
	return tl.pos
}
func (tl *TableLiteral) End() token.Pos {
	return tl.end
}

type TableField struct {
	Key   Expression
	Value Expression
	pos   token.Pos // Needed in case of bracketed keys
}

func (tf TableField) String() string {
	if tf.Key == nil {
		return tf.Value.String()
	}
	if ident, ok := tf.Key.(*Identifier); ok {
		return fmt.Sprintf("%s = %s", ident.String(), tf.Value.String())
	}
	return fmt.Sprintf("[%s] = %s", tf.Key.String(), tf.Value.String())
}
func (rf *TableField) Pos() token.Pos {
	return rf.pos
}
func (rf *TableField) End() token.Pos {
	return rf.Value.End()
}

type Vararg struct {
	pos token.Pos
}

func (va *Vararg) expressionNode() {}
func (va *Vararg) String() string  { return token.TokenStr[token.VARARG] }
func (va *Vararg) Pos() token.Pos {
	return va.pos
}
func (va *Vararg) End() token.Pos {
	return va.pos + len(token.VARARG.String())
}
