package ast

import "github.com/raiguard/luapls/lua/token"

type TableField interface {
	Node
	tableFieldNode()
}

type TableArrayField struct {
	Expr Expression
}

func (taf *TableArrayField) tableFieldNode() {}
func (taf *TableArrayField) Pos() token.Pos {
	return taf.Expr.Pos()
}
func (taf *TableArrayField) End() token.Pos {
	return taf.Expr.End()
}

type TableSimpleKeyField struct {
	Name      Identifier
	AssignTok Unit
	Expr      Expression
}

func (tf *TableSimpleKeyField) tableFieldNode() {}
func (tf *TableSimpleKeyField) Pos() token.Pos {
	return tf.Name.Pos()
}
func (tf *TableSimpleKeyField) End() token.Pos {
	return tf.Expr.End()
}

type TableExpressionKeyField struct {
	LeftBracket  Unit
	Name         Expression
	RightBracket Unit
	AssignTok    Unit
	Expr         Expression
}

func (tf *TableExpressionKeyField) tableFieldNode() {}
func (tf *TableExpressionKeyField) Pos() token.Pos {
	return tf.LeftBracket.Pos()
}
func (tf *TableExpressionKeyField) End() token.Pos {
	return tf.Expr.End()
}
