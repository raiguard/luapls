package types

import "fmt"

type Type interface {
	isType()
	String() string
}

type (
	Any      struct{}
	Boolean  struct{}
	Function struct {
		Params []FunctionParameter
		Return Type
	}
	Number  struct{}
	String  struct{}
	Unknown struct{}
)

func (a *Any) isType()      {}
func (b *Boolean) isType()  {}
func (f *Function) isType() {}
func (n *Number) isType()   {}
func (s *String) isType()   {}
func (u *Unknown) isType()  {}

func (b *Any) String() string     { return "any" }
func (b *Boolean) String() string { return "boolean" }
func (f *Function) String() string {
	output := "function("
	for _, param := range f.Params {
		output = output + param.String() + ", "
	}
	if len(f.Params) > 0 {
		output = output[0 : len(output)-2]
	}
	output = output + ")"
	if f.Return != nil {
		output = output + " -> " + f.Return.String()
	}
	return output
}
func (n *Number) String() string  { return "number" }
func (s *String) String() string  { return "string" }
func (u *Unknown) String() string { return "unknown" }

type FunctionParameter struct {
	Name string
	Type Type
}

func (f *FunctionParameter) String() string {
	typ := f.Type
	if typ == nil {
		typ = &Unknown{}
	}
	return fmt.Sprintf("%s: %s", f.Name, f.Type)
}
