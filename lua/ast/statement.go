package ast

import (
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
func (as *AssignmentStatement) Pos() token.Pos {
	return as.Vars[0].Pos()
}
func (as *AssignmentStatement) End() token.Pos {
	return as.Exps[len(as.Exps)-1].End()
}

type BreakStatement struct {
	StartPos token.Pos `json:"-"`
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) Pos() token.Pos {
	return bs.StartPos
}
func (bs *BreakStatement) End() token.Pos {
	return bs.StartPos + len(token.BREAK.String())
}

type DoStatement struct {
	Body     Block
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (ds *DoStatement) statementNode() {}
func (ds *DoStatement) Pos() token.Pos {
	return ds.StartPos
}
func (ds *DoStatement) End() token.Pos {
	return ds.EndPos
}

type ForStatement struct {
	Name     *Identifier
	Start    Expression
	Finish   Expression
	Step     Expression
	Body     Block
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (fs *ForStatement) statementNode() {}
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
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (fs *ForInStatement) statementNode() {}
func (fs *ForInStatement) Pos() token.Pos {
	return fs.StartPos
}
func (fs *ForInStatement) End() token.Pos {
	return fs.EndPos
}

type FunctionStatement struct {
	Left     Expression
	Params   []*Identifier
	Vararg   *Unit
	Body     Block
	IsLocal  bool
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (fs *FunctionStatement) statementNode() {}
func (fs *FunctionStatement) Pos() token.Pos {
	return fs.StartPos
}
func (fs *FunctionStatement) End() token.Pos {
	return fs.EndPos
}

type GotoStatement struct {
	Name     *Identifier
	StartPos token.Pos `json:"-"`
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) Pos() token.Pos {
	return gs.StartPos
}
func (gs *GotoStatement) End() token.Pos {
	return gs.Name.End()
}

type IfStatement struct {
	Clauses  []*IfClause
	StartPos token.Pos
	EndPos   token.Pos `json:"-"`
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) Pos() token.Pos {
	return is.StartPos
}
func (is *IfStatement) End() token.Pos {
	return is.EndPos
}

type IfClause struct {
	Condition Expression
	Body      Block
	StartPos  token.Pos `json:"-"`
	EndPos    token.Pos `json:"-"`
}

func (ic *IfClause) statementNode() {}
func (ic *IfClause) Pos() token.Pos {
	return ic.StartPos
}
func (ic *IfClause) End() token.Pos {
	return ic.EndPos
}

type LabelStatement struct {
	Name     *Identifier
	StartPos token.Pos `json:"-"`
	EndPos   token.Pos `json:"-"`
}

func (ls *LabelStatement) statementNode() {}
func (ls *LabelStatement) Pos() token.Pos {
	return ls.StartPos
}
func (ls *LabelStatement) End() token.Pos {
	return ls.EndPos
}
func (ls *LabelStatement) leaf() {}

type LocalStatement struct {
	Names    []*Identifier
	Exps     []Expression
	StartPos token.Pos `json:"-"`
}

func (ls *LocalStatement) statementNode() {}
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
	StartPos  token.Pos `json:"-"`
}

func (rs *RepeatStatement) statementNode() {}
func (rs *RepeatStatement) Pos() token.Pos {
	return rs.StartPos
}
func (rs *RepeatStatement) End() token.Pos {
	return rs.Condition.End()
}

type ReturnStatement struct {
	Exps     []Expression
	StartPos token.Pos `json:"-"`
}

func (rs *ReturnStatement) statementNode() {}
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
	StartPos  token.Pos `json:"-"`
	EndPos    token.Pos `json:"-"`
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) Pos() token.Pos {
	return ws.StartPos
}
func (ws *WhileStatement) End() token.Pos {
	return ws.EndPos
}
