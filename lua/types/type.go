package types

type Type interface {
	isType()
	String() string
}

type (
	Boolean  struct{}
	Function struct {
		Params []Type
		Return Type
	}
	Number  struct{}
	String  struct{}
	Unknown struct{}
)

func (b *Boolean) isType()  {}
func (f *Function) isType() {}
func (n *Number) isType()   {}
func (s *String) isType()   {}
func (u *Unknown) isType()  {}

func (b *Boolean) String() string { return "boolean" }
func (f *Function) String() string {
	output := "function("
	for _, param := range f.Params {
		output = output + param.String() + ", "
	}
	if len(f.Params) > 0 {
		output = output[0:len(output)-2] + ")"
	}
	output = output + ")"
	if f.Return != nil {
		output = output + " " + f.Return.String()
	}
	return output
}
func (n *Number) String() string  { return "number" }
func (s *String) String() string  { return "string" }
func (u *Unknown) String() string { return "unknown" }
