package lua

import protocol "github.com/tliron/glsp/protocol_3_16"

type Statement interface {
	Range() protocol.Range

	statement()
}

type (
	AssignmentStatement struct {
		Explist []Expression
		Varlist []Variable
	}
	BreakStatement Identifier
	BlockStatement struct {
		Members []Statement
	}
	ChunkStatement struct {
		Members []Statement
	}
	DoStatement struct {
		Block BlockStatement

		Do  Token
		End Token
	}
	ForStatement struct {
		Block    BlockStatement
		Init     Identifier
		InitExp  Expression
		LimitExp Expression
		DeltaExp *Expression // Might not exist

		For Token
		Do  Token
		End Token
	}
	ForInStatement struct {
		Block    BlockStatement
		Explist  []Expression
		Namelist []Identifier

		For Token
		Do  Token
		End Token
	}
	GotoStatement struct {
		Label Identifier

		Goto Token
	}
	IfStatement struct {
		Block BlockStatement
		Exp   Expression

		If   Token
		Then Token
		End  Token
	}
	LabelStatement Identifier
	LocalStatement struct {
		Explist  []Expression
		Namelist []Token

		Token
	}
	ReturnStatement struct {
		Explist []Expression

		Return Token
	}
	RepeatStatement struct {
		Block BlockStatement
		Exp   Expression

		Repeat Token
		Until  Token
	}
	WhileStatement struct {
		Block BlockStatement
		Exp   Expression

		While Token
		Do    Token
		End   Token
	}
)

// TODO:

type Expression interface {
	Range() protocol.Range

	expression()
}

type (
	IndexExpression struct {
		Base    Expression
		Indexer Token // TODO: Make a separate node for this?
		Key     Expression
	}
	MemberExpression struct {
		Base Expression
		Key  Identifier
	}
)

// Variable is a restricted subset of Expression.
type Variable interface {
	variable() // To differentiate variables.
}

type Identifier struct {
	Raw string
	Token
}
