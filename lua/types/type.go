package types

type Type interface {
	isType()
	String() string
}

type (
	Boolean  struct{}
	Function struct{}
	Number   struct{}
	String   struct{}
	Unknown  struct{}
)

func (b *Boolean) isType()  {}
func (f *Function) isType() {}
func (n *Number) isType()   {}
func (s *String) isType()   {}
func (u *Unknown) isType()  {}

func (b *Boolean) String() string  { return "boolean" }
func (f *Function) String() string { return "function" }
func (n *Number) String() string   { return "number" }
func (s *String) String() string   { return "string" }
func (u *Unknown) String() string  { return "unknown" }
