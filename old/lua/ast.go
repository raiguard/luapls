package lua

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Node is implemented by all types in the AST.
type Node interface {
	GetRange() protocol.Range
}

// Block is a list of statements to execute.
type Block struct {
	Range protocol.Range
	Stmts []Node
}

func (self *Block) GetRange() protocol.Range {
	return self.Range
}

type BreakStatement struct {
    Range protocol.Range
}

func (stmt *BreakStatement) GetRange() protocol.Range {
	return stmt.Range
}

type DoStatement struct {
    Block *Block
	Range protocol.Range
}

func (self *DoStatement) GetRange() protocol.Range {
    return self.Range
}

// An empty statement contains no logic and is denoted by a semicolon.
type EmptyStatement struct {
	Range protocol.Range
}

func (self *EmptyStatement) GetRange() protocol.Range {
    return self.Range
}

type GotoStatement struct {
    Label Identifier
	Range protocol.Range
}

func (self *GotoStatement) GetRange() protocol.Range {
    return self.Range
}

type Identifier struct {
	Range protocol.Range
    Raw string
}

func (self *Identifier) GetRange() protocol.Range {
    return self.Range
}

type Label struct {
	Range protocol.Range
    Raw string
}

func (self *Label) GetRange() protocol.Range {
    return self.Range
}
