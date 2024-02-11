package ast

type Type interface {
	isType()
}

type (
	Boolean  struct{}
	Function struct{}
	Nil      struct{}
	Number   struct{}
	String   struct{}
	Unknown  struct{}
)

func (b *Boolean) isType()  {}
func (f *Function) isType() {}
func (n *Nil) isType()      {}
func (n *Number) isType()   {}
func (s *String) isType()   {}
func (u *Unknown) isType()  {}

// Statements
func (fs *FunctionStatement) Type() Type { return &Function{} }

// func (as *AssignmentStatement) Type() Type { return as.Exps[0].Type() }
// func (ls *LocalStatement) Type() Type      { return ls.Exps[0].Type() }

// Expressions
func (self *BooleanLiteral) Type() Type     { return &Boolean{} }
func (self *FunctionCall) Type() Type       { return &Unknown{} }
func (self *FunctionExpression) Type() Type { return &Function{} }
func (self *Identifier) Type() Type         { return &Unknown{} }
func (self *IndexExpression) Type() Type    { return &Unknown{} }
func (self *Invalid) Type() Type            { return &Unknown{} }
func (self *InfixExpression) Type() Type    { return &Unknown{} }
func (self *NumberLiteral) Type() Type      { return &Number{} }
func (self *PrefixExpression) Type() Type   { return &Unknown{} }
func (self *StringLiteral) Type() Type      { return &String{} }
func (self *TableLiteral) Type() Type       { return &Unknown{} }
func (self *Vararg) Type() Type             { return &Unknown{} }
