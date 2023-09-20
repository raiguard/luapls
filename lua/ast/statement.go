package ast

import (
	"fmt"
	"github.com/raiguard/luapls/lua/token"
)

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
	StartPos token.Pos
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return token.BREAK.String()
}
func (bs *BreakStatement) Pos() token.Pos {
	return bs.StartPos
}
func (bs *BreakStatement) End() token.Pos {
	return bs.StartPos + len(bs.String())
}

type DoStatement struct {
	Body     Block
	StartPos token.Pos
	EndPos   token.Pos
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) String() string {
	return fmt.Sprintf("do\n%s\nend", ds.Body.String())
}
func (ds *DoStatement) Pos() token.Pos {
	return ds.StartPos
}
func (ds *DoStatement) End() token.Pos {
	return ds.EndPos
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
	Name     *Identifier
	Start    Expression
	Finish   Expression
	Step     Expression
	Body     Block
	StartPos token.Pos
	EndPos   token.Pos
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
	return fs.StartPos
}
func (fs *ForStatement) End() token.Pos {
	return fs.EndPos
}

type ForInStatement struct {
	Names    []*Identifier
	Exps     []Expression
	Body     Block
	StartPos token.Pos
	EndPos   token.Pos
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) String() string {
	return fmt.Sprintf("for %s in %s do\n%s\nend", nodeListToString(fs.Names), nodeListToString(fs.Exps), fs.Body.String())
}
func (fs *ForInStatement) Pos() token.Pos {
	return fs.StartPos
}
func (fs *ForInStatement) End() token.Pos {
	return fs.EndPos
}

type FunctionStatement struct {
	Left     Expression
	Params   []*Identifier
	Vararg   bool
	Body     Block
	IsLocal  bool
	StartPos token.Pos
	EndPos   token.Pos
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
	return fs.StartPos
}
func (fs *FunctionStatement) End() token.Pos {
	return fs.EndPos
}

type GotoStatement struct {
	Name     *Identifier
	StartPos token.Pos
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string {
	return fmt.Sprintf("goto %s", gs.Name.String())
}
func (gs *GotoStatement) Pos() token.Pos {
	return gs.StartPos
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
	StartPos  token.Pos
	EndPos    token.Pos
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
	return ic.StartPos
}
func (ic *IfClause) End() token.Pos {
	return ic.EndPos
}

type LabelStatement struct {
	Name     *Identifier
	StartPos token.Pos
	EndPos   token.Pos
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) String() string {
	return fmt.Sprintf("::%s::", ls.Name.String())
}
func (ls *LabelStatement) Pos() token.Pos {
	return ls.StartPos
}
func (ls *LabelStatement) End() token.Pos {
	return ls.EndPos
}

type LocalStatement struct {
	Names    []*Identifier
	Exps     []Expression
	StartPos token.Pos
}

func (ls *LocalStatement) statementNode() {}
func (ls *LocalStatement) String() string {
	return fmt.Sprintf("local %s = %s", nodeListToString(ls.Names), nodeListToString(ls.Exps))
}
func (ls *LocalStatement) Pos() token.Pos {
	return ls.StartPos
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
	StartPos  token.Pos
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) String() string {
	return fmt.Sprintf("repeat\n%s\nuntil %s", rs.Body.String(), rs.Condition.String())
}
func (rs *RepeatStatement) Pos() token.Pos {
	return rs.StartPos
}
func (rs *RepeatStatement) End() token.Pos {
	return rs.Condition.End()
}

type ReturnStatement struct {
	Exps     []Expression
	StartPos token.Pos
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", nodeListToString(rs.Exps))
}
func (rs *ReturnStatement) Pos() token.Pos {
	return rs.StartPos
}
func (rs *ReturnStatement) End() token.Pos {
	if len(rs.Exps) > 0 {
		return rs.Exps[len(rs.Exps)-1].End()
	}
	return rs.StartPos + len(token.RETURN.String())
}

type WhileStatement struct {
	Condition Expression
	Body      Block
	StartPos  token.Pos
	EndPos    token.Pos
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s do\n%s\nend", ws.Condition.String(), ws.Body.String())
}
func (ws *WhileStatement) Pos() token.Pos {
	return ws.StartPos
}
func (ws *WhileStatement) End() token.Pos {
	return ws.EndPos
}
