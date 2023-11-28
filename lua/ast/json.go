package ast

import (
	"encoding/json"

	"github.com/raiguard/luapls/lua/token"
)

func (node *AssignmentStatement) MarshalJSON() ([]byte, error) {
	type Alias AssignmentStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "AssignmentStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "Block",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *BooleanLiteral) MarshalJSON() ([]byte, error) {
	type Alias BooleanLiteral
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "BooleanLiteral",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *BreakStatement) MarshalJSON() ([]byte, error) {
	type Alias BreakStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "BreakStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *DoStatement) MarshalJSON() ([]byte, error) {
	type Alias DoStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "DoStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *ForInStatement) MarshalJSON() ([]byte, error) {
	type Alias ForInStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "ForInStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *ForStatement) MarshalJSON() ([]byte, error) {
	type Alias ForStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "ForStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *FunctionCall) MarshalJSON() ([]byte, error) {
	type Alias FunctionCall
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "FunctionCall",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *FunctionExpression) MarshalJSON() ([]byte, error) {
	type Alias FunctionExpression
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "FunctionExpression",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *FunctionStatement) MarshalJSON() ([]byte, error) {
	type Alias FunctionStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "FunctionStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *GotoStatement) MarshalJSON() ([]byte, error) {
	type Alias GotoStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "GotoStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *Identifier) MarshalJSON() ([]byte, error) {
	type Alias Identifier
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "Identifier",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *IfClause) MarshalJSON() ([]byte, error) {
	type Alias IfClause
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "IfClause",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *IfStatement) MarshalJSON() ([]byte, error) {
	type Alias IfStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "IfStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *IndexExpression) MarshalJSON() ([]byte, error) {
	type Alias IndexExpression
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "IndexExpression",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *InfixExpression) MarshalJSON() ([]byte, error) {
	type Alias InfixExpression
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "InfixExpression",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *LabelStatement) MarshalJSON() ([]byte, error) {
	type Alias LabelStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "LabelStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *LocalStatement) MarshalJSON() ([]byte, error) {
	type Alias LocalStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "LocalStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *NumberLiteral) MarshalJSON() ([]byte, error) {
	type Alias NumberLiteral
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "NumberLiteral",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *PrefixExpression) MarshalJSON() ([]byte, error) {
	type Alias PrefixExpression
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "PrefixExpression",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *RepeatStatement) MarshalJSON() ([]byte, error) {
	type Alias RepeatStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "RepeatStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *ReturnStatement) MarshalJSON() ([]byte, error) {
	type Alias ReturnStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "ReturnStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *StringLiteral) MarshalJSON() ([]byte, error) {
	type Alias StringLiteral
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "StringLiteral",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *TableField) MarshalJSON() ([]byte, error) {
	type Alias TableField
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "TableField",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *TableLiteral) MarshalJSON() ([]byte, error) {
	type Alias TableLiteral
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "TableLiteral",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *Vararg) MarshalJSON() ([]byte, error) {
	type Alias Vararg
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "Vararg",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}

func (node *WhileStatement) MarshalJSON() ([]byte, error) {
	type Alias WhileStatement
	return json.Marshal(&struct {
		Type  string
		Range token.Range
		*Alias
	}{
		Type:  "WhileStatement",
		Range: token.Range{node.Pos(), node.End()},
		Alias: (*Alias)(node),
	})
}
