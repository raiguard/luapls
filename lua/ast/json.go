package ast

import (
	"encoding/json"

	"github.com/raiguard/luapls/lua/token"
)

func (node *AssignmentStatement) MarshalJSON() ([]byte, error) {
	type Alias AssignmentStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "AssignmentStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *BinaryExpression) MarshalJSON() ([]byte, error) {
	type Alias BinaryExpression
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "BinaryExpression",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

// If we use a pointer receiver then it doesn't do it on the parser.File object
func (node Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		Alias
	}{
		NodeKind: "Block",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (Alias)(node),
	})
}

// func (node *Block) MarshalJSON() ([]byte, error) {
// 	type Alias Block
// 	return json.Marshal(&struct {
// 		NodeKind string
// 		Range    token.Range
// 		*Alias
// 	}{
// 		NodeKind: "Block",
// 		Range:    token.Range{node.Pos(), node.End()},
// 		Alias:    (*Alias)(node),
// 	})
// }

func (node *BooleanLiteral) MarshalJSON() ([]byte, error) {
	type Alias BooleanLiteral
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "BooleanLiteral",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *BreakStatement) MarshalJSON() ([]byte, error) {
	type Alias BreakStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "BreakStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *DoStatement) MarshalJSON() ([]byte, error) {
	type Alias DoStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "DoStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *ExpressionStatement) MarshalJSON() ([]byte, error) {
	type Alias ExpressionStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "ExpressionStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *ForInStatement) MarshalJSON() ([]byte, error) {
	type Alias ForInStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "ForInStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *ForStatement) MarshalJSON() ([]byte, error) {
	type Alias ForStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "ForStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *FunctionCall) MarshalJSON() ([]byte, error) {
	type Alias FunctionCall
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "FunctionCall",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *FunctionExpression) MarshalJSON() ([]byte, error) {
	type Alias FunctionExpression
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "FunctionExpression",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *FunctionStatement) MarshalJSON() ([]byte, error) {
	type Alias FunctionStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "FunctionStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *GotoStatement) MarshalJSON() ([]byte, error) {
	type Alias GotoStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "GotoStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *Identifier) MarshalJSON() ([]byte, error) {
	type Alias Identifier
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "Identifier",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *IfClause) MarshalJSON() ([]byte, error) {
	type Alias IfClause
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "IfClause",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *IfStatement) MarshalJSON() ([]byte, error) {
	type Alias IfStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "IfStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *IndexExpression) MarshalJSON() ([]byte, error) {
	type Alias IndexExpression
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "IndexExpression",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *LabelStatement) MarshalJSON() ([]byte, error) {
	type Alias LabelStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "LabelStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *LocalStatement) MarshalJSON() ([]byte, error) {
	type Alias LocalStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "LocalStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *NumberLiteral) MarshalJSON() ([]byte, error) {
	type Alias NumberLiteral
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "NumberLiteral",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *RepeatStatement) MarshalJSON() ([]byte, error) {
	type Alias RepeatStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "RepeatStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *ReturnStatement) MarshalJSON() ([]byte, error) {
	type Alias ReturnStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "ReturnStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *StringLiteral) MarshalJSON() ([]byte, error) {
	type Alias StringLiteral
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "StringLiteral",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *TableField) MarshalJSON() ([]byte, error) {
	type Alias TableField
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "TableField",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *TableLiteral) MarshalJSON() ([]byte, error) {
	type Alias TableLiteral
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "TableLiteral",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *UnaryExpression) MarshalJSON() ([]byte, error) {
	type Alias UnaryExpression
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "UnaryExpression",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *Vararg) MarshalJSON() ([]byte, error) {
	type Alias Vararg
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "Vararg",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}

func (node *WhileStatement) MarshalJSON() ([]byte, error) {
	type Alias WhileStatement
	return json.Marshal(&struct {
		NodeKind string
		Range    token.Range
		*Alias
	}{
		NodeKind: "WhileStatement",
		Range:    token.Range{node.Pos(), node.End()},
		Alias:    (*Alias)(node),
	})
}
