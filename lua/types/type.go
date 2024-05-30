package types

import (
	"fmt"
	"strings"
)

type Type interface {
	isType()
	String() string
}

type (
	Any      struct{}
	Boolean  struct{}
	Function struct {
		Params []NameAndType
		Return Type
	}
	Number struct{}
	String struct{}
	Table  struct {
		Fields []NameAndType
	}
	Unknown struct{}
)

func (a *Any) isType()      {}
func (b *Boolean) isType()  {}
func (f *Function) isType() {}
func (n *Number) isType()   {}
func (s *String) isType()   {}
func (t *Table) isType()    {}
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
		output = output + " → " + f.Return.String()
	}
	return output
}
func (n *Number) String() string { return "number" }
func (s *String) String() string { return "string" }
func (t *Table) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i := 0; i < len(t.Fields); i++ {
		if i > 0 {
			fmt.Fprint(&sb, ", ")
		}
		fmt.Fprintf(&sb, "%s", &t.Fields[i])
	}
	sb.WriteByte('}')

	return sb.String()
}
func (u *Unknown) String() string { return "unknown" }

type NameAndType struct {
	Name string
	Type Type
}

func (n *NameAndType) String() string {
	typ := n.Type
	if typ == nil {
		typ = &Unknown{}
	}
	return fmt.Sprintf("%s: %s", n.Name, n.Type)
}
